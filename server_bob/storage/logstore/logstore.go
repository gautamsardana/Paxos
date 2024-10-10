package logstore

import (
	common "GolandProjects/apaxos-gautamsardana/api_common"
)

type LogStore struct {
	Logs    map[string]*common.ProcessTxnRequest
	Balance float32
}

func NewLogStore(balance float32) *LogStore {
	return &LogStore{
		Balance: balance,
		Logs:    make(map[string]*common.ProcessTxnRequest),
	}
}

func (store *LogStore) GetBalance() float32 {
	return store.Balance
}

func (store *LogStore) AddTransactionLog(txn *common.ProcessTxnRequest) {
	store.Logs[txn.MsgID] = txn
	store.Balance -= txn.Amount
}
