package main

import (
	"fmt"
	"os"
	// "strconv"
	"log"
	"database/sql"
	_ "github.com/jackc/pgx"
	// "net/http"
	// "encoding/json"
	pkgpostgres "app/pkg/postgres"
)

type User struct {
	ID			int		`json:"id"`
	Balance		float64	`json:"balance"`
}

type DB *sql.DB

// func (db *DB) TransferAndWithdrawMoney(id int, sum float64) (error) {
// 	var op string

// 	if operation == transfer {
// 		op = "+"
// 	} else {
// 		op = "-"
// 	}

// 	b, err := db.GetBalance(id)

// 	if err != nil {
// 		return err
// 	}

// 	if b + sum < 0.00 {
// 		// e := fmt.Sprintf("Insufficient funds: %v", id)
// 		return fmt.Errorf("Insufficient funds: %v", id)
// 	}

// 	_, err = db.Query(`UPDATE user_balance SET balance = balanse $1 $2 WHERE id = $3`, op, sum, id)
// 	if err != nil {
// 	}
// }

// func (db *DB) TransferMoney(id_first, id_second int, sum float64) (error) {
// 	b, err := db.GetBalance(id_first)
// 	if err != nil {
// 		return err
// 	}

// 	if b -  sum < 0 {
// 		e := fmt.Sprintf("the client with id %v does not have enough funds", id)
// 		return error.New(e)
// 	}


// }

// func (db *DB) GetBalance(id int) (float64, error) {
// 	var byte_of_balance	[]byte
// 	var balance			float64

// 	rows, err := db.Query(`SELECT balance FROM user_balance WHERE id = $1`, id)
// 	if err != nil {
// 		fmt.Printf("GetBalance() -> Query(): %v", err)
// 		return 0, err
// 	}

// 	rows.Next()
// 	err = rows.Scan(&byte_of_balance)
// 	if err != nil {
// 		fmt.Printf("GetBalance() -> Scan(): %v", err)
// 		return 0, err
// 	}

// 	balance, err = strconv.ParseFloat(string(byte_of_balance), 64)
// 	if err != nil {
// 		fmt.Println("GetBalance -> strconv.ParseFloat:", err)
// 		return 0, err
// 	}

// 	return balance, nil
// }

func main() {
	var err error
	// var db DB

	// var (
	// 	psql_db = string(os.Getenv("POSTGRES_DB"))
	// 	psql_user = string(os.Getenv("POSTGRES_USER"))
	// 	psql_psw = string(os.Getenv("POSTGRES_PASSWORD"))
	// )

	// conn := string(fmt.Sprintf("dbname=%v user=%v password=%v host=localhost port=5432 sslmode=disable",
							// psql_db, psql_user, psql_psw))

	// conn := fmt.Sprintf("postgres://%v:%v@localhost:5432/%v?sslmode=disable",
	// 					psql_user, psql_psw, psql_db)

	fmt.Println(os.Getenv("DATABASE_URL"))

	pool, err := pkgpostgres.NewPool(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
		return
	}

	defer pool.Close()

	// for ; ; {
	// 	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 		fmt.Fprintf(w, "Hello!\n")
	// 	})
	// 	http.ListenAndServe(":8080", nil)
	// }
	// c := make(chan os.Signal, 1)
	// signal.Notify(c, os.Interrupt)
	// <-c
}
