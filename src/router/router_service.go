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

func allService(w http.ResponseWriter, res *http.Request) {

}

func serviceInRange(w http.ResponseWriter, res *http.Request) {
	var result struct {
		Status bool
		model.Services
		// Fuels    []model.Fuel
		Message string
	}

	result.Status = true

	query := res.URL.Query()
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

	validateParamsStage.NextStage(parseValueStage)
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
		services := model.ServicesInRange(location, max_range)

		for _, service := range services {
			a, isAtm := service.(model.Atm)
			f, isFuel := service.(model.Fuel)
			t, isToilet := service.(model.Toilet)

			if isAtm {
				result.Atms = append(result.Atms, a)
			}
			if isFuel {
				result.Fuels = append(result.Fuels, f)
			}
			if isToilet {
				result.Toilets = append(result.Toilets, t)
			}
		}
	}

	jsonData, _ := json.Marshal(result)

	w.Write(jsonData)
}

func HandleService(router *mux.Router) {
	log.Println("[Router]", "Handling service")

	s := router.PathPrefix("/service").Subrouter()
	s.HandleFunc("/all", allService).Methods("GET")
	s.HandleFunc("/range", serviceInRange).Methods("GET")
	HandleFuel(s)
}
