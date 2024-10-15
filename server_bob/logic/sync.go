package logic

import (
	"GolandProjects/apaxos-gautamsardana/server_bob/utils"
	"context"
	"fmt"

	common "GolandProjects/apaxos-gautamsardana/api_common"
	"GolandProjects/apaxos-gautamsardana/server_bob/config"
	"GolandProjects/apaxos-gautamsardana/server_bob/storage/datastore"
)

// SendSyncResponse i am the updated server. I need to send latest txns to slow servers
func SendSyncResponse(ctx context.Context, conf *config.Config, req *common.SyncRequest) error {
	latestTxns, err := datastore.GetTransactionsAfterTerm(conf.DataStore, req.LastCommittedTerm)
	if err != nil {
		return err
	}
	lastCommittedTerm := int32(0)

	if len(latestTxns) > 0 {
		lastCommittedTerm = latestTxns[len(latestTxns)-1].Term
	}

	commitRequest := &common.Commit{
		BallotNum:         conf.CurrBallot,
		AcceptVal:         latestTxns,
		LastCommittedTerm: lastCommittedTerm,
	}
	server, err := conf.Pool.GetServer(utils.MapServerNumberToAddress[req.ServerNo])
	if err != nil {
		fmt.Println(err)
	}

	_, err = server.Commit(ctx, commitRequest)
	if err != nil {
		fmt.Println(err)
	}
	return nil
}

// SyncRequest i am the updated server. I received this sync request from a slow server
func SyncRequest(ctx context.Context, conf *config.Config, req *common.SyncRequest) error {
	err := SendSyncResponse(ctx, conf, req)
	if err != nil {
		return err
	}
	return nil
}
