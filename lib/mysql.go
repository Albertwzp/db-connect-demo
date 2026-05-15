package lib

import (
	"context"
	"database/sql"
	"errors"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type MySQLDriver struct{
	db *sql.DB
}

func (m *MySQLDriver) Open(dsn string) error {
	if dsn == "" {
		return errors.New("empty dsn")
	}
	db, err := sql.Open("mysql", dsn)
	if err != nil { return err }
	db.SetMaxOpenConns(50)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(time.Minute*5)
	m.db = db
	return nil
}

func (m *MySQLDriver) Close() error {
	if m.db != nil { return m.db.Close() }
	return nil
}

func (m *MySQLDriver) RunOp(ctx context.Context, query string) (time.Duration, error) {
	if m.db == nil { return 0, errors.New("db not open") }
	start := time.Now()
	_, err := m.db.ExecContext(ctx, query)
	return time.Since(start), err
}

func (m *MySQLDriver) HealthCheck(ctx context.Context) error {
	if m.db == nil { return errors.New("db not open") }
	return m.db.PingContext(ctx)
}

func (m *MySQLDriver) Query(ctx context.Context, query string) ([]map[string]interface{}, error) {
	if m.db == nil { return nil, errors.New("db not open") }
	rows, err := m.db.QueryContext(ctx, query)
	if err != nil { return nil, err }
	defer rows.Close()
	cols, err := rows.Columns()
	if err != nil { return nil, err }
	results := []map[string]interface{}{}
	for rows.Next() {
		vals := make([]interface{}, len(cols))
		ptrs := make([]interface{}, len(cols))
		for i := range vals { ptrs[i] = &vals[i] }
		if err := rows.Scan(ptrs...); err != nil { return nil, err }
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
	Register("mysql", func() Driver { return &MySQLDriver{} })
}
