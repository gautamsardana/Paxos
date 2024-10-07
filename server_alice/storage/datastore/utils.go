package datastore

import (
	"database/sql"
	"fmt"

	"GolandProjects/apaxos-gautamsardana/server_alice/storage"
)

func GetBalance(db *sql.DB, user string) (float32, error) {
	var balance float32
	query := `SELECT balance FROM user WHERE user = ?`
	err := db.QueryRow(query, user).Scan(&balance)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("no balance found for user: %d", user)
		}
		return 0, err
	}
	return balance, nil
}

func UpdateBalance(tx *sql.Tx, user storage.User) error {
	query := `UPDATE user SET balance = ? WHERE user = ?`
	res, err := tx.Exec(query, user.Balance, user.User)
	if err != nil {
		return err
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("no rows updated for user: %d", user)
	}
	return nil
}

func GetTransactionByMsgID(db *sql.DB, msgID string) (*storage.Transaction, error) {
	transaction := &storage.Transaction{}
	query := `SELECT * FROM transaction WHERE msg_id = ?`
	err := db.QueryRow(query, msgID).Scan(transaction)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	} else {
		return nil, nil
	}
}

func InsertTransactionLog(tx *sql.Tx, transaction storage.Transaction) error {
	query := `INSERT INTO transaction (msg_id, sender, receiver, amount) VALUES (?, ?, ?, ?)`
	_, err := tx.Exec(query, transaction.MsgID, transaction.Sender, transaction.Receiver, transaction.Amount)
	if err != nil {
		return err
	}
	return nil
}
