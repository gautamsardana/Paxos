package outbound

import (
	"context"
	"fmt"

	common "GolandProjects/apaxos-gautamsardana/api_common"
	"GolandProjects/apaxos-gautamsardana/server_bob/config"
	"GolandProjects/apaxos-gautamsardana/server_bob/utils"
)

// i am a follower, i accepted the value. Now I need to send this accepted back to the
//leader and wait for commit

func Accepted(ctx context.Context, conf *config.Config, req *common.Accepted) {
	leaderAddress := utils.MapServerNumberToAddress[req.BallotNum.ServerNumber]
	server, err := conf.Pool.GetServer(leaderAddress)
	if err != nil {
		fmt.Println(err)
	}

	_, err = server.Accepted(ctx, req)
	if err != nil {
		fmt.Println(err)
	}
}
