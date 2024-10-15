package main

import (
	common "GolandProjects/apaxos-gautamsardana/api_common"
	"bufio"
	"encoding/csv"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"os"
	"strconv"
	"strings"
)

const inputFilePath = "input - Sheet1.csv"

var sets []*common.TxnSet
var activeServers map[string]bool // A map to keep track of which servers are live

// Function to load test cases from a CSV file
func loadCSV(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ','

	// Skip the header
	_, err = reader.Read()
	if err != nil {
		return err
	}

	for {
		// Read the set number row
		setRow, err := reader.Read()
		if err != nil {
			break
		}
		setNumber, _ := strconv.Atoi(setRow[0])

		// Read the transactions for the current set
		var transactions []*common.TxnRequest
		for _, txn := range strings.Split(setRow[1], "\n") {
			txn = strings.Trim(txn, "()")
			parts := strings.Split(txn, ",")
			amount, _ := strconv.ParseFloat(parts[2], 32)
			transactions = append(transactions, &common.TxnRequest{
				Sender:   parts[0],
				Receiver: parts[1],
				Amount:   float32(amount),
			})
		}
		// Read the live servers for the current set
		liveServers := strings.Split(strings.Trim(setRow[2], "[]"), ",")
		// Store the set
		sets = append(sets, &common.TxnSet{
			SetNo:       int32(setNumber),
			Txns:        transactions,
			LiveServers: liveServers,
		})
	}
	return nil
}

func main() {

	client := InitiateClient()

	err := loadCSV(inputFilePath)
	if err != nil {
		fmt.Println("Error loading CSV:", err)
		return
	}

	// Process sets one by one
	scanner := bufio.NewScanner(os.Stdin)
	for _, set := range sets {
		fmt.Println("Press Enter to process the next set of transactions...")
		scanner.Scan() // Wait for user input before processing the next set

		processSet(set, client)

		// Allow user to run functions while idle
		for {
			fmt.Println("Type 'next' to process the next set, 'printDB' to print database, 'printLog' to print log, balance to get balance:")
			scanner.Scan()
			input := scanner.Text()
			if input == "next" {
				break
			} else if input == "printDB" {
				printDB(client)
			} else if input == "printLog" {
				printLog(client)
			} else if input == "balance" {
				fmt.Println("Which user?")
				scanner.Scan()
				user := scanner.Text()
				GetBalance(client, user)

			} else {
				fmt.Println("Unknown command")
			}
		}
	}

	fmt.Println("All sets processed.")
}

func InitiateClient() common.PaxosClient {
	conn, err := grpc.NewClient("localhost:8085", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil
	}
	client := common.NewPaxosClient(conn)
	return client
}
