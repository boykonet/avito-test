package api

import (
	apimethods "app/api/methods"
	"github.com/jackc/pgx/v4/pgxpool"
)

type api struct {
	*apimethods.Methods
	pool	*pgxpool.Pool
}

func CreateApi(pgxPool *pgxpool.Pool) Api {
	return &api{
		Methods:	apimethods.New(pgxPool),
		pool:		pgxPool,
	}
}

type Api interface {
	GetBalance(id int) (int, float64, error)
	RefillAndWithdrawMoney(id int, sum float64) (int, float64, error)
	TransferMoney(from, to int, sum float64) (int, float64, int, float64, error)
}
