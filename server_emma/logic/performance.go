package logic

import (
	"context"
	"google.golang.org/protobuf/types/known/durationpb"
	"time"

	common "GolandProjects/apaxos-gautamsardana/api_common"
	"GolandProjects/apaxos-gautamsardana/server_emma/config"
)

func Performance(ctx context.Context, conf *config.Config, req *common.PerformanceRequest) *common.PerformanceResponse {
	var totalLatency time.Duration
	completedTxns := len(conf.LatencyQueue)
	for i := 0; i < completedTxns; i++ {
		totalLatency += conf.LatencyQueue[i]
	}

	var throughput float64
	if totalLatency > 0 {
		throughput = float64(completedTxns) / totalLatency.Seconds()
	}

	resp := &common.PerformanceResponse{
		Latency:    durationpb.New(totalLatency),
		Throughput: float32(throughput),
	}
	return resp
}
