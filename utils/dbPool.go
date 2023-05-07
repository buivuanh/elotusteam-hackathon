package utils

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

const maxAttempt = 20

func NewConnectionPool(ctx context.Context, dbConnectString string) (*pgxpool.Pool, error) {
	var dbPool *pgxpool.Pool
	if err := Do(func(attempt int) (retry bool, err error) {
		defer time.Sleep(1 * time.Second)
		dbPool, err = pgxpool.New(ctx, dbConnectString)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to connect to database: %v\n", err)
			return attempt < maxAttempt, err
		}

		// query to check connection
		var greeting string
		if err = dbPool.QueryRow(ctx, "select 'Hello, world!'").Scan(&greeting); err != nil {
			fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
			return attempt < maxAttempt, err
		}

		return false, nil
	}); err != nil {
		return nil, fmt.Errorf("Unable to create connection pool: %v\n", err)
	}

	return dbPool, nil
}
