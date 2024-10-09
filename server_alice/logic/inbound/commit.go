package inbound

import (
	"context"
	"fmt"
	"log"

	common "GolandProjects/apaxos-gautamsardana/api_common"
	"GolandProjects/apaxos-gautamsardana/server_alice/config"
	"GolandProjects/apaxos-gautamsardana/server_alice/storage"
	"GolandProjects/apaxos-gautamsardana/server_alice/storage/datastore"
)

//i am a follower - i received this commit message from the leader. I need to commit these messages
// in my db and delete those txns from log

func Commit(ctx context.Context, conf *config.Config, req *common.Commit) (err error) {
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

	for _, txnDetails := range req.AcceptVal {
		transaction := storage.Transaction{
			MsgID:    txnDetails.MsgID,
			Sender:   txnDetails.Sender,
			Receiver: txnDetails.Receiver,
			Amount:   txnDetails.Amount,
		}

		err = datastore.InsertTransaction(tx, transaction)
		if err != nil {
			log.Printf("error while inserting txn, err: %v", err)
			return fmt.Errorf("error while inserting txn: %v", err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("transaction commit failed: %v", err)
	}
	DeleteFromLogs(conf, req.AcceptVal)

	return nil
}

func DeleteFromLogs(conf *config.Config, transactions []*common.ProcessTxnRequest) {
	for _, txn := range transactions {
		if conf.LogStore.Logs[txn.MsgID] != nil {
			delete(conf.LogStore.Logs, txn.MsgID)
		}
	}
}
