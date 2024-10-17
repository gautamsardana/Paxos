package logic

import (
	"context"
	"fmt"
	"time"

	common "GolandProjects/apaxos-gautamsardana/api_common"
	"GolandProjects/apaxos-gautamsardana/server_dave/config"
	"GolandProjects/apaxos-gautamsardana/server_dave/utils"
)

func SendPrepare(ctx context.Context, conf *config.Config) {
	conf.MajorityAchieved = false
	config.ResetCurrVal(conf)
	config.ResetAcceptVal(conf)
	utils.UpdateBallot(conf, conf.CurrBallot.TermNumber+1, conf.ServerNumber)

	ballotDetails := &common.Ballot{
		TermNumber:   conf.CurrBallot.TermNumber,
		ServerNumber: conf.CurrBallot.ServerNumber,
	}
	prepareReq := &common.Prepare{BallotNum: ballotDetails, LastCommittedTerm: conf.LastCommittedTerm}

	fmt.Printf("Server %d: sending prepare with ballot:%v\n", conf.ServerNumber, conf.CurrBallot)
	go WaitForMajorityPromises(ctx, conf)

	for _, serverAddress := range conf.ServerAddresses {
		fmt.Println("\n---------", serverAddress, "==========\n")
		server, err := conf.Pool.GetServer(serverAddress)
		if err != nil {
			fmt.Println(err)
		}

		_, err = server.Prepare(ctx, prepareReq)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func WaitForMajorityPromises(ctx context.Context, conf *config.Config) {
	fmt.Printf("Server %d: waiting for promises...\n", conf.ServerNumber)

	// Set a timeout duration
	timeout := time.After(300 * time.Millisecond)

	// Collect promises until the timeout
	for {
		select {
		case <-timeout:
			fmt.Printf("Server %d: timeout reached, checking for majority...\n", conf.ServerNumber)

			// Check if the majority has been achieved
			if conf.CurrVal.CurrPromiseCount >= (conf.ServerTotal/2)+1 {
				fmt.Printf("Server %d: majority promises received\n", conf.ServerNumber)
				conf.MajorityAchieved = true
				if conf.CurrVal.MaxAcceptVal.TermNumber == 0 {
					AddLocalTxns(conf)
				}

				acceptRequest := &common.Accept{
					BallotNum:       conf.CurrBallot,
					AcceptVal:       conf.CurrVal.Transactions,
					ServerAddresses: conf.CurrVal.ServerAddresses,
				}
				if len(acceptRequest.AcceptVal) == 0 {
					fmt.Printf("Server %d: no transactions to send\n", conf.ServerNumber)
					return
				}
				SendAccept(ctx, conf, acceptRequest)
			} else {
				fmt.Printf("Server %d: not enough promises received, canceling...\n", conf.ServerNumber)
				config.ResetCurrVal(conf)
			}
			return
		}
	}
}

func ReceivePrepare(ctx context.Context, conf *config.Config, req *common.Prepare) error {
	fmt.Printf("Server %d: received prepare with request: %v\n", conf.ServerNumber, req)
	if !conf.IsAlive {
		return fmt.Errorf("Server %d: server not alive\n", conf.ServerNumber)
	}

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

	utils.UpdateBallot(conf, req.BallotNum.TermNumber, req.BallotNum.ServerNumber)
	SendPromise(ctx, conf, &common.Ballot{TermNumber: req.BallotNum.TermNumber, ServerNumber: req.BallotNum.ServerNumber})

	return nil
}

func IsValidBallot(req *common.Prepare, conf *config.Config) bool {
	if req.BallotNum.TermNumber < conf.CurrBallot.TermNumber {
		err := SendSyncResponse(context.Background(), conf, &common.SyncRequest{
			LastCommittedTerm: req.LastCommittedTerm,
			ServerNo:          req.BallotNum.ServerNumber})
		if err != nil {
			return false
		}
		return false
	}
	return true
}

func IsValidPrepare(ctx context.Context, req *common.Prepare, conf *config.Config) (bool, error) {
	// leader is slow, send new txns
	if req.LastCommittedTerm < conf.LastCommittedTerm {
		err := SendSyncResponse(ctx, conf, &common.SyncRequest{
			LastCommittedTerm: req.LastCommittedTerm,
			ServerNo:          req.BallotNum.ServerNumber})
		if err != nil {
			return false, err
		}
		return false, nil
	} else if req.LastCommittedTerm > conf.LastCommittedTerm { // receiver is slow, ask leader for new txns
		server, err := conf.Pool.GetServer(utils.MapServerNumberToAddress[req.BallotNum.ServerNumber])
		if err != nil {
			fmt.Println(err)
		}
		_, err = server.Sync(ctx, &common.SyncRequest{LastCommittedTerm: conf.LastCommittedTerm, ServerNo: conf.ServerNumber})
		if err != nil {
			fmt.Println(err)
			return false, err
		}
		return false, nil
	}
	return true, nil
}
