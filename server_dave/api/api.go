package api

import (
	"context"
	"fmt"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"

	common "GolandProjects/apaxos-gautamsardana/api_common"
	"GolandProjects/apaxos-gautamsardana/server_dave/config"
	"GolandProjects/apaxos-gautamsardana/server_dave/logic"
)

type Server struct {
	common.UnimplementedPaxosServer
	Config *config.Config
}

func (s *Server) EnqueueTxn(ctx context.Context, req *common.TxnRequest) (*emptypb.Empty, error) {
	_ = logic.EnqueueTxn(ctx, req, s.Config)
	return nil, nil
}

func (s *Server) Prepare(ctx context.Context, req *common.Prepare) (*emptypb.Empty, error) {
	err := logic.ReceivePrepare(ctx, s.Config, req)
	if err != nil {
		log.Printf("Error processing prepare from leader: %v", err)
		return nil, err
	}
	return nil, nil
}

func (s *Server) Promise(ctx context.Context, req *common.Promise) (*emptypb.Empty, error) {
	err := logic.ReceivePromise(ctx, s.Config, req)
	if err != nil {
		log.Printf("Error processing promise request from follower: %v", err)
		return nil, err
	}
	return nil, nil
}

func (s *Server) Accept(ctx context.Context, req *common.Accept) (*emptypb.Empty, error) {
	err := logic.ReceiveAccept(ctx, s.Config, req)
	if err != nil {
		log.Printf("Error processing accept request from leader: %v", err)
		return nil, err
	}
	return nil, nil
}

func (s *Server) Accepted(ctx context.Context, req *common.Accepted) (*emptypb.Empty, error) {
	err := logic.ReceiveAccepted(ctx, s.Config, req)
	if err != nil {
		log.Printf("Error processing accepted request from follower: %v", err)
		return nil, err
	}
	return nil, nil
}

func (s *Server) Commit(ctx context.Context, req *common.Commit) (*emptypb.Empty, error) {
	err := logic.ReceiveCommit(ctx, s.Config, req)
	if err != nil {
		log.Printf("Error processing commit request from leader: %v", err)
		return nil, err
	}
	return nil, nil
}

func (s *Server) Sync(ctx context.Context, req *common.SyncRequest) (*emptypb.Empty, error) {
	err := logic.SyncRequest(ctx, s.Config, req)
	if err != nil {
		log.Printf("Error processing sync request from slow follower: %v", err)
		return nil, err
	}
	return nil, nil
}

func (s *Server) IsAlive(ctx context.Context, req *common.IsAliveRequest) (*emptypb.Empty, error) {
	fmt.Printf("Server %d: IsAlive set to %t\n", s.Config.ServerNumber, req.IsAlive)
	s.Config.IsAlive = req.IsAlive
	return nil, nil
}

func (s *Server) PrintBalance(ctx context.Context, _ *common.GetBalanceRequest) (*common.GetBalanceResponse, error) {
	resp, err := logic.PrintBalance(ctx, s.Config)
	if err != nil {
		log.Printf("Error fetching balance: %v", err)
		return nil, err
	}
	return resp, nil
}

func (s *Server) GetServerBalance(ctx context.Context, req *common.GetServerBalanceRequest) (*common.GetServerBalanceResponse, error) {
	resp, err := logic.GetServerBalance(ctx, s.Config, req)
	if err != nil {
		log.Printf("Error fetching balance from other servers: %v", err)
		return nil, err
	}
	return resp, nil
}

func (s *Server) PrintLogs(ctx context.Context, req *common.PrintLogsRequest) (*common.PrintLogsResponse, error) {
	fmt.Printf("Server %d: received PrintLogs request\n", s.Config.ServerNumber)

	return &common.PrintLogsResponse{
		Logs: s.Config.LogStore.Logs,
	}, nil
}

func (s *Server) PrintDB(ctx context.Context, req *common.PrintDBRequest) (*common.PrintDBResponse, error) {
	fmt.Printf("Server %d: received PrintLogs request\n", s.Config.ServerNumber)
	resp, err := logic.PrintDB(ctx, s.Config, req)
	if err != nil {
		log.Printf("Error fetching balance from other servers: %v", err)
		return nil, err
	}
	return resp, nil
}
