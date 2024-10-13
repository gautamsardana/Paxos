package logic

import (
	"context"
	"fmt"
	"time"

	common "GolandProjects/apaxos-gautamsardana/api_common"
	"GolandProjects/apaxos-gautamsardana/server_alice/config"
	"GolandProjects/apaxos-gautamsardana/server_alice/storage/datastore"
	"GolandProjects/apaxos-gautamsardana/server_alice/utils"
)

func SendPrepare(ctx context.Context, conf *config.Config) {
	conf.CurrBallot.TermNumber++
	utils.UpdateBallot(conf, conf.CurrBallot.TermNumber, conf.CurrBallot.ServerNumber)

	lastCommittedTerm, dbErr := datastore.GetLatestTermNo(conf.DataStore)
	if dbErr != nil {
		return
	}

	ballotDetails := &common.Ballot{
		TermNumber:   conf.CurrBallot.TermNumber,
		ServerNumber: conf.ServerNumber,
	}
	prepareReq := &common.Prepare{BallotNum: ballotDetails, LastCommittedTerm: lastCommittedTerm}

	conf.MajorityHandler = config.NewMajorityHandler(50000 * time.Millisecond)
	fmt.Printf("Server %d: sending prepare with ballot:%v\n", conf.ServerNumber, conf.CurrBallot)
	go WaitForMajorityPromises(ctx, conf)

	for _, serverAddress := range conf.ServerAddresses {
		server, err := conf.Pool.GetServer(serverAddress)
		if err != nil {
			fmt.Println(err)
		}

		// todo : this will have more parameters here to notify servers what the last committed txn was
		_, err = server.Prepare(ctx, prepareReq)
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
		//time.Sleep(10 * time.Millisecond)
		AddLocalTxns(conf)
		acceptRequest := &common.Accept{
			BallotNum:       conf.CurrVal.BallotNumber,
			AcceptVal:       conf.CurrVal.Transactions,
			ServerAddresses: conf.CurrVal.ServerAddresses,
		}
		if len(acceptRequest.AcceptVal) == 0 {
			fmt.Printf("Server %d: no transactions to send\n", conf.ServerNumber)
			config.ResetCurrVal(conf)
			return
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
	fmt.Printf("Server %d: received prepare with request: %v\n", conf.ServerNumber, req)

	isValidPrepare, err := IsValidPrepare(context.Background(), req, conf)
	if err != nil {
		return err
	}
	if !isValidPrepare {
		return fmt.Errorf("not a valid prepare request, exit")
	}
	if !IsValidBallot(req, conf) {
		return fmt.Errorf("not a valid ballot, return error exit")
	}

	// valid prepare request -- proceed
	utils.UpdateBallot(conf, req.BallotNum.TermNumber, req.BallotNum.ServerNumber)

	// need to check for older promises here -- check whether to send old acceptNum,Val or nil with local txns

	fmt.Println(string(conf.ServerNumber)+": sending promise with request: %v", req)
	SendPromise(ctx, conf, &common.Ballot{TermNumber: req.BallotNum.TermNumber, ServerNumber: req.BallotNum.ServerNumber})

	return nil
}

func IsValidBallot(req *common.Prepare, conf *config.Config) bool {
	if req.BallotNum.TermNumber < conf.CurrBallot.TermNumber {
		return false
	}
	return true
}

func IsValidPrepare(ctx context.Context, req *common.Prepare, conf *config.Config) (bool, error) {
	lastCommittedTerm, dbErr := datastore.GetLatestTermNo(conf.DataStore)
	if dbErr != nil {
		return false, dbErr
	}

	// leader is slow, send new txns
	if req.LastCommittedTerm < lastCommittedTerm {
		err := SendSyncResponse(ctx, conf, &common.SyncRequest{
			LastCommittedTerm: req.LastCommittedTerm,
			ServerNo:          req.BallotNum.ServerNumber})
		if err != nil {
			return false, err
		}
		return false, nil
	} else if req.LastCommittedTerm > lastCommittedTerm { // receiver is slow, ask leader for new txns
		server, err := conf.Pool.GetServer(utils.MapServerNumberToAddress[req.BallotNum.ServerNumber])
		if err != nil {
			fmt.Println(err)
		}
		_, err = server.Sync(ctx, &common.SyncRequest{LastCommittedTerm: lastCommittedTerm, ServerNo: conf.ServerNumber})
		if err != nil {
			fmt.Println(err)
			return false, err
		}
		return false, nil
	}
	return true, nil
}
