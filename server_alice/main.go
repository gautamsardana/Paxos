package main

import (
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"

	common "GolandProjects/apaxos-gautamsardana/api_common"
	"GolandProjects/apaxos-gautamsardana/server_alice/api"
	"GolandProjects/apaxos-gautamsardana/server_alice/config"
	"GolandProjects/apaxos-gautamsardana/server_alice/storage/datastore"
	"GolandProjects/apaxos-gautamsardana/server_alice/storage/logstore"
	"GolandProjects/apaxos-gautamsardana/server_alice/utils"
	serverPool "GolandProjects/apaxos-gautamsardana/server_pool"
)

// todo: have a flag somewhere so that the servers who are not live will not send back the promise

// todo: when leader gets promises from majority -- the leader's accepted values will also be
// updated because it is implied that the leader already sent the accept message to itself AND
// received the accepted value from itself

func main() {
	conf := config.GetConfig()
	config.SetupDB(conf)
	utils.GetBallot(conf)
	InitiateLogStore(conf)
	InitiateServerPool(conf)
	ListenAndServe(conf)
}

func InitiateLogStore(conf *config.Config) {
	balance, err := datastore.GetBalance(conf.DataStore, conf.ClientName)
	if err != nil {
		fmt.Printf("error trying to fetch balance from datastore, err: %v", err)
		balance = 100
	}
	logStore := logstore.NewLogStore(balance)
	conf.LogStore = logStore
}

func InitiateServerPool(conf *config.Config) {
	pool, err := serverPool.NewServerPool(conf.ServerAddresses)
	if err != nil {
		// todo: change this
		fmt.Println(err)
	}
	conf.Pool = pool
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
