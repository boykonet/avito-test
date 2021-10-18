package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
)

func NewPool(conn string) (*pgxpool.Pool, error) {
	poolConfig, err := pgxpool.ParseConfig(conn)
	if err != nil {
		return nil, fmt.Errorf("ParseConfig: %w", err)
	}

	pool, err := pgxpool.ConnectConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, fmt.Errorf("ConnectConfig: %w", err)
	}
	return pool, nil
}
