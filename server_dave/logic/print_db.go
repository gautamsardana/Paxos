package logic

import (
	"context"

	common "GolandProjects/apaxos-gautamsardana/api_common"
	"GolandProjects/apaxos-gautamsardana/server_dave/config"
	"GolandProjects/apaxos-gautamsardana/server_dave/storage/datastore"
)

func PrintDB(ctx context.Context, conf *config.Config, req *common.PrintDBRequest) (*common.PrintDBResponse, error) {
	txns, err := datastore.GetAllTransactions(conf.DataStore)
	if err != nil {
		return nil, err
	}
	return &common.PrintDBResponse{Txns: txns}, nil
}
