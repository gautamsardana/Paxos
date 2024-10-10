package outbound

import (
	"context"
	"fmt"

	common "GolandProjects/apaxos-gautamsardana/api_common"
	"GolandProjects/apaxos-gautamsardana/server_bob/config"
)

// i am the leader - i got majority promises, now i need to send them the accept message

func Accept(ctx context.Context, conf *config.Config, req *common.Accept) {
	for _, serverAddress := range req.ServerAddresses {
		server, err := conf.Pool.GetServer(serverAddress)
		if err != nil {
			fmt.Println(err)
		}

		_, err = server.Accept(ctx, req)
		if err != nil {
			fmt.Println(err)
		}
	}
}
