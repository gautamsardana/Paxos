package outbound

import (
	"context"
	"fmt"

	common "GolandProjects/apaxos-gautamsardana/api_common"
	"GolandProjects/apaxos-gautamsardana/server_alice/config"
	"GolandProjects/apaxos-gautamsardana/server_alice/service"
	"GolandProjects/apaxos-gautamsardana/server_alice/utils"
)

// i am a leader - i got accepted requests from majority followers. Now I need to tell them to commit

func Commit(ctx context.Context, conf *config.Config, req *common.Commit) {
	dbErr := utils.CommitTransaction(ctx, conf, req)
	if dbErr != nil {
		return
	}
	utils.DeleteFromLogs(conf, req.AcceptVal)
	config.ResetCurrVal(conf)
	config.ResetAcceptVal(conf)
	fmt.Printf("Server %d: new txns committed, updated log:%v\n", conf.ServerNumber, conf.LogStore)

	go func() {
		_, err := service.ProcessTxn(ctx, conf)
		if err != nil {
			fmt.Printf("Server %d: error processing txn:%v\n", conf.ServerNumber, conf.CurrTxn)
		}
	}()

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
