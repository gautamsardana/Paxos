package logic

import (
	"context"
	"fmt"
	"time"

	common "GolandProjects/apaxos-gautamsardana/api_common"
	"GolandProjects/apaxos-gautamsardana/server_alice/config"
)

// i am the leader - i got majority promises, now i need to send them the accept message

func SendAccept(ctx context.Context, conf *config.Config, req *common.Accept) {
	fmt.Printf("Server %d: sending accept with request: %v\n", conf.ServerNumber, req)
	conf.MajorityHandler = config.NewMajorityHandler(50000 * time.Millisecond)
	go WaitForMajorityAccepted(ctx, conf)
	for _, serverAddress := range req.ServerAddresses {
		server, err := conf.Pool.GetServer(serverAddress)
		if err != nil {
			fmt.Println(err)
		}
		_, err = server.Accept(ctx, req)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func WaitForMajorityAccepted(ctx context.Context, conf *config.Config) {
	fmt.Printf("Server %d: waiting for accepted...\n", conf.ServerNumber)
	select {
	case <-conf.MajorityHandler.MajorityCh:
		fmt.Printf("Server %d: majority accepted received\n", conf.ServerNumber)
		commitRequest := &common.Commit{
			BallotNum:       conf.AcceptVal.BallotNumber,
			AcceptVal:       conf.AcceptVal.Transactions,
			ServerAddresses: conf.AcceptedServers.ServerAddresses,
		}
		SendCommit(context.Background(), conf, commitRequest)
		return
	case <-time.After(conf.MajorityHandler.Timeout):
		fmt.Printf("Server %d: timed out waiting for accepted\n", conf.ServerNumber)
		config.ResetCurrVal(conf)
		config.ResetAcceptVal(conf)
		conf.MajorityHandler.TimeoutCh <- true
		return
	}
}

//i am a follower - i'm getting this accept from the leader. I will accept this value and
//send accepted in return

// add your conf.acceptVal to whatever request you get

func ReceiveAccept(ctx context.Context, conf *config.Config, req *common.Accept) error {
	/* todo - any case where this ballot number can be greater than your current ballot?
	what happens when it is greater? Do you update your current ballot?
	*/
	fmt.Printf("Server %d: received accept from leader with request: %v\n", conf.ServerNumber, req)

	if req.BallotNum.TermNumber < conf.CurrBallot.TermNumber {
		return fmt.Errorf("outdated ballot number")
	}

	conf.AcceptVal = &config.AcceptValDetails{
		BallotNumber: req.BallotNum,
		Transactions: req.AcceptVal,
	}

	acceptedReq := &common.Accepted{
		BallotNum:    req.BallotNum,
		AcceptVal:    req.AcceptVal,
		ServerNumber: conf.ServerNumber,
	}

	SendAccepted(ctx, conf, acceptedReq)
	return nil
}
