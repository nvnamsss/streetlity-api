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
	var res struct {
		Response
		Atms []model.Atm
	}

	res.Status = true
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

	parseValueStage := pipeline.NewStage(func() (str struct {
		Lat   float64
		Lon   float64
		Range float64
	}, e error) {
		var (
			latErr   error
			lonErr   error
			rangeErr error
		)

		str.Lat, latErr = strconv.ParseFloat(query["location"][0], 64)
		str.Lon, lonErr = strconv.ParseFloat(query["location"][1], 64)
		str.Range, rangeErr = strconv.ParseFloat(query["range"][0], 64)
		if latErr != nil {
			return str, errors.New("cannot parse location[0] to float")
		}
		if lonErr != nil {
			return str, errors.New("cannot parse location[1] to float")
		}
		if rangeErr != nil {
			return str, errors.New("cannot parse range to float")
		}
		return str, nil
	})

	validateParamsStage.NextStage(parseValueStage)
	pipe.First = validateParamsStage

	res.Error(pipe.Run())

	if res.Status {
		lat := pipe.GetFloat("Lat")[0]
		lon := pipe.GetFloat("Lon")[0]
		max_range := pipe.GetFloat("Range")[0]
		var location r2.Point = r2.Point{X: lat, Y: lon}

		res.Atms = model.AtmsInRange(location, max_range)
	}

	Write(w, res)
}

func addAtm(w http.ResponseWriter, req *http.Request) {
	var res Response

	res.Status = true
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

	validateParamsStage.NextStage(parseValueStage)
	pipe.First = validateParamsStage

	err := pipe.Run()

	if err != nil {
		res.Status = false
		res.Message = err.Error()
	}

	if res.Status {
		var s model.Atm
		lat, _ := strconv.ParseFloat(form["location"][0], 64)
		lon, _ := strconv.ParseFloat(form["location"][1], 64)
		s.Lat = float32(lat)
		s.Lon = float32(lon)

		err := model.AddAtm(s)

		if err != nil {
			res.Status = false
			res.Message = err.Error()
		} else {
			res.Message = "Create new atm is succeed"
		}
	}

	jsonData, jsonErr := json.Marshal(res)
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
