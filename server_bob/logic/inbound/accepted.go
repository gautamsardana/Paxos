package inbound

import (
	"context"
	"fmt"
	"sync"

	common "GolandProjects/apaxos-gautamsardana/api_common"
	"GolandProjects/apaxos-gautamsardana/server_bob/config"
	"GolandProjects/apaxos-gautamsardana/server_bob/utils"
)

// i am the leader - i got the accepted from followers. Now I need to send them commit messages

//todo: need to send commit only if i get accepted from majority

func Accepted(ctx context.Context, conf *config.Config, req *common.Accepted) error {
	fmt.Printf("Server %d: received accepted with req: %v\n", conf.ServerNumber, req)

	if req.BallotNum.TermNumber != conf.CurrBallot.TermNumber ||
		req.BallotNum.ServerNumber != conf.CurrBallot.ServerNumber {
		return fmt.Errorf("invalid ballot")
	}

	var lock sync.Mutex
	lock.Lock()
	defer lock.Unlock()

	if conf.AcceptedServers == nil {
		conf.AcceptedServers = config.NewAcceptedServersInfo()
	}
	if conf.AcceptVal == nil {
		conf.AcceptVal = config.NewAcceptVal()
		conf.AcceptVal.BallotNumber = req.BallotNum
		conf.AcceptVal.Transactions = req.AcceptVal
	}

	conf.AcceptedServers.CurrAcceptedCount++
	conf.AcceptedServers.ServerAddresses = append(conf.AcceptedServers.ServerAddresses,
		utils.MapServerNumberToAddress[req.ServerNumber])

	if conf.AcceptedServers.CurrAcceptedCount >= (conf.ServerTotal/2)+1 {
		conf.MajorityHandler.MajorityCh <- true
	}

	return nil
}
