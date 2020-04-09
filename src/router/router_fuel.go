package router

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"streelity/v1/model"

	"github.com/gorilla/mux"
)

func getFuels(w http.ResponseWriter, req *http.Request) {
	var result struct {
		Status  bool
		Fuels   []model.Fuel
		Message []string
	}
	result.Status = true

	result.Fuels = model.AllFuels()

	log.Println("[GetFuels]", result.Fuels)

	jsonData, jsonErr := json.Marshal(result)

	if jsonErr != nil {
		log.Println(jsonErr)
	}

	w.Write(jsonData)
}

func getFuel(w http.ResponseWriter, req *http.Request) {
	var result struct {
		Status  bool
		Fuel    model.Fuel
		Message []string
	}

	result.Status = true
	result.Message = []string{}
	query := req.URL.Query()

	status, err := model.Auth(query["token"][0])
	if !status {
		result.Status = false
		result.Message = append(result.Message, err.Error())
		data, _ := json.Marshal(result)
		w.Write(data)
	}

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
		result.Fuel = model.FuelById(id)
		log.Println("[GetFuel]", result.Fuel)
	} else {
		log.Println("[GetFuel]", "Request failed")
	}

	jsonData, jsonErr := json.Marshal(result)

	if jsonErr != nil {
		log.Println(jsonErr)
	}
	w.Write(jsonData)
}

func updateFuel(w http.ResponseWriter, req *http.Request) {
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
	s.HandleFunc("/all", getFuels).Methods("GET")
	s.HandleFunc("/update", updateFuel).Methods("POST")
	s.HandleFunc("/id", getFuel).Methods("GET")
}
