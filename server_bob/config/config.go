package config

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
	"time"

	common "GolandProjects/apaxos-gautamsardana/api_common"
	"GolandProjects/apaxos-gautamsardana/server_bob/storage/logstore"
	pool "GolandProjects/apaxos-gautamsardana/server_pool"
)

const configPath = "/go/src/GolandProjects/apaxos-gautamsardana/server_bob/config/config.json"

type Config struct {
	Port            string  `json:"port"`
	ServerNumber    int32   `json:"server_number"`
	ClientName      string  `json:"client_name"`
	ServerTotal     int     `json:"server_total"`
	DBCreds         DBCreds `json:"db_creds"`
	DataStore       *sql.DB
	LogStore        *logstore.LogStore
	ServerAddresses []string `json:"server_addresses"`
	Pool            *pool.ServerPool
	CurrTxn         *common.ProcessTxnRequest
	CurrBallot      *common.Ballot       // for each server maintaining their ballots
	CurrVal         *CurrValDetails      // for leader getting promise requests
	AcceptVal       *AcceptValDetails    // for follower getting accept requests
	AcceptedServers *AcceptedServersInfo // for leader getting accepted requests
	MajorityHandler *MajorityHandlerDetails
	CurrRetryCount  int
	RetryLimit      int `json:"retry_limit"`
	StartTime       time.Time
}

type DBCreds struct {
	DSN      string `json:"dsn"`
	Host     string `json:"host"`
	User     string `json:"user"`
	Password string `json:"password"`
}

func NewCurrBallot() *common.Ballot {
	return &common.Ballot{}
}

type CurrValDetails struct {
	CurrPromiseCount int
	ServerAddresses  []string
	MaxAcceptVal     *common.Ballot
	BallotNumber     *common.Ballot
	Transactions     []*common.ProcessTxnRequest
}

func NewCurrentVal() *CurrValDetails {
	return &CurrValDetails{CurrPromiseCount: 1}
}

func ResetCurrVal(conf *Config) {
	conf.CurrVal = nil
}

type AcceptValDetails struct {
	//ServersRespNumber int
	//MaxAcceptVal      *common.Ballot
	BallotNumber *common.Ballot
	Transactions []*common.ProcessTxnRequest
}

func NewAcceptVal() *AcceptValDetails {
	return &AcceptValDetails{}
}

func ResetAcceptVal(conf *Config) {
	conf.AcceptVal = nil
}

type AcceptedServersInfo struct {
	CurrAcceptedCount int
	ServerAddresses   []string
}

func NewAcceptedServersInfo() *AcceptedServersInfo {
	return &AcceptedServersInfo{
		CurrAcceptedCount: 1,
		ServerAddresses:   make([]string, 0),
	}
}

type MajorityHandlerDetails struct {
	MajorityCh chan bool
	TimeoutCh  chan bool
	Timeout    time.Duration
}

func NewMajorityHandler(timeout time.Duration) *MajorityHandlerDetails {
	return &MajorityHandlerDetails{
		MajorityCh: make(chan bool),
		TimeoutCh:  make(chan bool),
		Timeout:    timeout,
	}
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
