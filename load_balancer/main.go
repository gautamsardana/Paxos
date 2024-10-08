package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const inputFilePath = "/input.csv"

type Transaction struct {
	Sender   string
	Receiver string
	Amount   int
}

type Set struct {
	SetNumber    int
	Transactions []Transaction
	LiveServers  []string
}

var sets []Set
var activeServers map[string]bool // A map to keep track of which servers are live

// Simulated functions
func printDB() {
	fmt.Println("Printing database state...")
}

func printLog() {
	fmt.Println("Printing transaction log...")
}

// Process transactions for a given set
func processSet(s Set) {
	fmt.Printf("Processing Set %d\n", s.SetNumber)

	// Mark live servers for the current set
	activeServers = make(map[string]bool)
	for _, server := range s.LiveServers {
		activeServers[server] = true
	}

	// Process each transaction
	for _, t := range s.Transactions {
		if activeServers[t.Sender] && activeServers[t.Receiver] {
			fmt.Printf("Processed Transaction: %s -> %s : %d\n", t.Sender, t.Receiver, t.Amount)
		} else {
			fmt.Printf("Skipped Transaction: %s -> %s : %d (Inactive server)\n", t.Sender, t.Receiver, t.Amount)
		}
	}

	// Simulate idle state after processing the set
	fmt.Println("All servers are idle. You can now use commands like printDB or printLog.")
}

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
		var transactions []Transaction
		for _, txn := range strings.Split(setRow[1], " ") {
			parts := strings.Split(txn, ",")
			amount, _ := strconv.Atoi(parts[2])
			transactions = append(transactions, Transaction{
				Sender:   parts[0],
				Receiver: parts[1],
				Amount:   amount,
			})
		}

		// Read the live servers for the current set
		liveServers := strings.Split(strings.Trim(setRow[2], "[]"), ",")

		// Store the set
		sets = append(sets, Set{
			SetNumber:    setNumber,
			Transactions: transactions,
			LiveServers:  liveServers,
		})
	}

	return nil
}

func main() {
	// Load the CSV file
	err := loadCSV("input - Sheet1.csv")
	if err != nil {
		fmt.Println("Error loading CSV:", err)
		return
	}

	// Process sets one by one
	scanner := bufio.NewScanner(os.Stdin)
	for _, set := range sets {
		fmt.Println("Press Enter to process the next set of transactions...")
		scanner.Scan() // Wait for user input before processing the next set

		processSet(set)

		// Allow user to run functions while idle
		for {
			fmt.Println("Type 'next' to process the next set, 'printDB' to print database, 'printLog' to print log:")
			scanner.Scan()
			input := scanner.Text()
			if input == "next" {
				break
			} else if input == "printDB" {
				printDB()
			} else if input == "printLog" {
				printLog()
			} else {
				fmt.Println("Unknown command")
			}
		}
	}

	fmt.Println("All sets processed.")
}
