package router

import (
	"encoding/json"
	"errors"
	"fmt"
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
		Response
		Fuels []model.Fuel
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
	var res struct {
		Response
		Fuel model.Fuel
	}

	res.Status = true
	query := req.URL.Query()

	pipe := pipeline.NewPipeline()
	validateParams := pipeline.NewStage(func() (str struct{ Id int64 }, e error) {
		ids, ok := query["id"]
		if ok {
			return str, errors.New("id param is missing")
		}
		var err error
		str.Id, err = strconv.ParseInt(ids[0], 10, 64)

		if err != nil {
			return str, errors.New("cannot parse id to int")
		}

		return str, nil
	})

	pipe.First = validateParams

	res.Error(pipe.Run())

	if res.Status {
		id := pipe.GetInt("Id")[0]
		res.Fuel = model.FuelById(id)
	}

	res.Write(w)
}

//getFuelInRange process the in-range query. the request must provide there
//
// Parameters:
// 	- `location`: X and Y coordinator
// 	- `range` : range to find
//
func getFuelInRange(w http.ResponseWriter, req *http.Request) {
	var res struct {
		Response
		Fuels []model.Fuel
	}

	res.Status = false
	query := req.URL.Query()

	fmt.Println(req.Header.Get("Auth"))

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

		res.Fuels = model.FuelsInRange(location, max_range)
	}

	res.Write(w)
}

func updateFuel(w http.ResponseWriter, req *http.Request) {
	var res Response

	res.Status = true
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

	validateParamsStage.NextStage(parseValueStage)
	pipe.First = validateParamsStage
	res.Error(pipe.Run())

	if res.Status {
		var f model.Fuel
		id, _ := strconv.ParseInt(form["id"][0], 10, 64)
		if err := model.Db.Where(&model.Fuel{Id: id}).First(&f).Error; err != nil {
			res.Status = false
			res.Message = err.Error()
		}

	}

	res.Write(w)
}

func addFuel(w http.ResponseWriter, req *http.Request) {
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

	parseValueStage := pipeline.NewStage(func() (str struct {
		Lat float64
		Lon float64
	}, e error) {
		var (
			latErr error
			lonErr error
		)

		str.Lat, latErr = strconv.ParseFloat(form["location"][0], 64)
		str.Lon, lonErr = strconv.ParseFloat(form["location"][1], 64)
		if latErr != nil {
			return str, errors.New("cannot parse location[0] to float")
		}
		if lonErr != nil {
			return str, errors.New("cannot parse location[1] to float")
		}
		return str, nil
	})

	validateParamsStage.NextStage(parseValueStage)
	pipe.First = validateParamsStage
	res.Error(pipe.Run())

	if res.Status {
		var f model.Fuel
		lat := pipe.GetFloat("Lat")[0]
		lon := pipe.GetFloat("Lon")[0]
		f.Lat = float32(lat)
		f.Lon = float32(lon)

		err := model.AddFuel(f)

		if err != nil {
			res.Status = false
			res.Message = err.Error()
		} else {
			res.Message = "Create new fuel is succeed"
		}
	}

	res.Write(w)
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
