package router

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"streelity/v1/model"
	"streelity/v1/pipeline"

	"github.com/golang/geo/r2"
	"github.com/gorilla/mux"
)

func addToilet(w http.ResponseWriter, req *http.Request) {
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
		var f model.Toilet
		lat := pipe.GetFloat("Lat")[0]
		lon := pipe.GetFloat("Lon")[0]
		f.Lat = float32(lat)
		f.Lon = float32(lon)

		err := model.AddToilet(f)

		if err != nil {
			res.Status = false
			res.Message = err.Error()
		} else {
			res.Message = "Create new fuel is succeed"
		}
	}

	Write(w, res)
}

func getAllToilets(w http.ResponseWriter, req *http.Request) {
	var res struct {
		Response
		Toilets []model.Toilet
	}
	res.Status = true

	res.Error(model.Auth(req.Header.Get("Auth")))
	if !res.Status {
		res.Write(w)
		return
	}

	if res.Status {
		res.Toilets = model.AllToilets()
	}

	Write(w, res)
}

func updateToilet(w http.ResponseWriter, req *http.Request) {

}

func getToiletInRange(w http.ResponseWriter, req *http.Request) {
	var res struct {
		Response
		Toilets []model.Toilet
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

		res.Toilets = model.ToiletsInRange(location, max_range)
	}

	Write(w, res)
}

func HandleToilet(router *mux.Router) {
	log.Println("[Router]", "Handling Toilet")
	s := router.PathPrefix("/toilet").Subrouter()
	s.HandleFunc("/all", getAllToilets).Methods("GET")
	s.HandleFunc("/add", addToilet).Methods("POST")
	s.HandleFunc("/update", updateToilet).Methods("POST")
	s.HandleFunc("/range", getToiletInRange).Methods("GET")
}
