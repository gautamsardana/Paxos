package outbound

import (
	"context"
	"fmt"
	"time"

	common "GolandProjects/apaxos-gautamsardana/api_common"
	"GolandProjects/apaxos-gautamsardana/server_chucky/config"
	"GolandProjects/apaxos-gautamsardana/server_chucky/logic"
)

func Prepare(ctx context.Context, conf *config.Config) {
	for _, serverAddress := range conf.ServerAddresses {
		server, err := conf.Pool.GetServer(serverAddress)
		if err != nil {
			fmt.Println(err)
		}
		ballotDetails := &common.Ballot{
			TermNumber:   conf.CurrBallot.TermNumber + 1,
			ServerNumber: conf.ServerNumber,
		}

		// todo : this will have more parameters here to notify servers what the last committed txn was
		_, err = server.Prepare(ctx, &common.Prepare{BallotNum: ballotDetails})
		if err != nil {
			fmt.Println(err)
		}
	}
	conf.MajorityHandler = config.NewMajorityHandler(200 * time.Millisecond)
	go WaitForMajority(ctx, conf)
}

func WaitForMajority(ctx context.Context, conf *config.Config) {
	select {
	case <-conf.MajorityHandler.MajorityCh:
		fmt.Println("Majority promises received, proceeding to accept phase")
		logic.AddLocalTxns(conf)
		acceptRequest := &common.Accept{
			BallotNum:       conf.CurrVal.BallotNumber,
			AcceptVal:       conf.CurrVal.Transactions,
			ServerAddresses: conf.CurrVal.ServerAddresses,
		}
		Accept(ctx, conf, acceptRequest)
	case <-time.After(conf.MajorityHandler.Timeout):
		config.ResetCurrVal(conf)
		conf.MajorityHandler.TimeoutCh <- true
	}
}
