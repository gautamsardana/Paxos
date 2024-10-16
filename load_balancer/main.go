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

const inputFilePath = "lab1_Test.csv"

var sets map[int32]*common.TxnSet // Use a map to store sets by their SetNo
var activeServers map[string]bool // A map to keep track of which servers are live
var totalSets int32

// Function to load test cases from a CSV file
func loadCSV(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ','

	// Initialize the map to store sets by their SetNo
	sets = make(map[int32]*common.TxnSet)

	var currentSetNo int32          // To store the current set number
	var currentLiveServers []string // To store the current live servers

	for {
		// Read each row of the CSV file
		setRow, err := reader.Read()
		if err != nil {
			break
		}

		// Check if the set number is provided in this row
		if setRow[0] != "" {
			// Parse the set number
			setNumber, _ := strconv.Atoi(setRow[0])
			currentSetNo = int32(setNumber)
			totalSets = currentSetNo

			// Parse the live servers for the current set
			currentLiveServers = strings.Split(strings.Trim(setRow[2], "[] "), ",")
			for i := range currentLiveServers {
				currentLiveServers[i] = strings.TrimSpace(currentLiveServers[i])
			}

			// Initialize the set if it doesn't exist
			if _, exists := sets[currentSetNo]; !exists {
				sets[currentSetNo] = &common.TxnSet{
					SetNo:       currentSetNo,
					Txns:        []*common.TxnRequest{},
					LiveServers: currentLiveServers,
				}
			}
		}

		// Parse the transactions for the current set
		txnStrings := strings.Split(setRow[1], ";") // Assuming each txn is separated by a semicolon
		for _, txn := range txnStrings {
			txn = strings.Trim(txn, "() ")
			parts := strings.Split(txn, ",")
			if len(parts) != 3 {
				continue
			}
			amount, _ := strconv.ParseFloat(strings.TrimSpace(parts[2]), 32)

			// Check if the currentSetNo already exists in the map
			if sets[currentSetNo] == nil {
				sets[currentSetNo] = &common.TxnSet{
					SetNo: currentSetNo,
					Txns:  []*common.TxnRequest{}, // Initialize the Txns slice to avoid nil panics
				}
			}

			// Append the transaction to the current set
			sets[currentSetNo].Txns = append(sets[currentSetNo].Txns, &common.TxnRequest{
				Sender:   strings.TrimSpace(parts[0]),
				Receiver: strings.TrimSpace(parts[1]),
				Amount:   float32(amount),
			})
		}
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

	var i int32
	for i = 1; i <= totalSets; i++ {
		fmt.Printf("Processing Set %d: Txns: %v LiveServers: %v\n", i, sets[i].Txns, sets[i].LiveServers)

		scanner.Scan() // Wait for user input before processing the next set

		processSet(sets[i], client)
		// Allow user to run additional functions
		for {
			fmt.Println("Type 'next' to process the next set, 'db' to print database, 'log' to print log, or 'balance' to get balance:")
			scanner.Scan()
			input := scanner.Text()
			if input == "next" {
				break
			} else if input == "db" {
				fmt.Println("Which user?")
				scanner.Scan()
				user := scanner.Text()
				printDB(client, user)
			} else if input == "log" {
				fmt.Println("Which user?")
				scanner.Scan()
				user := scanner.Text()
				printLogs(client, user)
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
