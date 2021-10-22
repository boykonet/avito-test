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
// 1. Input data:
//		Content-Type: application/json
//		request body: {"id":id}
//		---
//		id - user id
//		id > 0
// 2. Output:
//		Content-Type: application/json
//		response body: {"status":status,"id":id,"balance":balance}
//		---
//		status - response status
//		id - user id
//		balance - user balance
//		---
//		If successful:
//			status = 0, id > 0, balance >= 0.00
//		If data is not a valid:
//			status = 1, id = 0, balance = 0.00
//		If user ID does not exist:
//			status = 2, id = 0, balance = 0.00
//		If server error:
//			status = 3, id = 0, balance = 0.00
//
// ---
// RefillAndWithdrawMoney() method:
// 1. Input data:
//		Content-Type: application/json
//		request body: {"id":id,"sum":sum}
//		---
//		id - user id
//		sum - amount of money to refill or withdraw
//		id > 0, sum != 0
// 2. Output:
//		Content-Type: application/json
//		response body: {"status":status,"id":id,"balance":balance}
//		---
//		status - response status
//		id - user id
//		balance - user balance
//		---
//		If successful:
//			status = 0, id > 0, balance >= 0.00
//		If data is not a valid:
//			status = 1, id = 0, balance = 0.00
//		If insufficient funds:
//			status = 2, id = 0, balance = 0.00
//		If user ID does not exist:
//			status = 4, id = 0, balance = 0.00
//		If server error:
//			status = 3, id = 0, balance = 0.00
// ---
// TransferMoney() method:
// 1. Input data:
//		Content-Type: application/json
//		request body: {"from":from,"to":to,"sum":sum}
//		---
//		from - user id who transfer money
//		to - user id to whom money is transferred
//		sum - amount of money to transfer
//		from > 0, to > 0, sum != 0
// 2. Output
//		Content-Type: application/json
//		response body: {"status":status,"id":id,"balance":balance}
//		---
//		status - response status
//		id - "from" user id
//		balance - "from" user balance
//		---
//		If successful:
//			status = 0, id > 0, balance >= 0.00
//		If data is not a valid:
//			status = 1, id = 0, balance = 0.00
//		If insufficient funds:
//			status = 2, id = 0, balance = 0.00
//		If user ID does not exist:
//			status = 4, id = 0, balance = 0.00
//		If server error:
//			status = 3, id = 0, balance = 0.00

type Server struct {
	Router		*chi.Mux
}

type RequestGetBalance struct {
	ID			int		`json:"id"`
}

type RequestRefillWithdraw struct {
	ID			int		`json:"id"`
	Sum			float64	`json:"sum"`
}

type RequestTransfer struct {
	From		int		`json:"from"`
	To			int		`json:"to"`
	Sum			float64	`json:"sum"`
}

type ResponseToUser struct {
	Status		int		`json:"status"`
	ID			int		`json:"id"`
	Balance		float64	`json:"balance"`
}

func GetBalanceHandler(GetBalance func(int) (int, float64, error)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var request		RequestGetBalance
		var response	ResponseToUser

		w.Header().Set("Content-Type", "application/json")

		err := json.NewDecoder(r.Body).Decode(&request)
		p := r.Header.Get("Content-Type")

		switch {
		case err != nil:
			w.WriteHeader(http.StatusBadRequest)
			response = ResponseToUser{ Status: 1, ID: 0, Balance: 0.00 }
			log.Println("json.NewDecoder(): %w", err)
		case request.ID <= 0:
			w.WriteHeader(http.StatusBadRequest)
			response = ResponseToUser{ Status: 1, ID: 0, Balance: 0.00 }
		case p != "application/json":
			w.WriteHeader(http.StatusBadRequest)
			response = ResponseToUser{ Status: 1, ID: 0, Balance: 0.00 }
		default:
			uid, ub, err := GetBalance(request.ID)
			switch {
			case errors.Is(err, apimethods.UserNotFound):
				w.WriteHeader(http.StatusBadRequest)
				response = ResponseToUser{ Status: 2, ID: 0, Balance: 0.00 }
			default:
				w.WriteHeader(http.StatusOK)
				response = ResponseToUser{ Status: 0, ID: uid, Balance: ub }
			}
		}
		render.JSON(w, r, response)
	}
}

