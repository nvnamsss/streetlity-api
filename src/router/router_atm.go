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

func getAtms(w http.ResponseWriter, req *http.Request) {

}

func updateAtm(w http.ResponseWriter, req *http.Request) {

}

func getAtmById(w http.ResponseWriter, req *http.Request) {

}
func getAtmInRange(w http.ResponseWriter, req *http.Request) {
	var result struct {
		Status  bool
		Atms    []model.Atm
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

		result.Atms = model.AllAtmsInRange(location, max_range)
	}

	jsonData, jsonErr := json.Marshal(result)
	if jsonErr != nil {
		log.Println(jsonErr)
	}
	w.Write(jsonData)
}

func addAtm(w http.ResponseWriter, req *http.Request) {
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
		var s model.Atm
		lat, _ := strconv.ParseFloat(form["location"][0], 64)
		lon, _ := strconv.ParseFloat(form["location"][1], 64)
		s.Lat = float32(lat)
		s.Lon = float32(lon)

		err := model.AddAtm(s)

		if err != nil {
			result.Status = false
			result.Message = err.Error()
		} else {
			result.Message = "Create new atm is succeed"
		}
	}

	jsonData, jsonErr := json.Marshal(result)
	if jsonErr != nil {
		log.Println(jsonErr)
	}
	w.Write(jsonData)
}

func HandleAtm(router *mux.Router) {
	log.Println("[Router]", "Handling Atm")
	s := router.PathPrefix("/atm").Subrouter()
	s.HandleFunc("/all", getAtms).Methods("GET")
	s.HandleFunc("/update", updateAtm).Methods("POST")
	s.HandleFunc("/id", getFuel).Methods("GET")
	s.HandleFunc("/range", getAtmInRange).Methods("GET")
	s.HandleFunc("/add", addAtm).Methods("POST")
}
