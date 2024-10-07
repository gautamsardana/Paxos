package logic

import (
	"context"
	"fmt"

	common "GolandProjects/apaxos-gautamsardana/api_common"
	"GolandProjects/apaxos-gautamsardana/server_alice/config"
	logic2 "GolandProjects/apaxos-gautamsardana/server_alice/logic/outbound"
	"GolandProjects/apaxos-gautamsardana/server_alice/storage"
	"GolandProjects/apaxos-gautamsardana/server_alice/storage/datastore"
)

func ProcessTxn(ctx context.Context, req *common.ProcessTxnRequest, conf *config.Config) error {
	transaction, err := datastore.GetTransactionByMsgID(conf.DataStore, req.MsgID)
	if err != nil {
		return err
	}
	if transaction != nil {
		return fmt.Errorf("transaction with same msgID processed before")
	}

	balance := conf.LogStore.GetBalance()

	if balance >= req.Amount {
		err = ExecuteTxn(ctx, balance, req, conf)
		if err != nil {
			return err
		}
	} else {
		fmt.Println("this is where the magic happens!")
		// initiate paxos
		logic2.Prepare(ctx, conf)
	}
	return nil
}

func ExecuteTxn(ctx context.Context, balance float32, req *common.ProcessTxnRequest, conf *config.Config) error {
	txnDetails := storage.Transaction{
		MsgID:    req.MsgID,
		Sender:   req.Sender,
		Receiver: req.Receiver,
		Amount:   req.Amount,
	}
	fmt.Println(conf.LogStore.Logs, conf.LogStore.Balance)

	conf.LogStore.AddTransactionLog(txnDetails)

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
