package main

import (
	"database/sql"
	"fmt"
	"regexp"
	"strconv"

	"example.com/m/v2/Astar"
	_ "github.com/go-sql-driver/mysql"
	r2 "github.com/golang/geo/r2"
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

	results, err := Db.Query("SELECT * FROM streets")
	sspRegex := regexp.MustCompile(`;`)

	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	for results.Next() {
		var id int64
		var nodesData string
		var nodeIds []int64
		var cost float64
		err = results.Scan(&id, nil, nil, &nodesData, &cost, nil, nil)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}

		splits := sspRegex.Split(nodesData, -1)

		for _, e := range splits {
			i, err := strconv.ParseInt(e, 10, 64)

			if err != nil {
				panic(err)
			}

			nodeIds = append(nodeIds, i)
		}

		Astar.Streets[id] = *Astar.NewStreet(id, nodeIds)
	}

	results, err = Db.Query("SELECT * FROM  nodes")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	for results.Next() {
		var id int64
		var lat float64
		var lon float64
		var streetId int64

		err = results.Scan(&id, &lat, &lon, &streetId)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
		// and then print out the tag's Name attribute
		fmt.Println(id)
		fmt.Println(lat)
		fmt.Println(lon)
		fmt.Println(streetId)

		Astar.Nodes[id] = Astar.Node{Id: id, Location: r2.Point{X: lat, Y: lon}, StreetId: streetId}

		//how to make neighbor?
		//nodes on the same street are neighbors
		//the reference are n^2, too much
		//instead we can go into street and get the nodes out
		//better performance, better memory
	}

	fmt.Println(results)
}
