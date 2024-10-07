package outbound

import (
	common "GolandProjects/apaxos-gautamsardana/api_common"
	"GolandProjects/apaxos-gautamsardana/server_alice/storage"
)

type Pr struct {
	transactions storage.Transaction
}

func Promise(ack bool, ballotNumber, acceptNum *common.Ballot) {
	if !ack {
		// oops
	}

}
