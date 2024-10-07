package utils

import (
	"log"
	"os"
	"strconv"
	"strings"

	"GolandProjects/apaxos-gautamsardana/server_alice/config"
)

const ballotFilePath = "/server_alice/ballot.txt"

func GetBallot(conf *config.Config) {
	dir, _ := os.Getwd()
	fileContent, err := os.ReadFile(dir + ballotFilePath)
	if err != nil {
		log.Fatal(err)
	}
	ballotDetails := strings.Split(string(fileContent), ".")
	currTermNumber, err := strconv.Atoi(ballotDetails[0])
	currServerNumber, err := strconv.Atoi(ballotDetails[1])
	if err != nil {
		log.Fatal(err)
	}
	conf.CurrBallot.TermNumber = int32(currTermNumber)
	conf.CurrBallot.ServerNumber = int32(currServerNumber)
}

func UpdateBallot(conf *config.Config, updatedTermNumber, updatedServerNumber int32) {
	dir, _ := os.Getwd()
	updatedFileContent := strconv.Itoa(int(updatedTermNumber)) + "." + strconv.Itoa(int(updatedServerNumber))
	err := os.WriteFile(dir+ballotFilePath, []byte(updatedFileContent), 0777)
	if err != nil {
		log.Fatal(err)
	}
	conf.CurrBallot.TermNumber = updatedTermNumber
	conf.CurrBallot.ServerNumber = updatedServerNumber
}
