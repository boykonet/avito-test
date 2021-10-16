package methods

import (
	"context"
	"fmt"
)

const (
	transfer = `UPDATE user_balance SET balance = balance + $1 WHERE id = $2`
	withdraw = `UPDATE user_balance SET balance = balance - $1 WHERE id = $2`
	request = `SELECT balance FROM user_balance WHERE id = $1`
)

// The TransferAndWithdrawMoney method takes two parameters:
// the first parameter is the user ID,
// the second parameter is the amount of money that needs to be transferred or withdraw from user's account.
// If the amount of money is positive, the amount of money is transferred to user's account,
// otherwise, the amount of money is withdrawn from the user's account.
// On success, nil is returned. Otherwise, an error is returned
func (db *Methods) TransferAndWithdrawMoney(id int, sum float64) (error) {

	balance, err := db.GetBalance(id)
	if err != nil {
		return fmt.Errorf("GetBalance() error: %w", err)
	}

	fmt.Println()
	if balance + sum < 0.00 {
		return fmt.Errorf("Insufficient funds: %v", id)
	}

	_, err = db.pool.Exec(context.Background(), transfer, sum, id)
	if err != nil {
		return fmt.Errorf("Exec() error: %w", err)
	}

	return nil
}

// The TransferMoney method takes three parameters:
// thr first parameter is the user's id who transfers money,
// the second parameter is the user's id to whom the money is transferred,
// the third parameter is the amount of money to be transferred from the first user to the second user.
// The amount of money should be only positive. If the amount of money is negative, an error is returned.
// On success, nil is returned.
func (db *Methods) TransferMoney(first_id, second_id int, sum float64) (error) {

	if sum < 0.00 {
		return fmt.Errorf("Amount of money is negative: %v", sum)
	}

	balance, err := db.GetBalance(first_id)
	if err != nil {
		return fmt.Errorf("GetBalance error: %w", err)
	}

	if balance - sum < 0.00 {
		return fmt.Errorf("Insufficient funds: %v", first_id)
	}

	// Withdraw amount of money from first user
	_, err = db.pool.Exec(context.Background(), withdraw, sum, first_id)
	if err != nil {
		return fmt.Errorf("Exec(..., first_id) error: %w", err)
	}

	// Transfer amount of money to second user
	_, err = db.pool.Exec(context.Background(), transfer, sum, second_id)
	if err != nil {
		return fmt.Errorf("Exec(..., second_id) error: %w", err)
	}

	return nil
}


// The GetBalance method, in successful, returns the amount of the (user's money, nil)
// Otherwise, the returns an (0, error)
func (db *Methods) GetBalance(id int) (float64, error) {
	var balance float64

	rows, err := db.pool.Query(context.Background(), request, id)
	if err != nil {
		return 0, fmt.Errorf("Query() error: %w", err)
	}

	defer rows.Close()

	rows.Next()
	err = rows.Scan(&balance)
	if err != nil {
		return 0, fmt.Errorf("Scanf() error: %w", err)
	}
	return balance, nil
}

