package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

var Db *sql.DB

func connect() {
	db, err := sql.Open("mysql", "user:password@/dbname")
	if err != nil {
		fmt.Println("[Database]", "Error occurred", err)
	}

	Db = db
}

func PrepareData() {
	Db.Exec("Hi mom")
}
