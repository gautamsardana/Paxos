package logstore

import (
	common "GolandProjects/apaxos-gautamsardana/api_common"
	"sync"
)

type LogStore struct {
	Logs    []*common.ProcessTxnRequest
	Balance float32
	Mu      sync.Mutex
}

func NewLogStore(balance float32) *LogStore {
	return &LogStore{
		Balance: balance,
		Logs:    []*common.ProcessTxnRequest{},
	}
}

func (store *LogStore) GetBalance() float32 {
	store.Mu.Lock()
	defer store.Mu.Unlock()

	return store.Balance
}

func (store *LogStore) AddTransactionLog(log *common.ProcessTxnRequest) {
	store.Mu.Lock()
	defer store.Mu.Unlock()

	store.Logs = append(store.Logs, log)
	store.Balance -= log.Amount
}
