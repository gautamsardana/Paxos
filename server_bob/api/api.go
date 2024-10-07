package api

import (
	"context"
	"fmt"
	"google.golang.org/protobuf/types/known/emptypb"

	common "GolandProjects/apaxos-gautamsardana/api_common"
	"GolandProjects/apaxos-gautamsardana/server_bob/config"
)

type Server struct {
	common.UnimplementedPaxosServer
	Config *config.Config
}

func (s *Server) ProcessTxn(ctx context.Context, req *common.ProcessTxnRequest) (*emptypb.Empty, error) {
	fmt.Println("potty reached")
	fmt.Println(req)
	return nil, nil
}
