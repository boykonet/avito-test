package methods

import (
	"github.com/jackc/pgx/v4/pgxpool"
)

type Methods struct {
	pool	*pgxpool.Pool
}

func New(pgxPool *pgxpool.Pool) *Methods {
	return &Methods{ pool: pgxPool }
}
