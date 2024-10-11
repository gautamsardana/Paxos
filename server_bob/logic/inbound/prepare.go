package inbound

import (
	"context"
	"errors"
	"fmt"

	common "GolandProjects/apaxos-gautamsardana/api_common"
	"GolandProjects/apaxos-gautamsardana/server_bob/config"
	"GolandProjects/apaxos-gautamsardana/server_bob/logic/outbound"
	"GolandProjects/apaxos-gautamsardana/server_bob/utils"
)

func Prepare(ctx context.Context, conf *config.Config, req *common.Prepare) error {
	fmt.Println(string(conf.ServerNumber)+": received prepare with request: %v", req)
	if !isValidBallot(req, conf) {
		return errors.New("invalid ballot")
	}
	// valid prepare request -- process
	utils.UpdateBallot(conf, req.BallotNum.TermNumber, req.BallotNum.ServerNumber)

	// need to check for older promises here -- check whether to send old acceptNum,Val or nil with local txns

	fmt.Println(string(conf.ServerNumber)+": sending promise with request: %v", req)
	outbound.Promise(ctx, conf, &common.Ballot{TermNumber: req.BallotNum.TermNumber, ServerNumber: req.BallotNum.ServerNumber})

	return nil
}

func isValidBallot(req *common.Prepare, conf *config.Config) bool {
	if req.BallotNum.TermNumber < conf.CurrBallot.TermNumber {
		return false
	}
	return true
}
