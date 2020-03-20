package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

var Db *sql.DB

func connect() {
	db, err := sql.Open("mysql", "ixPwflxQAF:3satzXbjYd@tcp(remotemysql.com)/ixPwflxQAF")
	if err != nil {
		fmt.Println("[Database]", "Error occurred", err)
	}

	Db = db
}

func PrepareData() {
	connect()

	results, err := Db.Query("SELECT * FROM  nodes")

	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	for results.Next() {
		var id int
		var lat float32
		var lon float32
		var street string
        err = results.Scan(&id,&lat,&lon,&street)
        if err != nil {
            panic(err.Error()) // proper error handling instead of panic in your app
        }
                // and then print out the tag's Name attribute
		fmt.Println(id)
		fmt.Println(lat)
		fmt.Println(lon)
		fmt.Println(street)
	}
	
	fmt.Println(results)
}
