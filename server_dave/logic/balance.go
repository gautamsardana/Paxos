package logic

import (
	"context"
	"database/sql"
	"fmt"

	common "GolandProjects/apaxos-gautamsardana/api_common"
	"GolandProjects/apaxos-gautamsardana/server_dave/config"
	"GolandProjects/apaxos-gautamsardana/server_dave/storage/datastore"
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

		if len(resp.CommittedTxns) > 0 {
			serverLastCommittedTerm := resp.CommittedTxns[len(resp.CommittedTxns)-1].Term
			err = ReceiveCommit(ctx, conf, &common.Commit{
				BallotNum:         resp.BallotNum,
				AcceptVal:         resp.CommittedTxns,
				LastCommittedTerm: serverLastCommittedTerm,
			})
			if err != nil && err != ErrDuplicateTxns {
				return nil, err
			}
		}

		if len(resp.LogTxns) > 0 {
			for logTxnMsgID, logTxn := range resp.LogTxns {
				if logTxn.Receiver != conf.ClientName {
					continue
				}
				txn, logErr := datastore.GetTransactionByMsgID(conf.DataStore, logTxnMsgID)
				if logErr != nil && logErr != sql.ErrNoRows {
					return nil, logErr
				}
				if txn == nil {
					finalBalance += logTxn.Amount
				}
			}
		}
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

	var latestCommittedTxns []*common.TxnRequest
	if req.LastCommittedTerm < lastCommittedTerm {
		latestCommittedTxns, err = datastore.GetTransactionsAfterTerm(conf.DataStore, req.LastCommittedTerm)
		if err != nil {
			return nil, err
		}
	}
	resp.CommittedTxns = latestCommittedTxns
	resp.LogTxns = conf.LogStore.Logs

	fmt.Printf("Server %d: returning GetServerBalance with response: %v\n", conf.ServerNumber, resp)
	return resp, nil
}
