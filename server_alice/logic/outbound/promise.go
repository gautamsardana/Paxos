package outbound

import (
	"context"
	"fmt"

	common "GolandProjects/apaxos-gautamsardana/api_common"
	"GolandProjects/apaxos-gautamsardana/server_alice/config"
	"GolandProjects/apaxos-gautamsardana/server_alice/utils"
)

// i am a follower - i got the leader's prepare. I need to send back my promise
// todo - only send this promise if the server is live based on the input

func Promise(ctx context.Context, conf *config.Config, ballotNumber *common.Ballot) {
	fmt.Printf("Server %d: received prepare, sending promise\n", conf.ServerNumber)
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
		promiseReq.AcceptNum = conf.AcceptVal.BallotNumber
		promiseReq.AcceptVal = conf.AcceptVal.Transactions
	}
	SendPromise(ctx, conf, promiseReq)
}

func SendPromise(ctx context.Context, conf *config.Config, req *common.Promise) {
	leaderAddress := utils.MapServerNumberToAddress[req.BallotNum.ServerNumber]
	fmt.Println(string(conf.ServerNumber)+"sending promise to address: %v, request: %v ", leaderAddress, req)
	server, err := conf.Pool.GetServer(leaderAddress)
	if err != nil {
		fmt.Println(err)
	}

	_, err = server.Promise(ctx, req)
	if err != nil {
		fmt.Println(err)
	}
}
