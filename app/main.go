package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	apimethods "app/api/methods"
	pkgpostgres "app/pkg/postgres"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)
// GetBalance() method:
// 1. Input data:
// Content-Type: application/json
// {"id":uid}, uid > 0
// 2. Output data:
// Content-Type: application/json
// {"id":uid,"balance":ub}, uid - user id (uid > 0), ub - user balance (ub > 0)
// RefillAndWithdrawMoney() method:
// 1. Input data:
// Content-Type: application/json
// {"id":uid,"sum":usum}, uid > 0, usum

type User struct {
	ID			int		`json:"id"`
	Balance		float64	`json:"balance"`
}

type Output struct {
	Status		int		`json:"status"`
	Message		string	`json:"message"`
	UserID		int		`json:"id"`
	UserBalance	float64	`json:"balance"`
}

func RefillAndWithdrawHandler(api *apimethods.Methods) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var user User

		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			log.Fatal(err)
		}
		user.ID, user.Balance, err = api.RefillAndWithdrawMoney(user.ID, user.Balance)
		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type:", "application/json")
		b, err := json.Marshal(user)
		w.Write(b)
	}
}

func TransferHandler(api *apimethods.Methods) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		api.TransferMoney(1, 2, 33.33)
		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type:", "application/json")
	}
}

func GetBalanceHandler(api *apimethods.Methods) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		api.GetBalance(1)
		_, _ = w.Write([]byte("Hello from get balance\n"))
	}
}

func main() {
	pool, err := pkgpostgres.NewPool(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(fmt.Errorf("NewPool: %w", err))
	}

	defer pool.Close()

	api := apimethods.New(pool)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)

	r.Post("/refill", RefillAndWithdrawHandler(api))
	r.Post("/withdraw", RefillAndWithdrawHandler(api))
	r.Post("/transfer", TransferHandler(api))
	r.Get("/balance", GetBalanceHandler(api))

	log.Fatal(http.ListenAndServe(":8080", r))
	// c := make(chan os.Signal, 1)
	// signal.Notify(c, os.Interrupt)
	// <-c
}
