package lib

import (
	"context"
	"fmt"
)

var (
	backends      = map[string]Driver{}
	failedBackends = map[string]string{}
)

// RegisterBackend registers and opens the backend. Returns error on failure.
// Callers may choose to record the failure and continue.
func RegisterBackend(name, driverName, dsn string) error {
	if _, ok := backends[name]; ok {
		return fmt.Errorf("backend already registered: %s", name)
	}
	drv, err := NewDriver(driverName)
	if err != nil { return err }
	if err := drv.Open(dsn); err != nil { return err }
	backends[name] = drv
	return nil
}

// MarkBackendFailed records a backend that failed during registration so it appears in /ping
func MarkBackendFailed(name string, reason error) {
	failedBackends[name] = reason.Error()
}

// HealthAll returns health for registered and failed backends
func HealthAll(ctx context.Context) map[string]string {
	res := map[string]string{}
	for name, d := range backends {
		if err := d.HealthCheck(ctx); err != nil {
			res[name] = err.Error()
		} else {
			res[name] = "ok"
		}
	}
	for name, why := range failedBackends {
		res[name] = why
	}
	return res
}

// QueryBackend returns query results or an error if backend missing or unsupported
func QueryBackend(ctx context.Context, name, query string) ([]map[string]interface{}, error) {
	if why, ok := failedBackends[name]; ok {
		return nil, fmt.Errorf("backend %s registration failed: %s", name, why)
	}
	d, ok := backends[name]
	if !ok { return nil, fmt.Errorf("backend not found: %s", name) }
	return d.Query(ctx, query)
}

func CloseAllBackends() {
	for _, d := range backends {
		d.Close()
	}
	backends = map[string]Driver{}
	failedBackends = map[string]string{}
}
