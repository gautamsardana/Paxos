package inbound

import (
	"context"
	"fmt"
	"sync"

	common "GolandProjects/apaxos-gautamsardana/api_common"
	"GolandProjects/apaxos-gautamsardana/server_chucky/config"
	"GolandProjects/apaxos-gautamsardana/server_chucky/logic"
	"GolandProjects/apaxos-gautamsardana/server_chucky/utils"
)

// i am the leader - these are all the promises i got

func Promise(ctx context.Context, conf *config.Config, req *common.Promise) error {
	// todo: add timeout - if no majority promises within timeout - fail
	// todo: add the majority check and if majority arrived, add localtxns at the end
	// todo - leader will also have acceptVal which will be updated when majority is received)

	fmt.Printf("Server %d: received promise with request: %v\n", conf.ServerNumber, req)
	if !req.PromiseAck {
		return fmt.Errorf("promise not acknowledged, request canceled")
	}

	if req.BallotNum.TermNumber != conf.CurrBallot.TermNumber ||
		req.BallotNum.ServerNumber != conf.CurrBallot.ServerNumber {
		fmt.Println(req.BallotNum.TermNumber, req.BallotNum.ServerNumber)
		fmt.Println(conf.CurrBallot.TermNumber, conf.CurrBallot.ServerNumber)
		return fmt.Errorf("ballot number mismatch, request canceled")
	}

	if conf.CurrVal == nil {
		conf.CurrVal = config.NewCurrentVal()
		logic.AddBallot(conf, req)
	}
	//todo : this means that once the values are committed, you need to make currVal nil again

	var lock sync.Mutex
	lock.Lock()
	defer lock.Unlock()

	//accept num and accept val are nil -- just add local txns of the follower
	if req.AcceptNum == nil && req.AcceptVal == nil {
		logic.AddNewTxns(conf, req)
	} else {
		// accept num/val not empty -- update currVal of the leader to the acceptVal from follower
		if conf.CurrVal.MaxAcceptVal == nil || req.AcceptNum.TermNumber > conf.CurrVal.MaxAcceptVal.TermNumber {
			conf.CurrVal.MaxAcceptVal = req.AcceptNum
			conf.CurrVal.Transactions = req.AcceptVal
		}
	}

	conf.CurrVal.CurrPromiseCount++
	conf.CurrVal.ServerAddresses = append(conf.CurrVal.ServerAddresses, utils.MapServerNumberToAddress[req.ServerNumber])
	fmt.Println(conf.CurrVal.ServerAddresses)

	if conf.CurrVal.CurrPromiseCount >= (conf.ServerTotal/2)+1 {
		conf.MajorityHandler.MajorityCh <- true
	}

	return nil
}
