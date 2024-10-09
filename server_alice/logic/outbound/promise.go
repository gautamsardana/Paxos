package outbound

import (
	"context"
	"fmt"

	common "GolandProjects/apaxos-gautamsardana/api_common"
	"GolandProjects/apaxos-gautamsardana/server_alice/config"
	"GolandProjects/apaxos-gautamsardana/server_alice/utils"
)

// i am a follower - i got the leader's prepare. I need to send back my promise

func Promise(ctx context.Context, conf *config.Config, ballotNumber *common.Ballot) {
	promiseReq := &common.Promise{
		PromiseAck:   true,
		ServerNumber: conf.ServerNumber,
		BallotNum:    ballotNumber,
	}

	if conf.AcceptVal == nil {
		// send current local logs
		//conf.CurrPromiseSeq++
		//promiseReq.LocalVal = conf.LogStore.Logs[conf.CurrPromiseSeq]

		for _, txn := range conf.LogStore.Logs {
			promiseReq.LocalVal = append(promiseReq.LocalVal, txn)
		}

	} else {
		// send existing acceptNum and acceptVal
		promiseReq.AcceptNum = conf.CurrVal.BallotNumber
		promiseReq.AcceptVal = conf.CurrVal.Transactions
	}
	SendPromise(ctx, conf, promiseReq)
}

func SendPromise(ctx context.Context, conf *config.Config, req *common.Promise) {
	leaderAddress := utils.MapServerNumberToAddress[req.BallotNum.ServerNumber]
	server, err := conf.Pool.GetServer(leaderAddress)
	if err != nil {
		fmt.Println(err)
	}

	_, err = server.Promise(ctx, req)
	if err != nil {
		fmt.Println(err)
	}
}
