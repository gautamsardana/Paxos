package outbound

import (
	"context"
	"fmt"

	common "GolandProjects/apaxos-gautamsardana/api_common"
	"GolandProjects/apaxos-gautamsardana/server_alice/config"
)

func Prepare(ctx context.Context, conf *config.Config) {
	for _, serverAddress := range conf.ServerAddresses {
		server, err := conf.Pool.GetServer(serverAddress)
		if err != nil {
			fmt.Println(err)
		}
		ballotDetails := &common.Ballot{
			TermNumber:   conf.CurrBallot.TermNumber + 1,
			ServerNumber: conf.ServerNumber,
		}

		// todo : this will have more parameters here to notify servers what the last committed txn was
		_, err = server.Prepare(ctx, &common.Prepare{BallotNum: ballotDetails})
		if err != nil {
			fmt.Println(err)
		}
	}
}
