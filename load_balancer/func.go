package main

import (
	common "GolandProjects/apaxos-gautamsardana/api_common"
	"context"
	"fmt"
)

func GetBalance(client common.PaxosClient, user string) {
	fmt.Println("Printing balance for user...")
	resp, err := client.GetBalance(context.Background(), &common.GetBalanceRequest{User: user})
	if err != nil {
		fmt.Println("Error:", err)
	}
	fmt.Printf("Balance of %s: %f\n", user, resp.Balance)
}

func printDB(client common.PaxosClient) {
	fmt.Println("Printing database state...")
}

func printLog(client common.PaxosClient) {
	fmt.Println("Printing transaction log...")
}

// Process transactions for a given set
func processSet(s *common.TxnSet, client common.PaxosClient) {
	fmt.Printf("Processing Set ", s)

	_, err := client.ProcessTxn(context.Background(), s)
	if err != nil {
		return
	}

}
