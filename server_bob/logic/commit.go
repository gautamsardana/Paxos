package logic

import (
	"GolandProjects/apaxos-gautamsardana/server_bob/storage/datastore"
	"context"
	"fmt"

	common "GolandProjects/apaxos-gautamsardana/api_common"
	"GolandProjects/apaxos-gautamsardana/server_bob/config"
	"GolandProjects/apaxos-gautamsardana/server_bob/utils"
)

// i am a leader - i got accepted requests from majority followers. Now I need to tell them to commit

func SendCommit(ctx context.Context, conf *config.Config, req *common.Commit) {
	for _, txn := range req.AcceptVal {
		txn.Term = req.BallotNum.TermNumber
	}

	// if commitReq lastCommittedTerm is same as mine, meaning i have already pushed these txns before in my db
	req.LastCommittedTerm = req.BallotNum.TermNumber

	dbErr := CommitTransaction(ctx, conf, req)
	if dbErr != nil {
		return
	}
	DeleteFromLogs(conf, req.AcceptVal)
	config.ResetCurrVal(conf)
	config.ResetAcceptVal(conf)
	fmt.Printf("Server %d: new txns committed, updated log:%v\n", conf.ServerNumber, conf.LogStore)

	go func() {
		err := ProcessTxn(ctx, conf.CurrTxn, conf, true)
		if err != nil {
			fmt.Printf("Server %d: error processing txn:%v\n", conf.ServerNumber, conf.CurrTxn)
		}
	}()

	fmt.Printf("Server %d: sending commit request with req: %v\n", conf.ServerNumber, req)
	for _, serverAddress := range req.ServerAddresses {
		server, err := conf.Pool.GetServer(serverAddress)
		if err != nil {
			fmt.Println(err)
		}

		_, err = server.Commit(ctx, req)
		if err != nil {
			fmt.Println(err)
		}
	}
}

//i am a follower - i received this commit message from the leader. I need to commit these messages
// in my db and delete those txns from log

func ReceiveCommit(ctx context.Context, conf *config.Config, req *common.Commit) error {
	//todo, also check if the same txns are committed in your db already
	fmt.Printf("Server %d: received commit from leader with request: %v\n", conf.ServerNumber, req)

	lastCommittedTerm, err := datastore.GetLatestTermNo(conf.DataStore)
	if err != nil {
		return err
	}
	if req.LastCommittedTerm <= lastCommittedTerm {
		fmt.Printf("Server %d: outdated commit request: %v\n", conf.ServerNumber, req)
		return nil
	}

	if req.BallotNum.TermNumber > conf.CurrBallot.TermNumber {
		utils.UpdateBallot(conf, req.BallotNum.TermNumber, req.BallotNum.ServerNumber)
	}

	err = CommitTransaction(ctx, conf, req)
	if err != nil {
		return err
	}

	DeleteFromLogs(conf, req.AcceptVal)
	config.ResetAcceptVal(conf)
	fmt.Printf("Server %d: new txns committed, updated log:%v\n", conf.ServerNumber, conf.LogStore)

	if req.IsSync {
		go func() {
			err = ProcessTxn(context.Background(), conf.CurrTxn, conf, true)
			if err != nil {

			}
		}()
	}
	// todo update balance for all

	return nil
}
