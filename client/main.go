package main

import (
	"context"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"

	common "GolandProjects/apaxos-gautamsardana/api_common"
)

func main() {
	conn, err := grpc.NewClient("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to gRPC server at localhost:50051: %v", err)
	}
	defer conn.Close()
	c := common.NewPaxosClient(conn)

	//ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	//defer cancel()

	msg_id, _ := uuid.NewUUID()
	_, err = c.ProcessTxn(context.Background(), &common.ProcessTxnRequest{MsgID: msg_id.String(), Sender: "Alice", Receiver: "Chuck", Amount: 50})
	if err != nil {
		log.Fatalf("error calling function SayHello: %v", err)
	}
}
