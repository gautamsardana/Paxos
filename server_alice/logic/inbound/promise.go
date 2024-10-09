package inbound

import (
	"context"
	"fmt"

	common "GolandProjects/apaxos-gautamsardana/api_common"
	"GolandProjects/apaxos-gautamsardana/server_alice/config"
	"GolandProjects/apaxos-gautamsardana/server_alice/logic/outbound"
	"GolandProjects/apaxos-gautamsardana/server_alice/utils"
)

// i am the leader - these are all the promises i got

func Promise(ctx context.Context, conf *config.Config, req *common.Promise) error {
	// todo: add timeout - if no majority promises within timeout - fail
	// todo: add the majority check and if majority arrived, add localtxns at the end
	// todo - leader will also have acceptVal which will be updated when majority is received)

	if !req.PromiseAck {
		return fmt.Errorf("promise not acknowledged, request canceled")
	}

	if req.BallotNum.TermNumber != conf.CurrBallot.TermNumber ||
		req.BallotNum.ServerNumber != conf.CurrBallot.ServerNumber {
		return fmt.Errorf("ballot number mismatch, request canceled")
	}

	conf.CurrVal.CurrPromiseCount++
	conf.CurrVal.ServerAddresses = append(conf.CurrVal.ServerAddresses, utils.MapServerNumberToAddress[req.ServerNumber])

	if conf.CurrVal == nil {
		conf.CurrVal = config.CurrentValConstructor()
		AddBallot(conf, req)
	}
	//todo : this means that once the values are committed, you need to make currVal nil again

	//accept num and accept val are nil -- just add local txns of the follower
	if req.AcceptNum == nil && req.AcceptVal == nil {
		AddNewTxns(conf, req)
	} else {
		// accept num/val not empty -- update currVal of the leader to the acceptVal from follower
		if conf.CurrVal.MaxAcceptVal == nil || req.AcceptNum.TermNumber > conf.CurrVal.MaxAcceptVal.TermNumber {
			conf.CurrVal.MaxAcceptVal = req.AcceptNum
			conf.CurrVal.Transactions = req.AcceptVal
		}
	}

	// todo - once majority reached -------
	AddLocalTxns(conf)
	acceptRequest := &common.Accept{
		BallotNum:       conf.CurrVal.BallotNumber,
		AcceptVal:       conf.CurrVal.Transactions,
		ServerAddresses: conf.CurrVal.ServerAddresses,
	}
	outbound.Accept(ctx, conf, acceptRequest)
	return nil
}

func AddBallot(conf *config.Config, req *common.Promise) {
	conf.CurrVal.BallotNumber = req.BallotNum
}

func AddLocalTxns(conf *config.Config) {
	for _, txn := range conf.LogStore.Logs {
		conf.CurrVal.Transactions = append(conf.CurrVal.Transactions, txn)
	}
}

func AddNewTxns(conf *config.Config, req *common.Promise) {
	conf.CurrVal.Transactions = append(conf.CurrVal.Transactions, req.LocalVal...)
}
