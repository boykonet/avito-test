package main

import (
	apimethods "app/api/methods"
	pkgpostgres "app/pkg/postgres"
	handlers "app/handlers"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
	"os"
)

type Server struct {
	Router		*chi.Mux
}

func CreateNewServer() *Server {
	s := &Server{ Router: chi.NewRouter() }
	return s
}

func (s *Server) MountHandlers(api *apimethods.Methods) {
	s.Router.Use(middleware.RequestID)
	s.Router.Use(middleware.Logger)

	s.Router.Get("/balance", handlers.GetBalanceHandler(api.GetBalance))
	s.Router.Post("/refill", handlers.RefillAndWithdrawHandler(api.RefillAndWithdrawMoney))
	s.Router.Post("/withdraw", handlers.RefillAndWithdrawHandler(api.RefillAndWithdrawMoney))
	s.Router.Post("/transfer", handlers.TransferHandler(api.TransferMoney))
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

	log.Fatal(http.ListenAndServe(":8080", server.Router))
}
