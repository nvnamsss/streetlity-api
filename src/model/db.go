package model

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"streelity/v1/event"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

type Configuration struct {
	Server   string
	Database string `json:"dbname"`
	Username string
	Password string
}

type Service interface {
}

var Db *gorm.DB
var Config Configuration
var OnDisconnect *event.Event = event.NewEvent()

func loadConfig(path string) {
	file, fileErr := os.Open(path)
	if fileErr != nil {

		log.Panic(fileErr)
	}

	defer file.Close()
	decoder := json.NewDecoder(file)
	Config = Configuration{}

	err := decoder.Decode(&Config)

	if err != nil {
		log.Panic(err)
	}

}

func connect() {
	connectionString := fmt.Sprintf("%s:%s@tcp(%s)/%s",
		Config.Username, Config.Password, Config.Server, Config.Database)
	db, err := gorm.Open("mysql", connectionString)
	Db = db

	if err != nil {
		OnDisconnect.Invoke()
		log.Println(err.Error())
		// reconnect()
	}
}

func reconnect() {
	timer := time.NewTimer(10 * time.Second)
	<-timer.C

	connect()
}

func himom() {
	fmt.Println("Hi mom I'm subscriber")
}
func init() {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)
	configPath := filepath.Join(filepath.Dir(basepath), "config", "config.json")

	loadConfig(configPath)
	OnDisconnect.Subscribe(reconnect)
	OnDisconnect.Subscribe(himom)
	OnDisconnect.Unsubscribe(himom)
	go connect()

	log.Println("Hi mom init db")
}
