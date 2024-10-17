package logic

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	common "GolandProjects/apaxos-gautamsardana/api_common"
	"GolandProjects/apaxos-gautamsardana/server_alice/config"
	"GolandProjects/apaxos-gautamsardana/server_alice/storage/datastore"
	"GolandProjects/apaxos-gautamsardana/server_alice/utils"
)

var ErrDuplicateTxns = errors.New("duplicate txns")

// i am a leader - i got accepted requests from majority followers. Now I need to tell them to commit

func SendCommit(ctx context.Context, conf *config.Config, req *common.Commit) {
	if req.BallotNum.TermNumber < conf.CurrBallot.TermNumber {
		return
	}
	for _, txn := range req.AcceptVal {
		txn.Term = req.BallotNum.TermNumber
	}

	if len(req.AcceptVal) > 0 {
		msgID := req.AcceptVal[0].MsgID
		txn, dbErr := datastore.GetTransactionByMsgID(conf.DataStore, msgID)
		if dbErr != nil && dbErr != sql.ErrNoRows {
			return
		}
		if txn != nil {
			fmt.Printf("txn exists in db already")
			return
		}
	}

	dbErr := CommitTransaction(ctx, conf, req)
	if dbErr != nil {
		return
	}
	// if commitReq lastCommittedTerm is same as mine, meaning i have already pushed these txns before in my db
	conf.LastCommittedTerm = req.BallotNum.TermNumber
	req.LastCommittedTerm = req.BallotNum.TermNumber

	DeleteFromLogs(conf, req.AcceptVal)
	config.ResetCurrVal(conf)
	config.ResetAcceptVal(conf)
	config.NewAcceptedServersInfo(conf)
	fmt.Printf("Server %d: new txns committed, updated log:%v\n", conf.ServerNumber, conf.LogStore)

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

	if req.BallotNum.TermNumber > conf.CurrBallot.TermNumber {
		utils.UpdateBallot(conf, req.BallotNum.TermNumber, req.BallotNum.ServerNumber)
	}

	if req.LastCommittedTerm <= conf.LastCommittedTerm {
		fmt.Printf("Server %d: outdated commit request: %v\n", conf.ServerNumber, req)
		return nil
	}

	if len(req.AcceptVal) > 0 {
		msgID := req.AcceptVal[0].MsgID
		txn, dbErr := datastore.GetTransactionByMsgID(conf.DataStore, msgID)
		if dbErr != nil && dbErr != sql.ErrNoRows {
			return dbErr
		}
		if txn != nil {
			fmt.Printf("Server %d: duplicate txns %v\n", conf.ServerNumber, req)
			return ErrDuplicateTxns
		}
	}

	err := CommitTransaction(ctx, conf, req)
	if err != nil {
		return err
	}

	//set the last committed term here after committing and retrieve this in future
	lastCommittedTerm, dbErr := datastore.GetLatestTermNo(conf.DataStore)
	if dbErr != nil {
		return dbErr
	}
	conf.LastCommittedTerm = lastCommittedTerm

	DeleteFromLogs(conf, req.AcceptVal)
	config.ResetAcceptVal(conf)
	fmt.Printf("Server %d: new txns committed, updated log:%v\n", conf.ServerNumber, conf.LogStore)

	return nil
}
