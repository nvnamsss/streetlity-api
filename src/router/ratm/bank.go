package ratm

import (
	"errors"
	"net/http"
	"streelity/v1/model/atm"
	"streelity/v1/sres"

	"github.com/gorilla/mux"
	"github.com/nvnamsss/goinf/pipeline"
)

func GetBanks(w http.ResponseWriter, req *http.Request) {
	var res struct {
		sres.Response
		Banks []atm.Bank
	}
	res.Status = true

	if res.Status {
		res.Banks = atm.AllBanks()
	}

	sres.WriteJson(w, res)
}

func CreateBank(w http.ResponseWriter, req *http.Request) {
	var res struct {
		sres.Response
		Bank atm.Bank
	}

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
		var s atm.Bank
		s.Name = pipe.GetString("Name")[0]
		err := atm.AddBank(s)

		if err != nil {
			res.Status = false
			res.Message = err.Error()
		} else {
			res.Message = "Add new bank successfully"
			res.Bank, _ = atm.BankByName(s.Name)
		}
	}

	sres.WriteJson(w, res)
}

func HandleBank(router *mux.Router) {
	s := router.PathPrefix("/bank").Subrouter()
	s.HandleFunc("/all", GetBanks).Methods("GET")
	s.HandleFunc("/create", CreateBank).Methods("POST")
}
