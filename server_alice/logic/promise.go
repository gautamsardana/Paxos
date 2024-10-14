package logic

import (
	"context"
	"fmt"
	"sync"

	common "GolandProjects/apaxos-gautamsardana/api_common"
	"GolandProjects/apaxos-gautamsardana/server_alice/config"
	"GolandProjects/apaxos-gautamsardana/server_alice/utils"
)

// i am a follower - i got the leader's prepare. I need to send back my promise
// todo - only send this promise if the server is live based on the input

func SendPromise(ctx context.Context, conf *config.Config, ballotNumber *common.Ballot) {
	fmt.Printf("Server %d: received prepare, sending promise\n", conf.ServerNumber)
	promiseReq := &common.Promise{
		PromiseAck:   true,
		ServerNumber: conf.ServerNumber,
		BallotNum:    ballotNumber,
	}

	// send local txns
	if conf.AcceptVal == nil {
		for _, txn := range conf.LogStore.Logs {
			promiseReq.LocalVal = append(promiseReq.LocalVal, txn)
		}
	} else {
		// send existing acceptNum and acceptVal
		promiseReq.AcceptNum = conf.AcceptVal.BallotNumber
		promiseReq.AcceptVal = conf.AcceptVal.Transactions
	}
	SendPromiseToLeader(ctx, conf, promiseReq)
}

func SendPromiseToLeader(ctx context.Context, conf *config.Config, req *common.Promise) {
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

// i am the leader - these are all the promises i got

func ReceivePromise(ctx context.Context, conf *config.Config, req *common.Promise) error {
	fmt.Printf("Server %d: received promise with request: %v\n", conf.ServerNumber, req)
	if !req.PromiseAck {
		return fmt.Errorf("promise not acknowledged, request canceled")
	}

	if req.BallotNum.TermNumber != conf.CurrBallot.TermNumber ||
		req.BallotNum.ServerNumber != conf.CurrBallot.ServerNumber {
		return fmt.Errorf("ballot number mismatch, request canceled")
	}

	if conf.CurrVal == nil {
		conf.CurrVal = config.NewCurrentVal()
		AddBallot(conf, req)
	}
	//todo : this means that once the values are committed, you need to make currVal nil again

	var lock sync.Mutex
	lock.Lock()
	defer lock.Unlock()

	//accept num and accept val are nil -- just add local txns of the follower
	if req.AcceptNum == nil && req.AcceptVal == nil {
		AddNewTxnsToCurrVal(conf, req)
	} else {
		// accept num/val not empty -- update currVal of the leader to the acceptVal from follower
		if conf.CurrVal.MaxAcceptVal == nil || req.AcceptNum.TermNumber > conf.CurrVal.MaxAcceptVal.TermNumber {
			conf.CurrVal.MaxAcceptVal = req.AcceptNum
			conf.CurrVal.Transactions = req.AcceptVal
		}
	}

	conf.CurrVal.CurrPromiseCount++
	conf.CurrVal.ServerAddresses = append(conf.CurrVal.ServerAddresses, utils.MapServerNumberToAddress[req.ServerNumber])

	if conf.CurrVal.CurrPromiseCount >= (conf.ServerTotal/2)+1 {
		conf.MajorityHandler.MajorityCh <- true
	}

	return nil
}
