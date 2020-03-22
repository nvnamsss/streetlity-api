package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"

	"example.com/m/v2/Astar"
	_ "github.com/go-sql-driver/mysql"
	r2 "github.com/golang/geo/r2"
)

type Configuration struct {
	Server   string
	Database string
	Username string
	Password string
}

var Config Configuration
var Db *sql.DB

func connect() {
	connectionString := fmt.Sprintf("%s:%s@tcp(%s)/%s",
		Config.Username, Config.Password, Config.Server, Config.Database)
	fmt.Println(connectionString)
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		panic(err.Error())
		// fmt.Println("[Database]", "Error occurred", err)
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
		var generic_info string
		// var nodesData string
		// var nodeIds []int64
		var cost float64
		var oneway bool
		var direction string
		err = results.Scan(&id, &generic_info, &cost, &oneway, &direction)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}

		// splits := sspRegex.Split(nodesData, -1)
		// for _, e := range splits {
		// 	i, err := strconv.ParseInt(e, 10, 64)

		// 	if err != nil {
		// 		panic(err)
		// 	}

		// 	nodeIds = append(nodeIds, i)
		// }

		Astar.Streets[id] = *Astar.NewStreet(id, []int64{})
	}

	results, err = Db.Query("SELECT * FROM  nodes")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	for results.Next() {
		var id int64
		var lat float64
		var lon float64
		var street string
		var streetIds []int64 = []int64{}
		err = results.Scan(&id, &lat, &lon, &street)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}

		splits := sspRegex.Split(street, -1)
		for _, e := range splits {
			i, err := strconv.ParseInt(e, 10, 64)

			if err != nil {
				log.Print("[Data Puller] ", id, " is not belong to any street, it will be by passed")
				continue
			}

			streetIds = append(streetIds, i)
		}

		// and then print out the tag's Name attribute
		// fmt.Println(id)
		// fmt.Println(lat)
		// fmt.Println(lon)
		// fmt.Println(streetIds)

		Astar.Nodes[id] = Astar.Node{Id: id, Location: r2.Point{X: lat, Y: lon}, StreetId: streetIds}

		for _, e := range streetIds {
			str := Astar.Streets[e]
			str.NodeIds = append(str.NodeIds, id)
			Astar.Streets[e] = str
		}

		fmt.Println(Astar.Nodes[id])
		//how to make neighbor?
		//nodes on the same street are neighbors
		//the reference are n^2, too much
		//instead we can go into street and get the nodes out
		//better performance, better memory
	}

	fmt.Println(Astar.Streets)
}

func init() {
	file, _ := os.Open("config/config.json")

	defer file.Close()
	decoder := json.NewDecoder(file)
	Config = Configuration{}

	err := decoder.Decode(&Config)
	fmt.Println("[Config]", Config)

	if err != nil {
		panic(err)
	}

}
