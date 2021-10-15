package main

import (
	"fmt"
	"log"
	//"net/http"
	"os"

	// "encoding/json"
	pkgpostgres "app/pkg/postgres"
	apimethods "app/api/methods"
)

type User struct {
	ID			int		`json:"id"`
	Balance		float64	`json:"balance"`
}

func main() {
	var err error

	pool, err := pkgpostgres.NewPool(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
		return
	}

	defer pool.Close()

	api := apimethods.New(pool)

	balance, err := api.GetBalance(1)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("balance first user:", balance)

	balance, _ = api.GetBalance(2)
	fmt.Println("balance second user:", balance)

	err = api.TransferMoney(1, 2, 5.00)
	if err != nil {
		log.Fatal(err)
	}
	balance, _ = api.GetBalance(1)
	fmt.Println("balance first user:", balance)
	balance, _ = api.GetBalance(2)
	fmt.Println("balance second user:", balance)

	fmt.Println("TRANSFER")
	_ = api.TransferAndWithdrawMoney(1, 5.00)

	balance, _ = api.GetBalance(1)
	fmt.Println("balance first user:", balance)

	_ = api.TransferAndWithdrawMoney(1, -5.00)

	balance, _ = api.GetBalance(1)
	fmt.Println("balance first user:", balance)



	//for ; ; {
	//	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	//		fmt.Fprintf(w, "Hello!\n")
	//	})
	//	http.ListenAndServe(":8080", nil)
	//}
	// c := make(chan os.Signal, 1)
	// signal.Notify(c, os.Interrupt)
	// <-c
}
