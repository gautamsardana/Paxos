package api

import (
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"

	common "GolandProjects/apaxos-gautamsardana/api_common"
	"GolandProjects/apaxos-gautamsardana/server_chucky/config"
	"GolandProjects/apaxos-gautamsardana/server_chucky/logic/inbound"
)

type Server struct {
	common.UnimplementedPaxosServer
	Config *config.Config
}

func (s *Server) ProcessTxn(ctx context.Context, req *common.ProcessTxnRequest) (*emptypb.Empty, error) {
	err := inbound.ProcessTxn(ctx, req, s.Config)
	if err != nil {
		log.Printf("Error processing txn: %v", err)
		return nil, err
	}
	log.Printf("txn successful!")

	return nil, nil
}

func (s *Server) Prepare(ctx context.Context, req *common.Prepare) (*emptypb.Empty, error) {
	err := inbound.Prepare(ctx, s.Config, req)
	if err != nil {
		log.Printf("Error processing prepare from leader: %v", err)
		return nil, err
	}
	log.Printf("txn successful!")

	return nil, nil
}

func (s *Server) Promise(ctx context.Context, req *common.Promise) (*emptypb.Empty, error) {
	err := inbound.Promise(ctx, s.Config, req)
	if err != nil {
		log.Printf("Error processing promise request from follower: %v", err)
		return nil, err
	}
	log.Printf("txn successful!")

	return nil, nil
}

func (s *Server) Accept(ctx context.Context, req *common.Accept) (*emptypb.Empty, error) {
	err := inbound.Accept(ctx, s.Config, req)
	if err != nil {
		log.Printf("Error processing accept request from leader: %v", err)
		return nil, err
	}
	log.Printf("txn successful!")

	return nil, nil
}

func (s *Server) Accepted(ctx context.Context, req *common.Accepted) (*emptypb.Empty, error) {
	err := inbound.Accepted(ctx, s.Config, req)
	if err != nil {
		log.Printf("Error processing accepted request from follower: %v", err)
		return nil, err
	}
	log.Printf("txn successful!")

	return nil, nil
}

func (s *Server) Commit(ctx context.Context, req *common.Commit) (*emptypb.Empty, error) {
	err := inbound.Commit(ctx, s.Config, req)
	if err != nil {
		log.Printf("Error processing commit request from leader: %v", err)
		return nil, err
	}
	log.Printf("txn successful!")

	return nil, nil
}
