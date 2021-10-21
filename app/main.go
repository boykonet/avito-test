package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	apimethods "app/api/methods"
	pkgpostgres "app/pkg/postgres"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
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
//			status = 0, uid > 0, ub > 0.00
//		If data is not a valid:
//			status = 1, uid = 0, ub = 0
//		If user ID does not exist:
//			status = 2, uid = 0, ub = 0.00
//
// ---
// RefillAndWithdrawMoney() method:
// 1. Input data:
//		Content-Type: application/json
//		request body: {"id":uid,"sum":usum}, uid > 0, usum

type RequestGetBalance struct {
	ID			int		`json:"id"`
}

type ResponseGetBalance struct {
	Status		int		`json:"status"`
	ID			int		`json:"id"`
	Balance		float64	`json:"balance"`
}

type RequestRefillWithdraw struct {
	ID			int		`json:"id"`
	Sum			float64	`json:"sum"`
}

type ResponseRefillWithdraw struct {
	Status		int		`json:"status"`
	ID			int		`json:"id"`
	Balance		float64	`json:"balance"`
}

func GetBalanceHandler(GetBalance func(int) (int, float64, error)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var request		RequestGetBalance
		var response	ResponseGetBalance

		w.Header().Set("Content-Type", "application/json")

		err := json.NewDecoder(r.Body).Decode(&request)
		p := r.Header.Get("Content-Type")

		switch {
		case err != nil:
			w.WriteHeader(http.StatusBadRequest)
			response = ResponseGetBalance{ Status: 1, ID: 0, Balance: 0.00 }
			log.Println("json.NewDecoder(): %w", err)
		case request.ID <= 0:
			w.WriteHeader(http.StatusBadRequest)
			response = ResponseGetBalance{ Status: 1, ID: 0, Balance: 0.00 }
		case p != "application/json":
			w.WriteHeader(http.StatusBadRequest)
			response = ResponseGetBalance{ Status: 1, ID: 0, Balance: 0.00 }
		default:
			uid, ub, err := GetBalance(request.ID)
			switch {
			case errors.Is(err, apimethods.UserNotFound):
				w.WriteHeader(http.StatusBadRequest)
				response = ResponseGetBalance{ Status: 2, ID: 0, Balance: 0.00 }
			default:
				w.WriteHeader(http.StatusOK)
				response = ResponseGetBalance{ Status: 0, ID: uid, Balance: ub }
			}
		}
		render.JSON(w, r, response)
	}
}

func RefillAndWithdrawHandler(RefillAndWithdrawMoney func(int, float64) (int, float64, error)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var request		RequestRefillWithdraw
		var response	ResponseRefillWithdraw

		w.Header().Set("Content-Type", "application/json")

		err := json.NewDecoder(r.Body).Decode(&request)
		p := r.Header.Get("Content-Type")

		switch {
		case err != nil:
			w.WriteHeader(http.StatusBadRequest)
			response = ResponseRefillWithdraw{ Status: 1, ID: 0, Balance: 0.00 }
			log.Println("json.NewDecoder(): %w", err)
		case request.ID <= 0:
			w.WriteHeader(http.StatusBadRequest)
			response = ResponseRefillWithdraw{ Status: 1, ID: 0, Balance: 0.00 }
		case p != "application/json":
			w.WriteHeader(http.StatusBadRequest)
			response = ResponseRefillWithdraw{ Status: 1, ID: 0, Balance: 0.00 }
		default:
			uid, ub, err := RefillAndWithdrawMoney(request.ID, request.Sum)
			switch {
			case errors.Is(err, apimethods.InsufficientFunds):
				w.WriteHeader(http.StatusBadRequest)
				response = ResponseRefillWithdraw{ Status: 2, ID: 0, Balance: 0.00 }
			default:
				w.WriteHeader(http.StatusOK)
				response = ResponseRefillWithdraw{ Status: 0, ID: uid, Balance: ub }
			}
		}
		render.JSON(w, r, response)
		//u.ID, u.Balance, err = RefillAndWithdrawMoney(u.ID, u.Balance)
		//w.WriteHeader(http.StatusCreated)
		//w.Header().Set("Content-Type:", "application/json")
		//b, err := json.Marshal(u)
		//w.Write(b)
	}
}

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

	r.Get("/balance", GetBalanceHandler(api.GetBalance))
	r.Post("/refill", RefillAndWithdrawHandler(api.RefillAndWithdrawMoney))
	r.Post("/withdraw", RefillAndWithdrawHandler(api.RefillAndWithdrawMoney))
	//r.Post("/transfer", TransferHandler(api))

	log.Fatal(http.ListenAndServe(":8080", r))
	// c := make(chan os.Signal, 1)
	// signal.Notify(c, os.Interrupt)
	// <-c
}
