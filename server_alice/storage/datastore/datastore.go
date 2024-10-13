package datastore

import (
	"database/sql"
	"errors"
	"fmt"

	common "GolandProjects/apaxos-gautamsardana/api_common"
	"GolandProjects/apaxos-gautamsardana/server_alice/storage"
)

var ErrNoRowsUpdated = errors.New("no rows updated for user")

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
		return ErrNoRowsUpdated
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

func InsertTransaction(tx *sql.Tx, transaction storage.Transaction) error {
	query := `INSERT INTO transaction (msg_id, sender, receiver, amount, term, created_at) VALUES (?, ?, ?, ?, ?, ?)`
	_, err := tx.Exec(query, transaction.MsgID, transaction.Sender, transaction.Receiver,
		transaction.Amount, transaction.Term, transaction.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}

func GetLatestTermNo(db *sql.DB) (int32, error) {
	query := `SELECT term FROM transaction order by created_at desc limit 1`
	var term int32
	err := db.QueryRow(query).Scan(&term)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}
	return term, nil
}

func GetTransactionsAfterTerm(db *sql.DB, term int32) ([]*common.ProcessTxnRequest, error) {
	var transactions []*common.ProcessTxnRequest

	query := `SELECT msg_id, sender, receiver, amount, term FROM transaction WHERE term > ? ORDER BY created_at`
	rows, err := db.Query(query, term)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var txn common.ProcessTxnRequest
		if err = rows.Scan(&txn.MsgID, &txn.Sender, &txn.Receiver, &txn.Amount, &txn.Term); err != nil {
			return nil, err
		}
		transactions = append(transactions, &txn)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return transactions, nil
}
