package logic

import (
	"context"
	"fmt"
	"time"

	common "GolandProjects/apaxos-gautamsardana/api_common"
	"GolandProjects/apaxos-gautamsardana/server_chucky/config"
	"GolandProjects/apaxos-gautamsardana/server_chucky/storage"
	"GolandProjects/apaxos-gautamsardana/server_chucky/storage/datastore"
)

func AddBallot(conf *config.Config, req *common.Promise) {
	conf.CurrVal.BallotNumber = req.BallotNum
}

func AddLocalTxns(conf *config.Config) {
	for _, txn := range conf.LogStore.Logs {
		conf.CurrVal.Transactions = append(conf.CurrVal.Transactions, txn)
	}
}

func AddNewTxns(conf *config.Config, req *common.Promise) {
	if req.LocalVal == nil {
		return
	}
	conf.CurrVal.Transactions = append(conf.CurrVal.Transactions, req.LocalVal...)
}

func CommitTransaction(ctx context.Context, conf *config.Config, req *common.Commit) (err error) {
	tx, err := conf.DataStore.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		}
	}()

	currClientBalance := conf.LogStore.Balance

	for _, txnDetails := range req.AcceptVal {
		transaction := storage.Transaction{
			MsgID:     txnDetails.MsgID,
			Sender:    txnDetails.Sender,
			Receiver:  txnDetails.Receiver,
			Amount:    txnDetails.Amount,
			Term:      txnDetails.Term,
			CreatedAt: time.Now(),
		}

		err = datastore.InsertTransaction(tx, transaction)
		if err != nil {
			fmt.Printf("error while inserting txn, err: %v", err)
			return fmt.Errorf("error while inserting txn: %v", err)
		}
		if txnDetails.Receiver == conf.ClientName {
			currClientBalance += txnDetails.Amount
		}
	}
	err = datastore.UpdateBalance(tx, storage.User{User: conf.ClientName, Balance: currClientBalance})
	if err != nil && err != datastore.ErrNoRowsUpdated {
		fmt.Printf("error while updating balance, err: %v\n", err)
		return fmt.Errorf("error while updating balance, err: %v", err)
	}
	conf.LogStore.Balance = currClientBalance

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("transaction commit failed: %v", err)
	}
	return nil
}

func DeleteFromLogs(conf *config.Config, transactions []*common.ProcessTxnRequest) {
	for _, txn := range transactions {
		if conf.LogStore.Logs[txn.MsgID] != nil {
			delete(conf.LogStore.Logs, txn.MsgID)
		}
	}
}
