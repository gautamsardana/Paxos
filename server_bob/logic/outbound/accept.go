package outbound

import (
	"context"
	"fmt"
	"time"

	common "GolandProjects/apaxos-gautamsardana/api_common"
	"GolandProjects/apaxos-gautamsardana/server_bob/config"
)

// i am the leader - i got majority promises, now i need to send them the accept message

func Accept(ctx context.Context, conf *config.Config, req *common.Accept) {
	fmt.Printf("Server %d: sending accept with request: %v\n", conf.ServerNumber, req)
	conf.MajorityHandler = config.NewMajorityHandler(50000 * time.Millisecond)
	for _, serverAddress := range req.ServerAddresses {
		server, err := conf.Pool.GetServer(serverAddress)
		if err != nil {
			fmt.Println(err)
		}

		go WaitForMajorityAccepted(ctx, conf)
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
		Commit(ctx, conf, commitRequest)
	case <-time.After(conf.MajorityHandler.Timeout):
		fmt.Printf("Server %d: timed out waiting for accepted\n", conf.ServerNumber)
		config.ResetCurrVal(conf)
		config.ResetAcceptVal(conf)
		conf.MajorityHandler.TimeoutCh <- true
	}
	return
}
