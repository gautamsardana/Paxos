package logic

import (
	"context"
	"errors"
	"fmt"
	"time"

	common "GolandProjects/apaxos-gautamsardana/api_common"
	"GolandProjects/apaxos-gautamsardana/server_chucky/config"
	"GolandProjects/apaxos-gautamsardana/server_chucky/utils"
)

func SendPrepare(ctx context.Context, conf *config.Config) {
	conf.CurrBallot.TermNumber++
	utils.UpdateBallot(conf, conf.CurrBallot.TermNumber, conf.CurrBallot.ServerNumber)
	conf.MajorityHandler = config.NewMajorityHandler(50000 * time.Millisecond)
	fmt.Printf("Server %d: sending prepare with ballot:%v\n", conf.ServerNumber, conf.CurrBallot)
	go WaitForMajorityPromises(ctx, conf)

	for _, serverAddress := range conf.ServerAddresses {
		server, err := conf.Pool.GetServer(serverAddress)
		if err != nil {
			fmt.Println(err)
		}
		ballotDetails := &common.Ballot{
			TermNumber:   conf.CurrBallot.TermNumber,
			ServerNumber: conf.ServerNumber,
		}
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
		AddLocalTxns(conf)
		acceptRequest := &common.Accept{
			BallotNum:       conf.CurrVal.BallotNumber,
			AcceptVal:       conf.CurrVal.Transactions,
			ServerAddresses: conf.CurrVal.ServerAddresses,
		}
		SendAccept(context.Background(), conf, acceptRequest)
	case <-time.After(conf.MajorityHandler.Timeout):
		fmt.Printf("Server %d: timed out waiting for promises\n", conf.ServerNumber)
		config.ResetCurrVal(conf)
		conf.MajorityHandler.TimeoutCh <- true
	}
	return
}

func ReceivePrepare(ctx context.Context, conf *config.Config, req *common.Prepare) error {
	fmt.Println(string(conf.ServerNumber)+": received prepare with request: %v", req)
	if !isValidBallot(req, conf) {
		return errors.New("invalid ballot")
	}
	// valid prepare request -- process
	utils.UpdateBallot(conf, req.BallotNum.TermNumber, req.BallotNum.ServerNumber)

	// need to check for older promises here -- check whether to send old acceptNum,Val or nil with local txns

	fmt.Println(string(conf.ServerNumber)+": sending promise with request: %v", req)
	SendPromise(ctx, conf, &common.Ballot{TermNumber: req.BallotNum.TermNumber, ServerNumber: req.BallotNum.ServerNumber})

	return nil
}

func isValidBallot(req *common.Prepare, conf *config.Config) bool {
	if req.BallotNum.TermNumber < conf.CurrBallot.TermNumber {
		return false
	}
	return true
}
