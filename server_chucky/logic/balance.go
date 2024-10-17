package logic

import (
	"context"
	"database/sql"
	"fmt"

	common "GolandProjects/apaxos-gautamsardana/api_common"
	"GolandProjects/apaxos-gautamsardana/server_chucky/config"
	"GolandProjects/apaxos-gautamsardana/server_chucky/storage/datastore"
)

func PrintBalance(ctx context.Context, conf *config.Config) (*common.GetBalanceResponse, error) {
	fmt.Printf("Server %d: GetBalance request received\n", conf.ServerNumber)

	finalBalance := float32(0)

	for _, serverAddr := range conf.ServerAddresses {
		server, err := conf.Pool.GetServer(serverAddr)
		if err != nil {
			fmt.Println(err)
		}

		resp, err := server.GetServerBalance(ctx, &common.GetServerBalanceRequest{
			LastCommittedTerm: conf.LastCommittedTerm,
			User:              conf.ClientName,
		})
		if err != nil {
			fmt.Println(err)
		}

		if resp.CommittedTxns != nil && len(resp.CommittedTxns) > 0 {
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

		if resp.LogTxns != nil && len(resp.LogTxns) > 0 {
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

	var latestCommittedTxns []*common.TxnRequest
	var err error

	if req.LastCommittedTerm < conf.LastCommittedTerm {
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
