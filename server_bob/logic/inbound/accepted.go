package inbound

import (
	"context"
	"fmt"

	common "GolandProjects/apaxos-gautamsardana/api_common"
	"GolandProjects/apaxos-gautamsardana/server_bob/config"
	"GolandProjects/apaxos-gautamsardana/server_bob/logic/outbound"
	"GolandProjects/apaxos-gautamsardana/server_bob/utils"
)

// i am the leader - i got the accepted from followers. Now I need to send them commit messages

//todo: need to send commit only if i get accepted from majority

func Accepted(ctx context.Context, conf *config.Config, req *common.Accepted) error {
	if req.BallotNum.TermNumber != conf.CurrBallot.TermNumber ||
		req.BallotNum.ServerNumber != conf.CurrBallot.ServerNumber {
		return fmt.Errorf("invalid ballot")
	}
	conf.AcceptedServers.CurrentAcceptedCount++
	conf.AcceptedServers.ServerAddresses = append(conf.AcceptedServers.ServerAddresses,
		utils.MapServerNumberToAddress[req.ServerNumber])

	// if majority received -
	commitReq := &common.Commit{
		BallotNum:       req.BallotNum,
		AcceptVal:       req.AcceptVal,
		ServerAddresses: conf.AcceptedServers.ServerAddresses,
	}
	outbound.Commit(ctx, conf, commitReq)

	return nil
}
