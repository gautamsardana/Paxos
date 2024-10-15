package logic

import (
	common "GolandProjects/apaxos-gautamsardana/api_common"
	"GolandProjects/apaxos-gautamsardana/client/config"
	"context"
	"fmt"
	"github.com/google/uuid"
	"log"
)

var mapUserToServer = map[string]string{
	"S1": "localhost:8080",
	"S2": "localhost:8081",
	"S3": "localhost:8082",
	"S4": "localhost:8083",
	"S5": "localhost:8084",
}

func ProcessTxn(ctx context.Context, req *common.TxnSet, conf *config.Config) error {
	for _, txn := range req.Txns {
		serverAddr := mapUserToServer[txn.Sender]
		server, err := conf.Pool.GetServer(serverAddr)
		if err != nil {
			fmt.Println(err)
		}

		msgId, err := uuid.NewUUID()
		if err != nil {
			log.Fatalf("failed to generate UUID: %v", err)
		}

		txn.MsgID = msgId.String()
		_, err = server.EnqueueTxn(ctx, txn)
		if err != nil {
			return err
		}
	}
	return nil
}

func GetBalance(ctx context.Context, req *common.GetBalanceRequest, conf *config.Config) (*common.GetBalanceResponse, error) {
	serverAddr := mapUserToServer[req.User]
	server, err := conf.Pool.GetServer(serverAddr)
	if err != nil {
		fmt.Println(err)
	}
	resp, err := server.GetBalance(ctx, req)
	if err != nil {
		return nil, err
	}
	return &common.GetBalanceResponse{
		Balance: resp.Balance,
	}, nil
}
