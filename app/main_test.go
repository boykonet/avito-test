package main

import (
	apimethods "app/api/methods"
	pkgpostgres "app/pkg/postgres"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
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

func TestBalance(t *testing.T) {
	pool, err := pkgpostgres.NewPool(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(fmt.Errorf("NewPool() error: %w", err))
	}

	defer pool.Close()

	server := CreateNewServer()

	server.MountHandlers(apimethods.New(pool))

    // Create a New Request
    //req, _ := http.NewRequest("GET", "/balance", []byte("{\"id\":1}"))
	req, _ := http.NewRequest("GET", "/balance", nil)

    // Execute Request
    response := executeRequest(req, server)

    // Check the response code
    checkResponseCode(t, http.StatusOK, response.Code)

    // We can use testify/require to assert values, as it is more convenient
    require.Equal(t, "", response.Body.String())
}
