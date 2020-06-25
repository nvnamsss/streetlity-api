package rmaintenance

import (
	"errors"
	"log"
	"net/http"
	"streelity/v1/model/maintenance"
	"streelity/v1/sres"
	"streelity/v1/stages"

	"github.com/gorilla/mux"
	"github.com/nvnamsss/goinf/pipeline"
)

func GetHistory(w http.ResponseWriter, req *http.Request) {
	var res struct {
		sres.Response
		History maintenance.MaintenanceHistory
	}
	res.Status = true

	p := pipeline.NewPipeline()
	stage := stages.IdValidateStage(req.URL.Query())
	p.First = stage
	res.Error(p.Run())

	if res.Status {
		id := p.GetInt("Id")[0]
		if history, e := maintenance.HistoryById(id); e != nil {
			res.Error(e)
		} else {
			res.History = history
		}
	}
	sres.WriteJson(w, res)
}

func GetHistoriesC(w http.ResponseWriter, req *http.Request) {
	var res struct {
		sres.Response
		Histories []maintenance.MaintenanceHistory
	}
	res.Status = true

	p := pipeline.NewPipeline()
	stage := pipeline.NewStage(func() (str struct {
		CommonUser string
	}, e error) {
		query := req.URL.Query()
		users, ok := query["cuser"]
		if !ok {
			return str, errors.New("cuser param is missing")
		}
		str.CommonUser = users[0]
		return
	})
	p.First = stage
	res.Error(p.Run())

	if res.Status {
		c := p.GetString("CommonUser")[0]
		if histories, e := maintenance.HistoriesByCUser(c); e != nil {
			res.Error(e)
		} else {
			res.Histories = histories
		}
	}

	sres.WriteJson(w, res)
}

func GetHistoriesM(w http.ResponseWriter, req *http.Request) {
	var res struct {
		sres.Response
		Histories []maintenance.MaintenanceHistory
	}
	res.Status = true

	p := pipeline.NewPipeline()
	stage := pipeline.NewStage(func() (str struct {
		MaintenanceUser string
	}, e error) {
		query := req.URL.Query()

		users, ok := query["muser"]
		if !ok {
			return str, errors.New("muser param is missing")
		}

		str.MaintenanceUser = users[0]
		return
	})
	p.First = stage
	res.Error(p.Run())

	if res.Status {
		m := p.GetString("MaintenanceUser")[0]
		if histories, e := maintenance.HistoriesByMUser(m); e != nil {
			res.Error(e)
		} else {
			res.Histories = histories
		}
	}

	sres.WriteJson(w, res)
}

func HandleHistory(router *mux.Router) *mux.Router {
	log.Println("[Maintenance]", "Loading history")
	s := router.PathPrefix("/history").Subrouter()
	s.HandleFunc("/", GetHistory).Methods("GET")
	s.HandleFunc("/c", GetHistoriesC).Methods("GET")
	s.HandleFunc("/m", GetHistoriesM).Methods("GET")

	return s
}
