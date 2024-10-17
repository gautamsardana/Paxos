package logic

import (
	"context"
	"google.golang.org/protobuf/types/known/durationpb"
	"time"

	common "GolandProjects/apaxos-gautamsardana/api_common"
	"GolandProjects/apaxos-gautamsardana/server_chucky/config"
)

func Performance(ctx context.Context, conf *config.Config, req *common.PerformanceRequest) *common.PerformanceResponse {
	var totalTime time.Duration
	completedTxns := len(conf.LatencyQueue)

	for i := 0; i < completedTxns; i++ {
		totalTime += conf.LatencyQueue[i]
	}

	var avgLatency time.Duration
	if conf.TxnCount > 0 {
		avgLatency = totalTime / time.Duration(completedTxns)
	}

	var throughput float64
	if totalTime > 0 {
		throughput = float64(completedTxns) / totalTime.Seconds()
	}

	resp := &common.PerformanceResponse{
		Latency:    durationpb.New(avgLatency),
		Throughput: float32(throughput),
	}
	return resp
}
