package server_pool

import (
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	common "GolandProjects/apaxos-gautamsardana/api_common"
)

type ServerPool struct {
	servers map[string]common.PaxosClient
}

func NewServerPool(serverAddresses []string) (*ServerPool, error) {
	pool := &ServerPool{
		servers: make(map[string]common.PaxosClient),
	}

	for _, addr := range serverAddresses {
		conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			return nil, fmt.Errorf("failed to connect to server %s: %w", addr, err)
		}

		client := common.NewPaxosClient(conn)
		pool.servers[addr] = client
	}
	return pool, nil
}

func (sp *ServerPool) GetServer(addr string) (common.PaxosClient, error) {
	client, ok := sp.servers[addr]
	if !ok {
		return nil, fmt.Errorf("no server found for address: %s", addr)
	}
	return client, nil
}
