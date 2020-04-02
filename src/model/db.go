package model

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"runtime"

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

func init() {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)
	configPath := filepath.Join(filepath.Dir(basepath), "config", "config.json")

	loadConfig(configPath)

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
