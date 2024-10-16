package logic

import (
	"context"
	"fmt"

	common "GolandProjects/apaxos-gautamsardana/api_common"
	"GolandProjects/apaxos-gautamsardana/server_chucky/config"
	"GolandProjects/apaxos-gautamsardana/server_chucky/storage/datastore"
)

func PrintBalance(ctx context.Context, conf *config.Config) (*common.GetBalanceResponse, error) {
	fmt.Printf("Server %d: GetBalance request received\n", conf.ServerNumber)

	finalBalance := float32(0)
	lastCommittedTerm, dbErr := datastore.GetLatestTermNo(conf.DataStore)
	if dbErr != nil {
		return nil, dbErr
	}
	for _, serverAddr := range conf.ServerAddresses {
		server, err := conf.Pool.GetServer(serverAddr)
		if err != nil {
			fmt.Println(err)
		}

		resp, err := server.GetServerBalance(ctx, &common.GetServerBalanceRequest{
			LastCommittedTerm: lastCommittedTerm,
			User:              conf.ClientName,
		})
		if err != nil {
			fmt.Println(err)
		}

		if len(resp.NewTxns) > 0 {
			serverLastCommittedTerm := resp.NewTxns[len(resp.NewTxns)-1].Term
			err = ReceiveCommit(ctx, conf, &common.Commit{
				BallotNum:         resp.BallotNum,
				AcceptVal:         resp.NewTxns,
				ServerAddresses:   nil,
				LastCommittedTerm: serverLastCommittedTerm,
			})
			if err != nil && err != ErrDuplicateTxns {
				return nil, err
			}
		}
		finalBalance += resp.Balance
	}
	finalBalance += conf.Balance
	return &common.GetBalanceResponse{Balance: finalBalance}, nil
}

func GetServerBalance(ctx context.Context, conf *config.Config, req *common.GetServerBalanceRequest) (*common.GetServerBalanceResponse, error) {
	fmt.Printf("Server %d: received GetServerBalance with request: %v\n", conf.ServerNumber, req)

	resp := &common.GetServerBalanceResponse{
		BallotNum: conf.CurrBallot,
	}

	lastCommittedTerm, err := datastore.GetLatestTermNo(conf.DataStore)
	if err != nil {
		return nil, err
	}

	var latestTxns []*common.TxnRequest
	if req.LastCommittedTerm < lastCommittedTerm {
		latestTxns, err = datastore.GetTransactionsAfterTerm(conf.DataStore, req.LastCommittedTerm)
		if err != nil {
			return nil, err
		}
	}
	resp.NewTxns = latestTxns

	var balance float32
	for _, txn := range conf.LogStore.Logs {
		if txn.Receiver == req.User {
			balance += txn.Amount
		}
	}
	resp.Balance = balance

	fmt.Printf("Server %d: returning GetServerBalance with response: %v\n", conf.ServerNumber, resp)
	return resp, nil
}
