package outbound

import (
	"context"
	"fmt"

	common "GolandProjects/apaxos-gautamsardana/api_common"
	"GolandProjects/apaxos-gautamsardana/server_bob/config"
	"GolandProjects/apaxos-gautamsardana/server_bob/logic"
)

// i am a leader - i got accepted requests from majority followers. Now I need to tell them to commit

func Commit(ctx context.Context, conf *config.Config, req *common.Commit) {
	dbErr := logic.CommitTransaction(ctx, conf, req)
	if dbErr != nil {
		return
	}
	logic.DeleteFromLogs(conf, req.AcceptVal)
	config.ResetCurrVal(conf)
	config.ResetAcceptVal(conf)
	fmt.Printf("Server %d: new txns committed, updated log:%v\n", conf.ServerNumber, conf.LogStore)

	fmt.Printf("Server %d: sending commit request with req: %v\n", conf.ServerNumber, req)
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
