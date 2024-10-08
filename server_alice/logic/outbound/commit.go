package outbound

import (
	"GolandProjects/apaxos-gautamsardana/server_alice/config"
	"context"
)

// i am a leader - i got accepted requests from majority followers. Now I need to tell them to commit

func Commit(ctx context.Context, conf *config.Config) {

}
