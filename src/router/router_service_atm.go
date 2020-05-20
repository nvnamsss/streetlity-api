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
	var res Response = Response{Status: true}
	req.ParseForm()

	var pipe *pipeline.Pipeline = pipeline.NewPipeline()
	validateParamsStage := AddingServiceValidateStage(req)
	parseValueStage := AddingServiceParsingStage(req)
	bankValidateStage := pipeline.NewStage(func() (str struct{ BankId int64 }, e error) {
		form := req.PostForm
		bank, ok := form["bank"]
		if !ok {
			e = errors.New("bank param is missing")
		} else {
			str.BankId, e = strconv.ParseInt(bank[0], 10, 64)
			if e != nil {
				e = errors.New("cannot parse bank to float")
			}
		}

		return
	})

	validateParamsStage.NextStage(parseValueStage)
	parseValueStage.NextStage(bankValidateStage)
	pipe.First = validateParamsStage

	res.Error(pipe.Run())

	if res.Status {
		var s model.AtmUcf
		lat := pipe.GetFloat("Lat")[0]
		lon := pipe.GetFloat("Lon")[0]
		note := pipe.GetString("Note")[0]
		address := pipe.GetString("Address")[0]
		s.Lat = float32(lat)
		s.Lon = float32(lon)
		s.Note = note
		s.Address = address
		s.BankId = pipe.GetInt("BankId")[0]
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
	var res Response = Response{Status: true}

	req.ParseForm()
	form := req.PostForm

	var pipe *pipeline.Pipeline = pipeline.NewPipeline()
	validateParamsStage := pipeline.NewStage(func() (str struct{ Name string }, e error) {
		name, nameOk := form["name"]
		if !nameOk {
			e = errors.New("name param is missing")
		} else {
			str.Name = name[0]
		}

		return
	})

	pipe.First = validateParamsStage
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

func upvoteAtm(w http.ResponseWriter, req *http.Request) {
	var res Response = Response{Status: true}

	req.ParseForm()
	p := pipeline.NewPipeline()
	vStage := pipeline.NewStage(func() (str struct{ Id int64 }, e error) {
		form := req.PostForm
		_, ok := form["id"]
		if !ok {
			e = errors.New("id params is missing")
			return
		}

		str.Id, e = strconv.ParseInt(form["id"][0], 10, 64)
		return
	})
	p.First = vStage
	res.Error(p.Run())

	if res.Status {
		var id int64 = p.GetInt("Id")[0]
		res.Error(model.UpvoteAtmUcf(id))
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
	res.Status = true

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
	s.HandleFunc("/range", getAtmInRange).Methods("GET")
	s.HandleFunc("/bank/all", getBanks).Methods("GET")
	s.HandleFunc("/bank/add", addBank).Methods("POST")

	r := s.PathPrefix("/add").Subrouter()
	r.HandleFunc("", addAtm).Methods("POST")
	r.Use(Authenticate)

	r = s.PathPrefix("/upvote").Subrouter()
	r.HandleFunc("", upvoteAtm).Methods("POST")
	r.Use(Authenticate)

}
