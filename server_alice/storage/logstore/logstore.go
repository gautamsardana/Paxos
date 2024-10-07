package logstore

import (
	"sync"

	"GolandProjects/apaxos-gautamsardana/server_alice/storage"
)

type LogStore struct {
	Logs    []storage.Transaction
	Balance float32
	Mu      sync.Mutex
}

func NewLogStore(balance float32) *LogStore {
	return &LogStore{
		Balance: balance,
		Logs:    []storage.Transaction{},
	}
}

func (store *LogStore) GetBalance() float32 {
	store.Mu.Lock()
	defer store.Mu.Unlock()

	return store.Balance
}

func (store *LogStore) AddTransactionLog(log storage.Transaction) {
	store.Mu.Lock()
	defer store.Mu.Unlock()

	store.Logs = append(store.Logs, log)
	store.Balance -= log.Amount
}
