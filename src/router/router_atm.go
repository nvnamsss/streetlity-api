package router

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"streelity/v1/model"

	"github.com/golang/geo/r2"
	"github.com/gorilla/mux"
	"github.com/nvnamsss/goinf/pipeline"
)

/*AUTH REQUIRED*/
func updateAtm(w http.ResponseWriter, req *http.Request) {

}

func addAtm(w http.ResponseWriter, req *http.Request) {
	var res Response

	req.ParseForm()
	form := req.PostForm

	var pipe *pipeline.Pipeline = pipeline.NewPipeline()
	authStage := AuthStage(req)
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

	authStage.NextStage(validateParamsStage)
	validateParamsStage.NextStage(parseValueStage)
	pipe.First = authStage

	res.Error(pipe.Run())

	if res.Status {
		var s model.AtmUcf
		lat, _ := strconv.ParseFloat(form["location"][0], 64)
		lon, _ := strconv.ParseFloat(form["location"][1], 64)
		s.Lat = float32(lat)
		s.Lon = float32(lon)

		err := model.AddAtmUcf(s)

		if err != nil {
			res.Status = false
			res.Message = err.Error()
		} else {
			res.Message = "Create new atm is succeed"
		}
	}

	WriteJson(w, res)
}

func addBank(w http.ResponseWriter, req *http.Request) {
	var res Response
	res.Status = true

	req.ParseForm()
	form := req.PostForm

	var pipe *pipeline.Pipeline = pipeline.NewPipeline()
	authStage := AuthStage(req)
	validateParamsStage := pipeline.NewStage(func() (str struct{ Name string }, e error) {
		name, nameOk := form["name"]
		if !nameOk {
			e = errors.New("name param is missing")
		} else {
			str.Name = name[0]
		}

		return
	})

	authStage.NextStage(validateParamsStage)
	pipe.First = authStage
	res.Error(pipe.Run())

	if res.Status {
		var s model.Bank
		s.Name = pipe.GetString("Name")[0]
		err := model.AddBank(s)

		if err != nil {
			res.Status = false
			res.Message = err.Error()
		} else {
			res.Message = "Add new bank successfully"
		}
	}

	WriteJson(w, res)
}

/*NON-AUTH REQUIRED*/

func getAtms(w http.ResponseWriter, req *http.Request) {
	var res struct {
		Response
		Atms []model.Atm
	}
	res.Status = true

	res.Error(model.Auth(req.Header.Get("Auth")))
	if !res.Status {
		res.Write(w)
		return
	}

	if res.Status {
		res.Atms = model.AllAtms()
	}

	WriteJson(w, res)
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

	WriteJson(w, res)
}

func getBanks(w http.ResponseWriter, req *http.Request) {
	var res struct {
		Response
		Banks []model.Bank
	}

	if res.Status {
		res.Banks = model.AllBanks()
	}

	WriteJson(w, res)
}

func HandleAtm(router *mux.Router) {
	log.Println("[Router]", "Handling Atm")
	s := router.PathPrefix("/atm").Subrouter()
	s.HandleFunc("/all", getAtms).Methods("GET")
	s.HandleFunc("/update", updateAtm).Methods("POST")
	s.HandleFunc("/id", getFuel).Methods("GET")
	s.HandleFunc("/range", getAtmInRange).Methods("GET")
	s.HandleFunc("/add", addAtm).Methods("POST")
	s.HandleFunc("/bank/all", getBanks).Methods("GET")
	s.HandleFunc("/bank/add", addBank).Methods("POST")
}
