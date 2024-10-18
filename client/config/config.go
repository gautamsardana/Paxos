package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	serverPool "GolandProjects/apaxos-gautamsardana/server_pool"
)

const configPath = "/go/src/GolandProjects/apaxos-gautamsardana/client/config/config.json"

type Config struct {
	Port            string   `json:"port"`
	ServerAddresses []string `json:"server_addresses"`
	Pool            *serverPool.ServerPool
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

func InitiateServerPool(conf *Config) {
	pool, err := serverPool.NewServerPool(conf.ServerAddresses)
	if err != nil {
		//todo: change this
		fmt.Println(err)
	}
	conf.Pool = pool
}
