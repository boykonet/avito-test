package handler

import (
	apimethods "app/api/methods"
	"encoding/json"
	"errors"
	"github.com/go-chi/render"
	"log"
	"math"
	"net/http"
)

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

type ResponseTransfer struct {
	Status		int		`json:"status"`
	FromID		int		`json:"from_id"`
	FromBalance	float64	`json:"from_balance"`
	ToID		int		`json:"to_id"`
	ToBalance	float64	`json:"to_balance"`
}

// GetBalanceHandler method:
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
func GetBalanceHandler(GetBalance func(int) (int, float64, error)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var request		RequestGetBalance
		var response	ResponseToUser

		w.Header().Set("Content-Type", "application/json")

		err := json.NewDecoder(r.Body).Decode(&request)
		p := r.Header.Get("Content-Type")

		switch {
		case err != nil || p != "application/json":
			w.WriteHeader(http.StatusBadRequest)
			response = ResponseToUser{ Status: 1, ID: 0, Balance: 0.00 }
		default:
			uid, ub, err := GetBalance(request.ID)
			switch {
			case err == nil:
				ub = math.Round(ub * 100) / 100
				w.WriteHeader(http.StatusOK)
				response = ResponseToUser{ Status: 0, ID: uid, Balance: ub }
			case errors.Is(err, apimethods.WrongData):
				w.WriteHeader(http.StatusBadRequest)
				response = ResponseToUser{ Status: 1, ID: 0, Balance: 0.00 }
			case errors.Is(err, apimethods.UserNotFound):
				w.WriteHeader(http.StatusBadRequest)
				response = ResponseToUser{ Status: 2, ID: 0, Balance: 0.00 }
			default:
				w.WriteHeader(http.StatusInternalServerError)
				response = ResponseToUser{ Status: 4, ID: 0, Balance: 0 }
				log.Println(err)
			}
		}
		render.JSON(w, r, response)
	}
}

// RefillAndWithdrawHandler method:
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
func RefillAndWithdrawHandler(RefillAndWithdrawMoney func(int, float64) (int, float64, error)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var request		RequestRefillWithdraw
		var response	ResponseToUser

		w.Header().Set("Content-Type", "application/json")

		err := json.NewDecoder(r.Body).Decode(&request)
		p := r.Header.Get("Content-Type")

		switch {
		case err != nil || p != "application/json":
			w.WriteHeader(http.StatusBadRequest)
			response = ResponseToUser{ Status: 1, ID: 0, Balance: 0.00 }
		default:
			uid, ub, err := RefillAndWithdrawMoney(request.ID, request.Sum)
			switch {
			case err == nil:
				ub = math.Round(ub * 100) / 100
				w.WriteHeader(http.StatusOK)
				response = ResponseToUser{ Status: 0, ID: uid, Balance: ub }
			case errors.Is(err, apimethods.WrongData):
				w.WriteHeader(http.StatusBadRequest)
				response = ResponseToUser{ Status: 1, ID: 0, Balance: 0.00 }
			case errors.Is(err, apimethods.UserNotFound):
				w.WriteHeader(http.StatusBadRequest)
				response = ResponseToUser{ Status: 2, ID: 0, Balance: 0.00 }
			case errors.Is(err, apimethods.InsufficientFunds):
				w.WriteHeader(http.StatusBadRequest)
				response = ResponseToUser{ Status: 3, ID: 0, Balance: 0.00 }
			default:
				w.WriteHeader(http.StatusInternalServerError)
				response = ResponseToUser{ Status: 4, ID: 0, Balance: 0.00 }
				log.Println(err)
			}
		}
		render.JSON(w, r, response)
	}
}

// TransferHandler method:
// 1. Input data:
//		Content-Type: application/json
//		request body: {"from":from,"to":to,"sum":sum}
//		---
//		from - user id who transfer money
//		to - user id to whom money is transferred
//		sum - amount of money to transfer
//		from > 0, to > 0, sum > 0
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
func TransferHandler(TransferMoney func(int, int, float64)(int, float64, int, float64, error)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var request		RequestTransfer
		var response	ResponseTransfer

		w.Header().Set("Content-Type", "application/json")

		err := json.NewDecoder(r.Body).Decode(&request)
		p := r.Header.Get("Content-Type")

		switch {
		case err != nil || p != "application/json":
			w.WriteHeader(http.StatusBadRequest)
			response = ResponseTransfer{ Status: 1, FromID: 0, FromBalance: 0.00, ToID: 0, ToBalance: 0.00 }
		default:
			from_id, from_balance, to_id, to_balance, err := TransferMoney(request.From, request.To, request.Sum)
			switch {
			case err == nil:
				from_balance = math.Round(from_balance * 100) / 100
				to_balance = math.Round(to_balance * 100) / 100
				w.WriteHeader(http.StatusOK)
				response = ResponseTransfer{ Status: 0, FromID: from_id, FromBalance: from_balance, ToID: to_id, ToBalance: to_balance }
			case errors.Is(err, apimethods.WrongData):
				w.WriteHeader(http.StatusBadRequest)
				response = ResponseTransfer{ Status: 1, FromID: 0, FromBalance: 0.00, ToID: 0, ToBalance: 0.00 }
			case errors.Is(err, apimethods.UserNotFound):
				w.WriteHeader(http.StatusBadRequest)
				response = ResponseTransfer{ Status: 2, FromID: 0, FromBalance: 0.00, ToID: 0, ToBalance: 0.00 }
			case errors.Is(err, apimethods.InsufficientFunds):
				w.WriteHeader(http.StatusBadRequest)
				response = ResponseTransfer{ Status: 3, FromID: 0, FromBalance: 0.00, ToID: 0, ToBalance: 0.00 }
			default:
				w.WriteHeader(http.StatusInternalServerError)
				response = ResponseTransfer{ Status: 4, FromID: 0, FromBalance: 0.00, ToID: 0, ToBalance: 0.00 }
				log.Println(err)
			}
		}
		render.JSON(w, r, response)
	}
}
