package inbound

import (
	"context"
	"fmt"

	common "GolandProjects/apaxos-gautamsardana/api_common"
	"GolandProjects/apaxos-gautamsardana/server_alice/config"
	"GolandProjects/apaxos-gautamsardana/server_alice/logic/outbound"
)

//i am a follower - i'm getting this accept from the leader. I will accept this value and
//send accepted in return

// add your conf.acceptVal to whatever request you get

func Accept(ctx context.Context, conf *config.Config, req *common.Accept) error {
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

	outbound.Accepted(ctx, conf, acceptedReq)
	return nil
}
