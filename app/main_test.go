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
	"strconv"
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

	sum := 5.0
	id, balance, _ :=  api.GetBalance(2)
	expBalance := fmt.Sprintf(`{"status":0,"id":%v,"balance":%.2f}`, id, math.Round((sum + balance) * 100) / 100)

	// Correct request
	request := fmt.Sprintf(`{"id":%v,"sum":%.2f}`, id, sum)

	checkMethods(t, server,
		request,
		`POST`,
		`/refill`,
		`application/json`,
		http.StatusOK,
		expBalance)

	// Negative user id
	request = fmt.Sprintf(`{"id":%v,"sum":%.2f}`, -id, sum)

	checkMethods(t, server,
		request,
		`POST`,
		`/refill`,
		`application/json`,
		http.StatusBadRequest,
		`{"status":1,"id":0,"balance":0}`)

	// User id is not a number (string)
	request = fmt.Sprintf(`{"id":"%v","sum":%.2f}`, id, sum)

	checkMethods(t, server,
		request,
		`POST`,
		`/refill`,
		`application/json`,
		http.StatusBadRequest,
		`{"status":1,"id":0,"balance":0}`)

	// Sum is not a float (string)
	request = fmt.Sprintf(`{"id":%v,"sum":"%.2f"}`, id, sum)

	checkMethods(t, server,
		request,
		`POST`,
		`/refill`,
		`application/json`,
		http.StatusBadRequest,
		`{"status":1,"id":0,"balance":0}`)

	// Incorrect content type (not application/json)
	request = fmt.Sprintf(`{"id":%v,"sum":%.2f}`, id, sum)

	checkMethods(t, server,
		request,
		`POST`,
		`/refill`,
		`blablabla`,
		http.StatusBadRequest,
		`{"status":1,"id":0,"balance":0}`)

	// Incorrect body
	request = `{"beleberda":1}`

	checkMethods(t, server,
		request,
		`POST`,
		`/refill`,
		`application/json`,
		http.StatusBadRequest,
		`{"status":1,"id":0,"balance":0}`)


	// User id doesn't exist
	request = fmt.Sprintf(`{"id":%v,"sum":%.2f}`, id + 100, sum)

	checkMethods(t, server,
		request,
		`POST`,
		`/refill`,
		`application/json`,
		http.StatusBadRequest,
		`{"status":2,"id":0,"balance":0}`)

	// Method not allowed
	request = fmt.Sprintf(`{"id":%v,"sum":%.2f}`, id, sum)

	checkMethods(t, server,
		request,
		`GET`,
		`/refill`,
		`application/json`,
		http.StatusMethodNotAllowed,
		``)

	//-------------------- WITHDRAW --------------------

	id, sum, _ =  api.GetBalance(2)
	expBalance = fmt.Sprintf(`{"status":0,"id":%v,"balance":%.2f}`, id, math.Round((sum - 5.0) * 100) / 100)

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
		`{"status":3,"id":0,"balance":0}`)

	// Method not allowed
	checkMethods(t, server,
		`{"id":2,"sum":5}`,
		`GET`,
		`/withdraw`,
		`application/json`,
		http.StatusMethodNotAllowed,
		``)

	//-------------------- TRANSFER --------------------

	from_id, from_balance, _ :=  api.GetBalance(2)
	to_id, to_balance, _ :=  api.GetBalance(3)
	sum = 5.

	expBalance = fmt.Sprintf(`{"status":0,"from_id":%v,"from_balance":%.2f,"to_id":%v,"to_balance":%.2f}`,
		from_id, math.Round((from_balance - sum) * 100) / 100,
		to_id, math.Round((to_balance + sum) * 100) / 100 )

	// Correct request
	request := fmt.Sprintf(`{"from":%v,"to":%v,"sum":%.2f}`, from_id, to_id, sum)

	checkMethods(t, server,
		request,
		`POST`,
		`/transfer`,
		`application/json`,
		http.StatusOK,
		expBalance)

	// User id is not a number (string)
	request = fmt.Sprintf(`{"from":"%v","to":%v,"sum":%.2f}`,
		strconv.Itoa(from_id), to_id, sum)

	checkMethods(t, server,
		request,
		`POST`,
		`/transfer`,
		`application/json`,
		http.StatusBadRequest,
		`{"status":1,"from_id":0,"from_balance":0,"to_id":0,"to_balance":0}`)

	request = fmt.Sprintf(`{"from":%v,"to":"%v","sum":%.2f}`,
		from_id, strconv.Itoa(to_id), sum)

	checkMethods(t, server,
		request,
		`POST`,
		`/transfer`,
		`application/json`,
		http.StatusBadRequest,
		`{"status":1,"from_id":0,"from_balance":0,"to_id":0,"to_balance":0}`)

	// Sum is not a float (string)
	request = fmt.Sprintf(`{"from":%v,"to":%v,"sum":"%v"}`,
		from_id, to_id, strconv.FormatFloat(sum, 'f', 2, 64))

	checkMethods(t, server,
		request,
		`POST`,
		`/transfer`,
		`application/json`,
		http.StatusBadRequest,
		`{"status":1,"from_id":0,"from_balance":0,"to_id":0,"to_balance":0}`)

	// Sum is negative
	request = fmt.Sprintf(`{"from":%v,"to":%v,"sum":%.2f}`, from_id, to_id, -sum)

	checkMethods(t, server,
		request,
		`POST`,
		`/transfer`,
		`application/json`,
		http.StatusBadRequest,
		`{"status":1,"from_id":0,"from_balance":0,"to_id":0,"to_balance":0}`)

	// Incorrect content type (not application/json)
	request = fmt.Sprintf(`{"from":%v,"to":%v,"sum":%.2f}`, from_id, to_id, sum)

	checkMethods(t, server,
		request,
		`POST`,
		`/transfer`,
		`blablabla`,
		http.StatusBadRequest,
		`{"status":1,"from_id":0,"from_balance":0,"to_id":0,"to_balance":0}`)

	// Incorrect body
	checkMethods(t, server,
		`{"beleberda":1}`,
		`POST`,
		`/transfer`,
		`application/json`,
		http.StatusBadRequest,
		`{"status":1,"from_id":0,"from_balance":0,"to_id":0,"to_balance":0}`)

	// Method not allowed
	request = fmt.Sprintf(`{"from":%v,"to":%v,"sum":%.2f}`, from_id, to_id, sum)

	checkMethods(t, server,
		request,
		`GET`,
		`/transfer`,
		`application/json`,
		http.StatusMethodNotAllowed,
		``)
}
