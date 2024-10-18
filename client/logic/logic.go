package logic

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"log"

	common "GolandProjects/apaxos-gautamsardana/api_common"
	"GolandProjects/apaxos-gautamsardana/client/config"
)

var mapUserToServer = map[string]string{
	"S1": "localhost:8080",
	"S2": "localhost:8081",
	"S3": "localhost:8082",
	"S4": "localhost:8083",
	"S5": "localhost:8084",
}

func ProcessTxnSet(ctx context.Context, req *common.TxnSet, conf *config.Config) error {
	isServerAlive := map[string]bool{
		"localhost:8080": false,
		"localhost:8081": false,
		"localhost:8082": false,
		"localhost:8083": false,
		"localhost:8084": false,
	}

	for _, aliveServer := range req.LiveServers {
		isServerAlive[mapUserToServer[aliveServer]] = true
	}

	for _, serverAddr := range conf.ServerAddresses {
		server, err := conf.Pool.GetServer(serverAddr)
		if err != nil {
			fmt.Println(err)
		}
		server.IsAlive(ctx, &common.IsAliveRequest{IsAlive: isServerAlive[serverAddr]})
	}

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

func PrintBalance(ctx context.Context, req *common.GetBalanceRequest, conf *config.Config) (*common.GetBalanceResponse, error) {
	serverAddr := mapUserToServer[req.User]
	server, err := conf.Pool.GetServer(serverAddr)
	if err != nil {
		fmt.Println(err)
	}
	resp, err := server.PrintBalance(ctx, req)
	if err != nil {
		return nil, err
	}
	return &common.GetBalanceResponse{
		Balance: resp.Balance,
	}, nil
}

func PrintLogs(ctx context.Context, req *common.PrintLogsRequest, conf *config.Config) (*common.PrintLogsResponse, error) {
	serverAddr := mapUserToServer[req.User]
	server, err := conf.Pool.GetServer(serverAddr)
	if err != nil {
		fmt.Println(err)
	}
	resp, err := server.PrintLogs(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func PrintDB(ctx context.Context, req *common.PrintDBRequest, conf *config.Config) (*common.PrintDBResponse, error) {
	serverAddr := mapUserToServer[req.User]
	server, err := conf.Pool.GetServer(serverAddr)
	if err != nil {
		fmt.Println(err)
	}
	resp, err := server.PrintDB(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func Performance(ctx context.Context, req *common.PerformanceRequest, conf *config.Config) (*common.PerformanceResponse, error) {
	serverAddr := mapUserToServer[req.User]
	server, err := conf.Pool.GetServer(serverAddr)
	if err != nil {
		fmt.Println(err)
	}
	resp, err := server.Performance(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
