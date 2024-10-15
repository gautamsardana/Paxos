package config

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
	"sync"
	"time"

	common "GolandProjects/apaxos-gautamsardana/api_common"
	"GolandProjects/apaxos-gautamsardana/server_chucky/storage/datastore"
	"GolandProjects/apaxos-gautamsardana/server_chucky/storage/logstore"
	serverPool "GolandProjects/apaxos-gautamsardana/server_pool"
)

const configPath = "/go/src/GolandProjects/apaxos-gautamsardana/server_chucky/config/config.json"

type Config struct {
	Port             string  `json:"port"`
	ServerNumber     int32   `json:"server_number"`
	ClientName       string  `json:"client_name"`
	ServerTotal      int     `json:"server_total"`
	DBCreds          DBCreds `json:"db_creds"`
	DataStore        *sql.DB
	LogStore         *logstore.LogStore
	ServerAddresses  []string `json:"server_addresses"`
	Pool             *serverPool.ServerPool
	CurrBallot       *common.Ballot       // for each server maintaining their ballots
	CurrVal          *CurrValDetails      // for leader getting promise requests
	AcceptVal        *AcceptValDetails    // for follower getting accept requests
	AcceptedServers  *AcceptedServersInfo // for leader getting accepted requests
	MajorityHandler  *MajorityHandlerDetails
	MajorityAchieved bool
	TxnQueue         []*common.TxnRequest
	QueueMutex       sync.Mutex
	Balance          float32
	StartTime        time.Time
}

func InitiateConfig(conf *Config) {
	InitiateCurrVal(conf)
	InitiateTxnQueue(conf)
	InitiateLogStore(conf)
	InitiateBalance(conf)
	InitiateServerPool(conf)
}

func InitiateCurrVal(conf *Config) {
	conf.CurrVal = &CurrValDetails{
		CurrPromiseCount: 1,
		ServerAddresses:  make([]string, 0),
		MaxAcceptVal:     &common.Ballot{},
		Transactions:     make([]*common.TxnRequest, 0),
	}
}

func ResetCurrVal(conf *Config) {
	conf.CurrVal.CurrPromiseCount = 1
	conf.CurrVal.ServerAddresses = make([]string, 0)
	conf.CurrVal.MaxAcceptVal = &common.Ballot{}
	conf.CurrVal.Transactions = make([]*common.TxnRequest, 0)
}

func InitiateTxnQueue(conf *Config) {
	conf.TxnQueue = make([]*common.TxnRequest, 0)
}

func InitiateLogStore(conf *Config) {
	logStore := logstore.NewLogStore()
	conf.LogStore = logStore
}

func InitiateBalance(conf *Config) {
	balance, err := datastore.GetBalance(conf.DataStore, conf.ClientName)
	if err != nil {
		fmt.Printf("error trying to fetch balance from datastore, err: %v", err)
		balance = 100
	}
	conf.Balance = balance
}

func InitiateServerPool(conf *Config) {
	pool, err := serverPool.NewServerPool(conf.ServerAddresses)
	if err != nil {
		// todo: change this
		fmt.Println(err)
	}
	conf.Pool = pool
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
	Transactions     []*common.TxnRequest
}

type AcceptValDetails struct {
	//ServersRespNumber int
	//MaxAcceptVal      *common.Ballot
	BallotNumber *common.Ballot
	Transactions []*common.TxnRequest
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

func NewAcceptedServersInfo(conf *Config) {
	conf.AcceptedServers = &AcceptedServersInfo{
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
