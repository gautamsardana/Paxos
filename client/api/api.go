package api

import (
	"GolandProjects/apaxos-gautamsardana/client/config"
	"GolandProjects/apaxos-gautamsardana/client/logic"
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"

	common "GolandProjects/apaxos-gautamsardana/api_common"
)

type Client struct {
	common.UnimplementedPaxosServer
	Config *config.Config
}

func (c *Client) ProcessTxn(ctx context.Context, req *common.TxnSet) (*emptypb.Empty, error) {
	err := logic.ProcessTxn(ctx, req, c.Config)
	if err != nil {
		log.Printf("Error processing txn from load balancer: %v", err)
		return nil, err
	}
	return nil, nil
}

func (c *Client) GetBalance(ctx context.Context, req *common.GetBalanceRequest) (*common.GetBalanceResponse, error) {
	resp, err := logic.GetBalance(ctx, req, c.Config)
	if err != nil {
		log.Printf("Error processing txn from load balancer: %v", err)
		return nil, err
	}
	return resp, nil
}
