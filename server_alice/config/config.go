package config

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"

	common "GolandProjects/apaxos-gautamsardana/api_common"
	"GolandProjects/apaxos-gautamsardana/server_alice/storage/logstore"
	pool "GolandProjects/apaxos-gautamsardana/server_pool"
)

const configPath = "/go/src/GolandProjects/CSE535_Project1/server_alice/config/config.json"

type Config struct {
	Port            string  `json:"port"`
	ServerNumber    int32   `json:"server_number"`
	ClientName      string  `json:"client_name"`
	ServerTotal     int32   `json:"server_total"`
	DBCreds         DBCreds `json:"db_creds"`
	DataStore       *sql.DB
	LogStore        *logstore.LogStore
	ServerAddresses []string `json:"server_addresses"`
	Pool            *pool.ServerPool
	CurrBallot      BallotDetails   // for each server maintaining their ballots
	CurrVal         *CurrValDetails // for leader getting promise requests
	//CurrPromiseSeq  int                  `json:"curr_promise_seq"` // for follower to maintain the txns which are promised
	AcceptVal       *AcceptValDetails    // for follower getting accept requests
	AcceptedServers *AcceptedServersInfo // for leader getting accepted requests
}

type DBCreds struct {
	DSN      string `json:"dsn"`
	Host     string `json:"host"`
	User     string `json:"user"`
	Password string `json:"password"`
}

type BallotDetails struct {
	TermNumber   int32
	ServerNumber int32
}

type CurrValDetails struct {
	CurrPromiseCount int
	ServerAddresses  []string
	MaxAcceptVal     *common.Ballot
	BallotNumber     *common.Ballot
	Transactions     []*common.ProcessTxnRequest
}

func CurrentValConstructor() *CurrValDetails {
	return &CurrValDetails{CurrPromiseCount: 1}
}

func ResetCurrVal(conf *Config) {
	conf.CurrVal = nil
}

type AcceptValDetails struct {
	ServersRespNumber int
	MaxAcceptVal      *common.Ballot
	BallotNumber      *common.Ballot
	Transactions      []*common.ProcessTxnRequest
}

func ResetAcceptValDetails(conf *Config) {
	conf.AcceptVal = nil
}

type AcceptedServersInfo struct {
	CurrentAcceptedCount int
	ServerAddresses      []string
}

func GetConfig() *Config {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	jsonConfig, err := os.ReadFile(homeDir + configPath)
	if err != nil {
		log.Fatal(err)
	}
	conf := &Config{}
	if err = json.Unmarshal(jsonConfig, conf); err != nil {
		log.Fatal(err)
	}
	return conf
}

func SetupDB(config *Config) {
	db, err := sql.Open("mysql", config.DBCreds.DSN)
	if err != nil {
		log.Fatal(err)
	}
	config.DataStore = db
	//defer db.Close()
	fmt.Println("MySQL Connected!!")
}