func RefillAndWithdrawHandler(RefillAndWithdrawMoney func(int, float64) (int, float64, error)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var request		RequestRefillWithdraw
		var response	ResponseToUser

		w.Header().Set("Content-Type", "application/json")

		err := json.NewDecoder(r.Body).Decode(&request)
		p := r.Header.Get("Content-Type")

		switch {
		case err != nil:
			w.WriteHeader(http.StatusBadRequest)
			response = ResponseToUser{ Status: 1, ID: 0, Balance: 0.00 }
			log.Println("json.NewDecoder(): %w", err)
		case request.ID <= 0:
			w.WriteHeader(http.StatusBadRequest)
			response = ResponseToUser{ Status: 1, ID: 0, Balance: 0.00 }
		case p != "application/json":
			w.WriteHeader(http.StatusBadRequest)
			response = ResponseToUser{ Status: 1, ID: 0, Balance: 0.00 }
		default:
			uid, ub, err := RefillAndWithdrawMoney(request.ID, request.Sum)
			switch {
			case errors.Is(err, apimethods.InsufficientFunds):
				w.WriteHeader(http.StatusBadRequest)
				response = ResponseToUser{ Status: 2, ID: 0, Balance: 0.00 }
			default:
				w.WriteHeader(http.StatusOK)
				response = ResponseToUser{ Status: 0, ID: uid, Balance: ub }
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

func TransferHandler(TransferMoney func(int, int, float64)(int, float64, error)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var request		RequestTransfer
		var response	ResponseToUser

		w.Header().Set("Content-Type", "application/json")

		err := json.NewDecoder(r.Body).Decode(&request)
		p := r.Header.Get("Content-Type")

		switch {
		case err != nil:
			w.WriteHeader(http.StatusBadRequest)
			response = ResponseToUser{ Status: 1, ID: 0, Balance: 0.00 }
			log.Println("json.NewDecoder(): %w", err)
		case request.From <= 0 || request.To <= 0:
			w.WriteHeader(http.StatusBadRequest)
			response = ResponseToUser{ Status: 1, ID: 0, Balance: 0.00 }
		case p != "application/json":
			w.WriteHeader(http.StatusBadRequest)
			response = ResponseToUser{ Status: 1, ID: 0, Balance: 0.00 }
		default:
			uid, ub, err := TransferMoney(request.From, request.To, request.Sum)
			switch {
			case errors.Is(err, apimethods.InsufficientFunds):
				w.WriteHeader(http.StatusBadRequest)
				response = ResponseToUser{ Status: 2, ID: 0, Balance: 0.00 }
			default:
				w.WriteHeader(http.StatusOK)
				response = ResponseToUser{ Status: 0, ID: uid, Balance: ub }
			}
		}
		render.JSON(w, r, response)
	}
}

func CreateNewServer() *Server {
	s := &Server{ Router: chi.NewRouter() }
	return s
}

func (s *Server) MountHandlers(api *apimethods.Methods) {
	s.Router.Use(middleware.RequestID)
	s.Router.Use(middleware.Logger)

	s.Router.Get("/balance", GetBalanceHandler(api.GetBalance))
	s.Router.Post("/refill", RefillAndWithdrawHandler(api.RefillAndWithdrawMoney))
	s.Router.Post("/withdraw", RefillAndWithdrawHandler(api.RefillAndWithdrawMoney))
	s.Router.Post("/transfer", TransferHandler(api.TransferMoney))
}

func main() {
	pool, err := pkgpostgres.NewPool(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(fmt.Errorf("NewPool: %w", err))
	}

	defer pool.Close()

	api := apimethods.New(pool)

	server := CreateNewServer()

	server.MountHandlers(api)
	//r.Use(middleware.RequestID)
	//r.Use(middleware.Logger)
	//
	//r.Get("/balance", GetBalanceHandler(api.GetBalance))
	//r.Post("/refill", RefillAndWithdrawHandler(api.RefillAndWithdrawMoney))
	//r.Post("/withdraw", RefillAndWithdrawHandler(api.RefillAndWithdrawMoney))
	//r.Post("/transfer", TransferHandler(api.TransferMoney))

	log.Fatal(http.ListenAndServe(":8080", server.Router))
	// c := make(chan os.Signal, 1)
	// signal.Notify(c, os.Interrupt)
	// <-c
}
