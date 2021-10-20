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
// 1. Input:
//		Content-Type: application/json
//		request body: {"id":uid}
//		---
//		uid > 0
// 2. Output:
//		Content-Type: application/json
//		response body: {"status":status,"message":message,:"id":uid,"balance":ub}
//		---
//		status - response status, message - error message, uid - user id, ub - user balance
//		If successful:
//			status = 200, message = "OK", uid > 0, ub > 0.00
//		Otherwise:
//			uid = 0, ub = 0.00
//		status = 500, message = "Internal Server Error"
// ---
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

func GetBalanceHandler(GetBalance func(int) (int, float64, error)) func(w http.ResponseWriter, r *http.Request) {
	var out Output
	var u User

	return func(w http.ResponseWriter, r *http.Request) {
		err := json.NewDecoder(r.Body).Decode(&u)
		if err != nil {
			log.Fatal(err)
		}
		uid, ub, err := GetBalance(u.ID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Println(err)
			out = Output{
				Status:			500,
				Message:		"Internal Server Error",
				UserID:			0,
				UserBalance:	0.00,
			}
		} else {
			w.WriteHeader(http.StatusOK)
			out = Output{
				Status:			200,
				Message:		"OK",
				UserID:			uid,
				UserBalance:	ub,
			}
		}
		w.Header().Set("Content-Type:", "application/json")
		b, err := json.Marshal(out)
		_, err = w.Write(b)
		if err != nil {
			log.Fatal(err)
		}
	}
}

//func RefillAndWithdrawHandler(RefillAndWithdrawMoney func(int, float64)) func(w http.ResponseWriter, r *http.Request) {
//	return func(w http.ResponseWriter, r *http.Request) {
//		var user User
//
//		err := json.NewDecoder(r.Body).Decode(&user)
//		if err != nil {
//			log.Fatal(err)
//		}
//		user.ID, user.Balance, err = RefillAndWithdrawMoney(user.ID, user.Balance)
//		w.WriteHeader(http.StatusCreated)
//		w.Header().Set("Content-Type:", "application/json")
//		b, err := json.Marshal(user)
//		w.Write(b)
//	}
//}
//
//func TransferHandler(api *apimethods.Methods) func(w http.ResponseWriter, r *http.Request) {
//	return func(w http.ResponseWriter, r *http.Request) {
//		api.TransferMoney(1, 2, 33.33)
//		w.WriteHeader(http.StatusCreated)
//		w.Header().Set("Content-Type:", "application/json")
//	}
//}

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

	//r.Post("/refill", RefillAndWithdrawHandler(api))
	//r.Post("/withdraw", RefillAndWithdrawHandler(api))
	//r.Post("/transfer", TransferHandler(api))
	r.Get("/balance", GetBalanceHandler(api.GetBalance))

	log.Fatal(http.ListenAndServe(":8080", r))
	// c := make(chan os.Signal, 1)
	// signal.Notify(c, os.Interrupt)
	// <-c
}
