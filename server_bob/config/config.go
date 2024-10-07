package config

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

type Config struct {
	Port     string   `json:"port"`
	Database Database `json:"database"`
}

type Database struct {
	DSN      string `json:"dsn"`
	Host     string `json:"host"`
	User     string `json:"user"`
	Password string `json:"password"`
}

func GetConfig() *Config {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	jsonConfig, err := os.ReadFile(homeDir + "/go/src/GolandProjects/CSE535_Project1/Server_B/config/config.json")
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
	db, err := sql.Open("mysql", config.Database.DSN)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	fmt.Println("MySQL Connected!!")

}
