package main

import (
	"context"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"time"

	common "GolandProjects/apaxos-gautamsardana/api_common"
)

func main() {
	conn, err := grpc.Dial("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to gRPC server at localhost:8080: %v", err)
	}
	defer conn.Close()

	c := common.NewPaxosClient(conn)

	// Generate a new UUID
	msg_id, err := uuid.NewUUID()
	if err != nil {
		log.Fatalf("failed to generate UUID: %v", err)
	}

	// Use a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Call the ProcessTxn RPC method
	_, err = c.ProcessTxn(ctx, &common.ProcessTxnRequest{
		MsgID:    msg_id.String(),
		Sender:   "Alice",
		Receiver: "Chucky",
		Amount:   20,
	})
	if err != nil {
		log.Fatalf("error calling function ProcessTxn: %v", err)
	}
}
