package main

import (
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"

	common "GolandProjects/apaxos-gautamsardana/api_common"
	"GolandProjects/apaxos-gautamsardana/server_bob/api"
	"GolandProjects/apaxos-gautamsardana/server_bob/config"
	"GolandProjects/apaxos-gautamsardana/server_bob/logic"
	"GolandProjects/apaxos-gautamsardana/server_bob/utils"
)

// todo: when leader gets accepted from majority -- the leader's accepted values will also be
// updated because it is implied that the leader already sent the accept message to itself AND
// received the accepted value from itself

func main() {
	conf := config.GetConfig()
	config.SetupDB(conf)
	config.InitiateConfig(conf)
	utils.GetBallot(conf)

	logic.TransactionWorker(conf)

	ListenAndServe(conf)
}

func ListenAndServe(conf *config.Config) {
	lis, err := net.Listen("tcp", ":"+conf.Port)
	if err != nil {
		panic(err)
	}
	s := grpc.NewServer()
	common.RegisterPaxosServer(s, &api.Server{Config: conf})
	fmt.Printf("gRPC server running on port %v...\n", conf.Port)
	if err := s.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
