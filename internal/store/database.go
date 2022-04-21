package store

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

type Row interface {
	Scan(dest ...interface{}) error
}

func NewDatabase(dbSourceName string, timeout time.Duration) (*sql.DB, error) {
	ticker := time.NewTicker(1000 * time.Millisecond)
	defer ticker.Stop()

	timeoutExceeded := time.After(timeout)

	for {
		select {
		case <-timeoutExceeded:
			return nil, fmt.Errorf("db connection failed after %s timeout", timeout)
		case <-ticker.C:
			db, err := sql.Open("postgres", dbSourceName)
			if err != nil {
				continue
			}

			if err := db.Ping(); err != nil {
				continue
			}

			if err == nil {
				return db, nil
			}
		}
	}
}
