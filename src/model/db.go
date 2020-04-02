package model

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"reflect"

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

func loadConfig() {
	file, fileErr := os.Open("config/config.json")
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

func init() {
	loadConfig()

	connectionString := fmt.Sprintf("%s:%s@tcp(%s)/%s",
		Config.Username, Config.Password, Config.Server, Config.Database)
	log.Println(connectionString)
	db, err := gorm.Open("mysql", connectionString)
	Db = db
	log.Println(reflect.TypeOf(db))

	if err != nil {
		log.Println(err.Error())
		//panic(err)
	}
	log.Println("Hi mom init db")
}
