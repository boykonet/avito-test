package main

import (
	"fmt"
	"database/sql"
	"log"
	"strconv"
	// "encoding/binary"
	// "math"
	_ "github.com/lib/pq"
)

type User struct {
	ID			int
	Balance		[]byte
}

// func Float64fromBytes(bytes []byte) float64 {
//     bits := binary.LittleEndian.Uint64(bytes)
//     float := math.Float64frombits(bits)
//     return float
// }

func main() {
	db, err := sql.Open("postgres", "host=localhost port=5432 user=gkarina_user password=gkarina_password dbname=gkarina_database sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	// rows, err := db.Query("SELECT id, balance FROM user_balance")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// for rows.Next() {
	// 	var user User
	// 	err := rows.Scan(&user.ID, &user.Balance)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	balance, _ := strconv.ParseFloat(string(user.Balance), 64)
	// 	// fmt.Printf("%T\n", balance)
	// 	fmt.Println(user.ID, balance)
	// }

	// SELECT
	var id int
	var balance []byte
	rows, err := db.Query(`SELECT id, balance FROM user_balance WHERE id = $1`, 4)
	if err != nil {
		fmt.Println(err)
	}

	rows.Next()
	rows.Scan(&id, &balance)
	b, _ := strconv.ParseFloat(string(balance), 64)
	fmt.Println(id, b)


	// UPDATE
	_, err = db.Query(`UPDATE user_balance SET balance=$1 WHERE id=$2`, 66.66, 7)


	// rows, err = db.Query(`INSERT INTO user_balance(balance) VALUES($1) RETURNING id, balance`, 33.33)
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// rows.Next()
	// rows.Scan(&id, &balance)
	// b, _ = strconv.ParseFloat(string(balance), 64)
	// fmt.Println(id, b)

	rows, err = db.Query(`SELECT id, balance FROM user_balance WHERE id = $1`, 7)
	if err != nil {
		fmt.Println(err)
	}

	rows.Next()
	rows.Scan(&id, &balance)
	b, _ = strconv.ParseFloat(string(balance), 64)
	fmt.Println(id, b)

}
