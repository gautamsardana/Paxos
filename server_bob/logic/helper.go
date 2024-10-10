package logic

import (
	common "GolandProjects/apaxos-gautamsardana/api_common"
	"GolandProjects/apaxos-gautamsardana/server_bob/config"
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
	conf.CurrVal.Transactions = append(conf.CurrVal.Transactions, req.LocalVal...)
}
