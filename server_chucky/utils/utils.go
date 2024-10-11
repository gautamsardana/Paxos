package utils

import (
	"log"
	"os"
	"strconv"
	"strings"

	"GolandProjects/apaxos-gautamsardana/server_chucky/config"
)

const ballotFilePath = "/go/src/GolandProjects/apaxos-gautamsardana/server_chucky/ballot.txt"

var MapServerNumberToAddress = map[int32]string{
	1: "localhost:8080",
	2: "localhost:8081",
	3: "localhost:8082",
	4: "localhost:8083",
	5: "localhost:8084",
}

func GetBallot(conf *config.Config) {
	homeDir, err := os.UserHomeDir()
	fileContent, err := os.ReadFile(homeDir + ballotFilePath)
	if err != nil {
		log.Fatal(err)
	}
	conf.CurrBallot.TermNumber, conf.CurrBallot.ServerNumber = GetTermAndServerNumber(string(fileContent))
}

func GetTermAndServerNumber(ballot string) (int32, int32) {
	ballotDetails := strings.Split(ballot, ".")
	currTermNumber, err := strconv.Atoi(ballotDetails[0])
	currServerNumber, err := strconv.Atoi(ballotDetails[1])
	if err != nil {
		log.Fatal(err)
	}
	return int32(currTermNumber), int32(currServerNumber)
}

func UpdateBallot(conf *config.Config, updatedTermNumber, updatedServerNumber int32) {
	homeDir, _ := os.UserHomeDir()
	updatedFileContent := strconv.Itoa(int(updatedTermNumber)) + "." + strconv.Itoa(int(updatedServerNumber))
	err := os.WriteFile(homeDir+ballotFilePath, []byte(updatedFileContent), 0777)
	if err != nil {
		log.Fatal(err)
	}
	conf.CurrBallot.TermNumber = updatedTermNumber
	conf.CurrBallot.ServerNumber = updatedServerNumber
}
