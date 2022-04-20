package db

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/lib/pq"
)

func SetPgConn(ctx context.Context, dbUrl string) (*sql.DB, error) {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		return nil, err
	}

	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}

	return db, nil
}
