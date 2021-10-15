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
	TransferAndWithdrawMoney(id int, sum float64) (error)
	TransferMoney(first_id, second_id int, sum float64) (error)
	GetBalance(id int) (float64, error)
}
