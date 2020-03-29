package router

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"example.com/m/v2/model"

	"github.com/gorilla/mux"
)

func getFuels(w http.ResponseWriter, req *http.Request) {
	var result struct {
		Status  bool
		Fuels   []model.Fuel
		Message []string
	}
	result.Status = true

	var f []model.Fuel
	model.Db.Find(&f)
	result.Fuels = f

	log.Println("[GetFuels]", f)

	jsonData, jsonErr := json.Marshal(result)

	if jsonErr != nil {
		log.Println(jsonErr)
	}

	w.Write(jsonData)
}

func getFuel(w http.ResponseWriter, req *http.Request) {
	var result struct {
		Status  bool
		Fuels   []model.Fuel
		Message []string
	}
	result.Status = true
	result.Fuels = []model.Fuel{}
	result.Message = []string{}
	query := req.URL.Query()

	var id int64
	var idErr error
	log.Println("[GetFuel]", query)
	_, idReady := query["id"]
	if !idReady {
		result.Status = false
		result.Message = append(result.Message, "Id is missing")
	} else {
		id, idErr = strconv.ParseInt(query["id"][0], 10, 64)
		if idErr != nil {
			result.Status = false
			result.Message = append(result.Message, "Id is invalid")
		}
	}

	if result.Status {
		var f []model.Fuel
		model.Db.Find(&f, id)
		result.Fuels = f	
		log.Println("[GetFuel]", f)
	} else {
		log.Println("[GetFuel]", "Request failed")
	}

	jsonData, jsonErr := json.Marshal(result)

	if jsonErr != nil {
		log.Println(jsonErr)
	}
	w.Write(jsonData)
}

func update(w http.ResponseWriter, req *http.Request) {
	var result struct {
		Status  bool
		Message []string
	}
	result.Status = true
	result.Message = []string{}
	req.ParseForm()
	id, idErr := strconv.ParseInt(req.PostFormValue("id"), 10, 64)
	lat, latErr := strconv.ParseFloat(req.PostFormValue("lat"), 64)
	lon, lonErr := strconv.ParseFloat(req.PostFormValue("lon"), 64)

	if idErr != nil {
		result.Status = false
		result.Message = append(result.Message, "Id is invalid")
	}

	if latErr != nil {
		result.Status = false
		result.Message = append(result.Message, "Lat is invalid")
	}

	if lonErr != nil {
		result.Status = false
		result.Message = append(result.Message, "Lon is invalid")
	}

	fmt.Println(id, lat, lon)

	fmt.Println(result)

	if result.Status {
		var f model.Fuel
		if err := model.Db.Where(&model.Fuel{Id: id}).First(&f).Error; err != nil {
			result.Status = false
			result.Message = append(result.Message, err.Error())
		}

	}

	jsonData, jsonErr := json.Marshal(result)
	if jsonErr != nil {
		log.Println(jsonErr)
	}
	fmt.Println(string(jsonData))
	w.Write(jsonData)
}

func HandleFuel(router *mux.Router) {
	log.Println("[Router]", "Handling fuel")
	s := router.PathPrefix("/fuel").Subrouter()
	s.HandleFunc("/getFuels", getFuels).Methods("GET")
	s.HandleFunc("/update", update).Methods("POST")
	s.HandleFunc("/getFuel", getFuel).Methods("GET")
}
