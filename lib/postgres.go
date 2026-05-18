package lib

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

// Driver is a minimal interface for load testing operations
type Driver interface {
	Open(dsn string) error
	Close() error
	RunOp(ctx context.Context, query string) (time.Duration, error)
	HealthCheck(ctx context.Context) error
	Query(ctx context.Context, query string) ([]map[string]interface{}, error)
}

var registry = map[string]func() Driver{}

func Register(name string, f func() Driver) {
	registry[name] = f
}

func NewDriver(name string) (Driver, error) {
	if f, ok := registry[name]; ok {
		return f(), nil
	}
	return nil, fmt.Errorf("driver not found: %s", name)
}

// --- Postgres driver implementation (example) ---

type PostgresDriver struct {
	db *sql.DB
}

func (p *PostgresDriver) Open(dsn string) error {
	if dsn == "" {
		return errors.New("empty dsn")
	}
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return err
	}
	// optional: set connection pool defaults
	db.SetMaxOpenConns(50)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(time.Minute * 5)
	p.db = db
	return nil
}

func (p *PostgresDriver) Close() error {
	if p.db != nil {
		return p.db.Close()
	}
	return nil
}

func (p *PostgresDriver) RunOp(ctx context.Context, query string) (time.Duration, error) {
	if p.db == nil {
		return 0, errors.New("db not open")
	}
	start := time.Now()
	_, err := p.db.ExecContext(ctx, query)
	return time.Since(start), err
}

func (p *PostgresDriver) HealthCheck(ctx context.Context) error {
	if p.db == nil {
		return errors.New("db not open")
	}
	return p.db.PingContext(ctx)
}

func (p *PostgresDriver) Query(ctx context.Context, query string) ([]map[string]interface{}, error) {
	if p.db == nil {
		return nil, errors.New("db not open")
	}
	rows, err := p.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	results := []map[string]interface{}{}
	for rows.Next() {
		vals := make([]interface{}, len(cols))
		ptrs := make([]interface{}, len(cols))
		for i := range vals {
			ptrs[i] = &vals[i]
		}
		if err := rows.Scan(ptrs...); err != nil {
			return nil, err
		}
		row := make(map[string]interface{})
		for i, c := range cols {
			v := vals[i]
			if b, ok := v.([]byte); ok {
				row[c] = string(b)
			} else {
				row[c] = v
			}
		}
		results = append(results, row)
	}
	return results, nil
}

func init() {
	Register("postgres", func() Driver { return &PostgresDriver{} })
}
