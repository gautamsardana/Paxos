package service

import (
	"GolandProjects/apaxos-gautamsardana/server_alice/config"
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"

	"GolandProjects/apaxos-gautamsardana/server_alice/logic"
)

func ProcessTxn(ctx context.Context, conf *config.Config) (*emptypb.Empty, error) {
	err := logic.ProcessTxn(ctx, conf.CurrTxn, conf)
	if err != nil {
		log.Printf("Error processing txn: %v", err)
		return nil, err
	}
	log.Printf("txn successful!")

	return nil, nil
}
