package logic

import (
	"context"
	"fmt"
	"time"

	common "GolandProjects/apaxos-gautamsardana/api_common"
	"GolandProjects/apaxos-gautamsardana/server_alice/config"
	"GolandProjects/apaxos-gautamsardana/server_alice/storage/datastore"
)

// todo - also add the balances of each server

func ProcessTxn(ctx context.Context, req *common.ProcessTxnRequest, conf *config.Config, retryFlag bool) error {
	fmt.Printf("Server %d: received process txn request:%v\n", conf.ServerNumber, req)

	conf.StartTime = time.Now()

	transaction, err := datastore.GetTransactionByMsgID(conf.DataStore, req.MsgID)
	if err != nil {
		return err
	}
	if transaction != nil {
		return fmt.Errorf("transaction with same msgID processed before")
	}

	balance := conf.LogStore.Balance

	if balance >= req.Amount {
		err = ExecuteTxn(ctx, req, conf)
		if err != nil {
			return err
		}
	} else {
		fmt.Println("this is where the magic happens!")

		conf.CurrTxn = req
		if retryFlag && conf.CurrRetryCount >= conf.RetryLimit {
			fmt.Println("txn failed after 1 attempt")
			return fmt.Errorf("txn failed after 1 attempt")
		}
		conf.CurrRetryCount++
		SendPrepare(context.Background(), conf)
	}
	fmt.Printf("-------- %s\n", time.Since(conf.StartTime))
	return nil
}

func ExecuteTxn(ctx context.Context, req *common.ProcessTxnRequest, conf *config.Config) error {
	fmt.Printf("Server %d: executing txn: %v\n", conf.ServerNumber, req)
	conf.LogStore.AddTransactionLog(req)

	fmt.Println(conf.LogStore.Logs, conf.LogStore.Balance)
	return nil
}

//func ExecuteTxn(ctx context.Context, balance float32, req *common.ProcessTxnRequest, conf *config.Config) error {
//	tx, err := conf.DataStore.BeginTx(ctx, nil)
//	if err != nil {
//		return fmt.Errorf("failed to begin transaction: %v", err)
//	}
//
//	defer func() {
//		if err != nil {
//			tx.Rollback()
//		}
//	}()
//
//	txnDetails := storage.Transaction{
//		MsgID:    req.MsgId,
//		Sender:   req.Sender,
//		Receiver: req.Receiver,
//		Amount:   req.Amount,
//	}
//
//	conf.LogStore.AddTransactionLog(txnDetails)
//
//	updatedBalance := balance - req.Amount
//	updatedUser := storage.User{
//		User:    req.Sender,
//		Balance: updatedBalance,
//	}
//	err = datastore.UpdateBalance(tx, updatedUser)
//	if err != nil {
//		log.Printf("error while updating user balance, err: %v", err)
//		return err
//	}
//
//	err = tx.Commit()
//	if err != nil {
//		return fmt.Errorf("transaction commit failed: %v", err)
//	}
//
//	fmt.Println(conf.LogStore)
//	return nil
//}
