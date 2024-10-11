package inbound

import (
	"context"
	"fmt"

	common "GolandProjects/apaxos-gautamsardana/api_common"
	"GolandProjects/apaxos-gautamsardana/server_chucky/config"
	"GolandProjects/apaxos-gautamsardana/server_chucky/logic"
)

//i am a follower - i received this commit message from the leader. I need to commit these messages
// in my db and delete those txns from log

func Commit(ctx context.Context, conf *config.Config, req *common.Commit) error {
	//todo, also check if the same txns are committed in your db already
	fmt.Printf("Server %d: received commit from leader with request: %v\n", conf.ServerNumber, req)

	err := logic.CommitTransaction(ctx, conf, req)
	if err != nil {
		return err
	}

	logic.DeleteFromLogs(conf, req.AcceptVal)
	config.ResetAcceptVal(conf)
	fmt.Printf("Server %d: new txns committed, updated log:%v\n", conf.ServerNumber, conf.LogStore)
	// todo update balance for all

	return nil
}
