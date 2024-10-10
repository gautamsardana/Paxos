package outbound

import (
	"context"
	"fmt"

	common "GolandProjects/apaxos-gautamsardana/api_common"
	"GolandProjects/apaxos-gautamsardana/server_bob/config"
)

// i am a leader - i got accepted requests from majority followers. Now I need to tell them to commit

func Commit(ctx context.Context, conf *config.Config, req *common.Commit) {
	for _, serverAddress := range req.ServerAddresses {
		server, err := conf.Pool.GetServer(serverAddress)
		if err != nil {
			fmt.Println(err)
		}

		_, err = server.Commit(ctx, req)
		if err != nil {
			fmt.Println(err)
		}
	}
}
