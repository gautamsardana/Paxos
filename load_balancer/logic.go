package main

import (
	common "GolandProjects/apaxos-gautamsardana/api_common"
	"context"
	"fmt"
)

func printBalance(client common.PaxosClient, user string) {
	resp, err := client.PrintBalance(context.Background(), &common.GetBalanceRequest{User: user})
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Printf("Balance of %s: %.2f\n", user, resp.Balance)
}

func printDB(client common.PaxosClient, user string) {
	resp, err := client.PrintDB(context.Background(), &common.PrintDBRequest{User: user})
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Printf("DB txns of user %s: %+v\n", user, resp.Txns)
}

func printLogs(client common.PaxosClient, user string) {
	resp, err := client.PrintLogs(context.Background(), &common.PrintLogsRequest{User: user})
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Printf("Logs of user %s: %+v\n", user, resp.Logs)
}

func performance(client common.PaxosClient, user string) {
	resp, err := client.Performance(context.Background(), &common.PerformanceRequest{User: user})
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	latency := resp.Latency.AsDuration()

	fmt.Printf("Total Latency till now: %s\n", latency)
	fmt.Printf("Throughput: %.2f transactions/sec\n", resp.Throughput)
}

func processSet(s *common.TxnSet, client common.PaxosClient) {
	_, err := client.ProcessTxnSet(context.Background(), s)
	if err != nil {
		return
	}
}
