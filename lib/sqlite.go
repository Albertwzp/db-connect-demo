package lib

import (
	"context"
	"database/sql"
	"errors"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type SQLiteDriver struct{
	db *sql.DB
}

func (s *SQLiteDriver) Open(dsn string) error {
	if dsn == "" { return errors.New("empty dsn") }
	// For sqlite, dsn is a file path ("file::memory:?cache=shared")
	db, err := sql.Open("sqlite3", dsn)
	if err != nil { return err }
	db.SetMaxOpenConns(1)
	s.db = db
	return nil
}

func (s *SQLiteDriver) Close() error { if s.db!=nil { return s.db.Close() }; return nil }

func (s *SQLiteDriver) RunOp(ctx context.Context, query string) (time.Duration, error) {
	if s.db==nil { return 0, errors.New("db not open") }
	start := time.Now()
	_, err := s.db.ExecContext(ctx, query)
	return time.Since(start), err
}

func (s *SQLiteDriver) HealthCheck(ctx context.Context) error { if s.db==nil { return errors.New("db not open") }; return s.db.PingContext(ctx) }

func (s *SQLiteDriver) Query(ctx context.Context, query string) ([]map[string]interface{}, error) {
	if s.db==nil { return nil, errors.New("db not open") }
	rows, err := s.db.QueryContext(ctx, query)
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

func init() { Register("sqlite", func() Driver { return &SQLiteDriver{} }) }
