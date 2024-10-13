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
			fmt.Println("txn failed after 3 attempts")
			return fmt.Errorf("txn failed after 3 attempts")
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
