package outbound

import (
	"GolandProjects/apaxos-gautamsardana/server_bob/utils"
	"context"
	"fmt"
	"time"

	common "GolandProjects/apaxos-gautamsardana/api_common"
	"GolandProjects/apaxos-gautamsardana/server_bob/config"
	"GolandProjects/apaxos-gautamsardana/server_bob/logic"
)

func Prepare(ctx context.Context, conf *config.Config) {
	conf.CurrBallot.TermNumber++
	utils.UpdateBallot(conf, conf.CurrBallot.TermNumber, conf.CurrBallot.ServerNumber)
	conf.MajorityHandler = config.NewMajorityHandler(2000 * time.Millisecond)
	fmt.Printf("Server %d: sending prepare with ballot:%v\n", conf.ServerNumber, conf.CurrBallot)

	for _, serverAddress := range conf.ServerAddresses {
		server, err := conf.Pool.GetServer(serverAddress)
		if err != nil {
			fmt.Println(err)
		}
		ballotDetails := &common.Ballot{
			TermNumber:   conf.CurrBallot.TermNumber,
			ServerNumber: conf.ServerNumber,
		}
		go WaitForMajorityPromises(ctx, conf)
		// todo : this will have more parameters here to notify servers what the last committed txn was
		_, err = server.Prepare(ctx, &common.Prepare{BallotNum: ballotDetails})
		if err != nil {
			fmt.Println(err)
		}
	}
}

func WaitForMajorityPromises(ctx context.Context, conf *config.Config) {
	fmt.Printf("Server %d: waiting for promises...\n", conf.ServerNumber)
	select {
	case <-conf.MajorityHandler.MajorityCh:
		fmt.Printf("Server %d: majority promises received\n", conf.ServerNumber)
		logic.AddLocalTxns(conf)
		acceptRequest := &common.Accept{
			BallotNum:       conf.CurrVal.BallotNumber,
			AcceptVal:       conf.CurrVal.Transactions,
			ServerAddresses: conf.CurrVal.ServerAddresses,
		}
		Accept(ctx, conf, acceptRequest)
	case <-time.After(conf.MajorityHandler.Timeout):
		fmt.Printf("Server %d: timed out waiting for promises\n", conf.ServerNumber)
		config.ResetCurrVal(conf)
		conf.MajorityHandler.TimeoutCh <- true
	}
	return
}
