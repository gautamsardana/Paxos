package inbound

import (
	"context"
	"fmt"

	common "GolandProjects/apaxos-gautamsardana/api_common"
	"GolandProjects/apaxos-gautamsardana/server_alice/config"
	"GolandProjects/apaxos-gautamsardana/server_alice/storage"
)

func Promise(ctx context.Context, req *common.Promise, conf *config.Config) error {
	// todo: add timeout - if no majority promises within timeout - fail

	if !req.PromiseAck {
		return fmt.Errorf("promise not acknowledged, request canceled")
	}

	if req.BallotNum.TermNumber != conf.CurrBallot.TermNumber ||
		req.BallotNum.ServerNumber != conf.CurrBallot.ServerNumber {
		return fmt.Errorf("ballot number mismatch, request canceled")
	}

	if conf.CurrVal == nil {
		conf.CurrVal = config.CurrentValConstructor()
		AddBallot(conf, req)
		AddLocalTxns(conf)
	}
	//todo : this means that once the values are committed, you need to make currVal nil again

	if req.AcceptNum == nil && req.AcceptVal == nil {
		conf.CurrVal.ServersRespNumber++
		AddNewTxns(conf, req)
	} else {
		conf.CurrVal.ServersRespNumber++
		UpdateCurrVal(conf, req)
	}
	return nil
}

func AddBallot(conf *config.Config, req *common.Promise) {
	conf.CurrVal.BallotNumber.TermNumber = req.BallotNum.TermNumber
	conf.CurrVal.BallotNumber.ServerNumber = req.BallotNum.ServerNumber
}

func AddLocalTxns(conf *config.Config) {
	conf.CurrVal.Transactions = append(conf.CurrVal.Transactions, conf.LogStore.Logs...)
}

func AddNewTxns(conf *config.Config, req *common.Promise) {
	for _, localVal := range req.LocalVal {
		txn := storage.Transaction{
			MsgID:    localVal.MsgID,
			Sender:   localVal.Sender,
			Receiver: localVal.Receiver,
			Amount:   localVal.Amount,
		}
		conf.CurrVal.Transactions = append(conf.CurrVal.Transactions, txn)
	}
}

func UpdateCurrVal(conf *config.Config, req *common.Promise) {
	var newBlock []storage.Transaction

	for _, transaction := range req.LocalVal {
		txn := storage.Transaction{
			MsgID:    transaction.MsgID,
			Sender:   transaction.Sender,
			Receiver: transaction.Receiver,
			Amount:   transaction.Amount,
		}
		newBlock = append(newBlock, txn)
	}
	conf.CurrVal.Transactions = newBlock
}
