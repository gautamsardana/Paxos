package logic

import (
	common "GolandProjects/apaxos-gautamsardana/api_common"
	"context"
	"database/sql"
	"fmt"
	"time"

	"GolandProjects/apaxos-gautamsardana/server_alice/config"
)

func TransactionWorker(conf *config.Config) {
	ticker := time.NewTicker(5 * time.Second)

	go func() {
		for {
			select {
			case <-ticker.C:
				QueueTransaction(conf)
			}
		}
	}()

}

func QueueTransaction(conf *config.Config) {
	conf.QueueMutex.Lock()
	defer conf.QueueMutex.Unlock()

	if len(conf.TxnQueue) > 0 {
		txn := conf.TxnQueue[0]

		err := ValidateTxn(conf, txn)
		if err != nil {
			fmt.Printf("Server %d: txn not valid", conf.ServerNumber)
			return
		}

		go func() {
			err = ProcessTxn(context.Background(), txn, conf)
			if err != nil {
				return
			}
		}()
	}
}

func ValidateTxn(conf *config.Config, txn *common.TxnRequest) error {
	dbTxn, err := ValidateTxnInDB(conf, txn)
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	if dbTxn != nil {
		err = fmt.Errorf("Server %d: duplicate txn found in db. Moving on to the next one..", conf.ServerNumber)
		conf.TxnQueue = conf.TxnQueue[1:]
		return err
	}

	logTxn := ValidateTxnInLogs(conf, txn)
	if logTxn != nil {
		err = fmt.Errorf("Server %d: duplicate txn found in logs. Moving on to the next one...\n", conf.ServerNumber)
		conf.TxnQueue = conf.TxnQueue[1:]
		return err
	}
	return nil
}
