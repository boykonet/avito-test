package main

import (
	apimethods "app/api/methods"
	pkgpostgres "app/pkg/postgres"
	"bytes"
	"fmt"
	"github.com/stretchr/testify/require"
	"log"
	"math"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func executeRequest(req *http.Request, s *Server) *httptest.ResponseRecorder {
    rr := httptest.NewRecorder()
    s.Router.ServeHTTP(rr, req)

    return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
    if expected != actual {
        t.Errorf("Expected response code %d. Got %d\n", expected, actual)
    }
}

func checkMethods(t *testing.T, server *Server, body, method, url, headerValue string, expectedStatus int, expectedBody string) {
	var actualBody = ""
	var jsonStr = []byte(body)

	req, _ := http.NewRequest(method, url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", headerValue)

	response := executeRequest(req, server)

	checkResponseCode(t, expectedStatus, response.Code)

	lenResponseBody := len(response.Body.String())

	if lenResponseBody > 0 {
		actualBody = response.Body.String()[:len(response.Body.String()) - 1]
	}

	require.Equal(t, expectedBody, actualBody)
}

func TestBalance(t *testing.T) {
	pool, err := pkgpostgres.NewPool(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(fmt.Errorf("NewPool() error: %w", err))
	}

	defer pool.Close()

	server := CreateNewServer()

	server.MountHandlers(apimethods.New(pool))

	//-------------------- GET BALANCE --------------------

	// Correct request
	checkMethods(t, server,
		`{"id":1}`,
		`GET`,
		`/balance`,
		`application/json`,
		http.StatusOK,
		`{"status":0,"id":1,"balance":56.99}`)

	// Negative user id
	checkMethods(t, server,
		`{"id":-1}`,
		`GET`,
		`/balance`,
		`application/json`,
		http.StatusBadRequest,
		`{"status":1,"id":0,"balance":0}`)

	// User id is not a number (string)
	checkMethods(t, server,
		`{"id":"1"}`,
		`GET`,
		`/balance`,
		`application/json`,
		http.StatusBadRequest,
		`{"status":1,"id":0,"balance":0}`)

	// User id doesn't exist
	checkMethods(t, server,
		`{"id":6}`,
		`GET`,
		`/balance`,
		`application/json`,
		http.StatusBadRequest,
		`{"status":2,"id":0,"balance":0}`)

	// Incorrect content type (not application/json)
	checkMethods(t, server,
		`{"id":1}`,
		`GET`,
		`/balance`,
		`blablabla`,
		http.StatusBadRequest,
		`{"status":1,"id":0,"balance":0}`)

	// Incorrect body
	checkMethods(t, server,
		`{"beleberda":1}`,
		`GET`,
		`/balance`,
		`blablabla`,
		http.StatusBadRequest,
		`{"status":1,"id":0,"balance":0}`)

	// Method not allowed
	checkMethods(t, server,
		`{"id":1}`,
		`POST`,
		`/balance`,
		`application/json`,
		http.StatusMethodNotAllowed,
		``)

	//-------------------- REFILL --------------------

	api := apimethods.New(pool)

	_, sum, _ :=  api.GetBalance(2)
	expBalance := fmt.Sprintf(`{"status":0,"id":2,"balance":%.2f}`, math.Round((5.0 + sum) * 100) / 100)

	// Correct request
	checkMethods(t, server,
		`{"id":2,"sum":5}`,
		`POST`,
		`/refill`,
		`application/json`,
		http.StatusOK,
		expBalance)

	// Check incorrect user id
	checkMethods(t, server,
		`{"id":-1,"sum":5}`,
		`POST`,
		`/refill`,
		`application/json`,
		http.StatusBadRequest,
		`{"status":1,"id":0,"balance":0}`)

	// User id is not a number (string)
	checkMethods(t, server,
		`{"id":"2","sum":5}`,
		`POST`,
		`/refill`,
		`application/json`,
		http.StatusBadRequest,
		`{"status":1,"id":0,"balance":0}`)

	// Sum is not a float (string)
	checkMethods(t, server,
		`{"id":2,"sum":"5"}`,
		`POST`,
		`/refill`,
		`application/json`,
		http.StatusBadRequest,
		`{"status":1,"id":0,"balance":0}`)

	// Incorrect content type (not application/json)
	checkMethods(t, server,
		`{"id":2,"sum":5}`,
		`POST`,
		`/refill`,
		`blablabla`,
		http.StatusBadRequest,
		`{"status":1,"id":0,"balance":0}`)

	// Incorrect body
	checkMethods(t, server,
		`{"beleberda":1}`,
		`POST`,
		`/refill`,
		`application/json`,
		http.StatusBadRequest,
		`{"status":1,"id":0,"balance":0}`)


	// User id doesn't exist
	checkMethods(t, server,
		`{"id":6,"sum":5}`,
		`POST`,
		`/refill`,
		`application/json`,
		http.StatusBadRequest,
		`{"status":4,"id":0,"balance":0}`)

	// Method not allowed
	checkMethods(t, server,
		`{"id":2,"sum":5}`,
		`GET`,
		`/refill`,
		`application/json`,
		http.StatusMethodNotAllowed,
		``)

	//-------------------- WITHDRAW --------------------

	_, sum, _ =  api.GetBalance(2)
	expBalance = fmt.Sprintf(`{"status":0,"id":2,"balance":%.2f}`, math.Round((sum - 5.0) * 100) / 100)

	// Correct request
	checkMethods(t, server,
		`{"id":2,"sum":-5}`,
		`POST`,
		`/withdraw`,
		`application/json`,
		http.StatusOK,
		expBalance)

	// Correct request with sum = 0
	checkMethods(t, server,
		`{"id":2,"sum":0}`,
		`POST`,
		`/withdraw`,
		`application/json`,
		http.StatusOK,
		expBalance)

	// Insufficient funds
	// User with ID=5 has balance=0.00
	checkMethods(t, server,
		`{"id":5,"sum":-5}`,
		`POST`,
		`/withdraw`,
		`application/json`,
		http.StatusBadRequest,
		`{"status":2,"id":0,"balance":0}`)

	// Method not allowed
	checkMethods(t, server,
		`{"id":2,"sum":5}`,
		`GET`,
		`/withdraw`,
		`application/json`,
		http.StatusMethodNotAllowed,
		``)

	//-------------------- TRANSFER --------------------

	_, sum, _ =  api.GetBalance(2)
	expBalance = fmt.Sprintf(`{"status":0,"id":2,"balance":%.2f}`, math.Round((sum - 5.0) * 100) / 100)

	// Correct request
	checkMethods(t, server,
		`{"from":2,"to":3,"sum":5}`,
		`POST`,
		`/transfer`,
		`application/json`,
		http.StatusOK,
		expBalance)

	// User id is not a number (string)
	checkMethods(t, server,
		`{"from":"2","to":3,"sum":5}`,
		`POST`,
		`/transfer`,
		`application/json`,
		http.StatusBadRequest,
		`{"status":1,"id":0,"balance":0}`)

	checkMethods(t, server,
		`{"from":2,"to":"3","sum":5}`,
		`POST`,
		`/transfer`,
		`application/json`,
		http.StatusBadRequest,
		`{"status":1,"id":0,"balance":0}`)

	// Sum is not a float (string)
	checkMethods(t, server,
		`{"from":2,"to":3,"sum":"5"}`,
		`POST`,
		`/transfer`,
		`application/json`,
		http.StatusBadRequest,
		`{"status":1,"id":0,"balance":0}`)

	// Sum is negative
	checkMethods(t, server,
		`{"from":2,"to":3,"sum":-5}`,
		`POST`,
		`/transfer`,
		`application/json`,
		http.StatusBadRequest,
		`{"status":1,"id":0,"balance":0}`)

	// Incorrect content type (not application/json)
	checkMethods(t, server,
		`{"from":2,"to":3,"sum":5}`,
		`POST`,
		`/transfer`,
		`blablabla`,
		http.StatusBadRequest,
		`{"status":1,"id":0,"balance":0}`)

	// Incorrect body
	checkMethods(t, server,
		`{"beleberda":1}`,
		`POST`,
		`/transfer`,
		`application/json`,
		http.StatusBadRequest,
		`{"status":1,"id":0,"balance":0}`)

	// Method not allowed
	checkMethods(t, server,
		`{"from":2,"to":3,"sum":5}`,
		`GET`,
		`/transfer`,
		`application/json`,
		http.StatusMethodNotAllowed,
		``)
}
