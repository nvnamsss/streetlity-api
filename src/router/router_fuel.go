package router

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"streelity/v1/model"
	"streelity/v1/pipeline"

	"github.com/golang/geo/r2"
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

func getFuelInRange(w http.ResponseWriter, req *http.Request) {
	var result struct {
		Status  bool
		Fuels   []model.Fuel
		Message string
	}
	result.Status = false
	query := req.URL.Query()
	var pipe *pipeline.Pipeline = pipeline.NewPipeline()
	validateParamsStage := pipeline.NewStage(func() error {
		location, locationOk := query["location"]
		if !locationOk {
			return errors.New("location param is missing")
		} else {
			if len(location) < 2 {
				return errors.New("location param must have 2 values")
			}
		}

		_, rangeOk := query["range"]
		if !rangeOk {
			return errors.New("range param is missing")
		}

		return nil
	})

	parseValueStage := pipeline.NewStage(func() error {
		_, latErr := strconv.ParseFloat(query["location"][0], 64)

		_, lonErr := strconv.ParseFloat(query["location"][1], 64)
		_, rangeErr := strconv.ParseFloat(query["range"][0], 64)

		if latErr != nil {
			return errors.New("cannot parse location[0] to float")
		}

		if lonErr != nil {
			return errors.New("cannot parse location[1] to float")
		}

		if rangeErr != nil {
			return errors.New("cannot parse range to float")
		}

		return nil
	})

	validateParamsStage.Next(parseValueStage)
	pipe.First = validateParamsStage

	err := pipe.Run()

	if err != nil {
		result.Status = false
		result.Message = err.Error()
	}

	if result.Status {
		lat, _ := strconv.ParseFloat(query["location"][0], 64)
		lon, _ := strconv.ParseFloat(query["location"][1], 64)
		max_range, _ := strconv.ParseFloat(query["range"][0], 64)
		var location r2.Point = r2.Point{X: lat, Y: lon}

		result.Fuels = model.FuelsInRange(location, max_range)
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
		Message string
	}
	result.Status = true
	req.ParseForm()
	form := req.PostForm

	var pipe *pipeline.Pipeline = pipeline.NewPipeline()
	validateParamsStage := pipeline.NewStage(func() error {
		_, idOk := form["id"]
		location, locationOk := form["location"]

		if !idOk {
			return errors.New("id param is missing")
		}

		if !locationOk {
			return errors.New("location param is missing")
		} else {
			if len(location) < 2 {
				return errors.New("location param must have 2 values")
			}
		}

		return nil
	})

	parseValueStage := pipeline.NewStage(func() error {
		_, idErr := strconv.ParseInt(form["id"][0], 10, 64)
		_, latErr := strconv.ParseFloat(form["location"][0], 64)
		_, lonErr := strconv.ParseFloat(form["location"][1], 64)

		if idErr != nil {
			return errors.New("cannot parse id to int")
		}

		if latErr != nil {
			return errors.New("cannot parse location[0] to float")
		}

		if lonErr != nil {
			return errors.New("cannot parse location[1] to float")
		}

		return nil
	})

	validateParamsStage.Next(parseValueStage)
	pipe.First = validateParamsStage

	err := pipe.Run()

	if err != nil {
		result.Status = false
		result.Message = err.Error()
	}

	if result.Status {
		var f model.Fuel
		id, _ := strconv.ParseInt(form["id"][0], 10, 64)
		if err := model.Db.Where(&model.Fuel{Id: id}).First(&f).Error; err != nil {
			result.Status = false
			result.Message = err.Error()
		}

	}

	jsonData, jsonErr := json.Marshal(result)
	if jsonErr != nil {
		log.Println(jsonErr)
	}
	w.Write(jsonData)
}

func addFuel(w http.ResponseWriter, req *http.Request) {
	var result struct {
		Status  bool
		Message string
	}

	result.Status = true
	req.ParseForm()
	form := req.PostForm

	var pipe *pipeline.Pipeline = pipeline.NewPipeline()
	validateParamsStage := pipeline.NewStage(func() error {
		location, locationOk := form["location"]
		if !locationOk {
			return errors.New("location param is missing")
		} else {
			if len(location) < 2 {
				return errors.New("location param must have 2 values")
			}
		}

		return nil
	})

	parseValueStage := pipeline.NewStage(func() error {
		_, latErr := strconv.ParseFloat(form["location"][0], 64)

		_, lonErr := strconv.ParseFloat(form["location"][1], 64)

		if latErr != nil {
			return errors.New("cannot parse location[0] to float")
		}

		if lonErr != nil {
			return errors.New("cannot parse location[1] to float")
		}

		return nil
	})

	validateParamsStage.Next(parseValueStage)
	pipe.First = validateParamsStage

	err := pipe.Run()

	if err != nil {
		result.Status = false
		result.Message = err.Error()
	}

	if result.Status {
		var f model.Fuel
		lat, _ := strconv.ParseFloat(form["location"][0], 64)
		lon, _ := strconv.ParseFloat(form["location"][1], 64)
		f.Lat = float32(lat)
		f.Lon = float32(lon)

		err := model.AddFuel(f)

		if err != nil {
			result.Status = false
			result.Message = err.Error()
		} else {
			result.Message = "Create new fuel is succeed"
		}
	}

	jsonData, jsonErr := json.Marshal(result)
	if jsonErr != nil {
		log.Println(jsonErr)
	}
	w.Write(jsonData)
}

func HandleFuel(router *mux.Router) {
	log.Println("[Router]", "Handling fuel")
	s := router.PathPrefix("/fuel").Subrouter()
	s.HandleFunc("/all", getFuels).Methods("GET")
	s.HandleFunc("/update", updateFuel).Methods("POST")
	s.HandleFunc("/id", getFuel).Methods("GET")
	s.HandleFunc("/range", getFuelInRange).Methods("GET")
	s.HandleFunc("/add", addFuel).Methods("POSt")
}
