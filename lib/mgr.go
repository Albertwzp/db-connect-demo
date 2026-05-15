package lib

import (
	"context"
	"fmt"
)

var backends = map[string]Driver{}

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

func HealthAll(ctx context.Context) map[string]string {
	res := map[string]string{}
	for name, d := range backends {
		if err := d.HealthCheck(ctx); err != nil {
			res[name] = err.Error()
		} else {
			res[name] = "ok"
		}
	}
	return res
}

func QueryBackend(ctx context.Context, name, query string) ([]map[string]interface{}, error) {
	d, ok := backends[name]
	if !ok { return nil, fmt.Errorf("backend not found: %s", name) }
	return d.Query(ctx, query)
}

func CloseAllBackends() {
	for _, d := range backends {
		d.Close()
	}
	backends = map[string]Driver{}
}
