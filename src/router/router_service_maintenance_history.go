package router

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"streelity/v1/model"

	"github.com/gorilla/mux"
	"github.com/nvnamsss/goinf/pipeline"
)

func allMaintenanceHistory(w http.ResponseWriter, req *http.Request) {

}

func getMaintenanceHistories(w http.ResponseWriter, req *http.Request) {
	var res struct {
		Response
		Histories []model.MaintenanceHistory
	}
	res.Status = true

	p := pipeline.NewPipeline()
	stage := pipeline.NewStage(func() (str struct {
		MaintenanceUser string
	}, e error) {
		query := req.URL.Query()
		users, ok := query["mUser"]
		if !ok {
			return str, errors.New("mUser param is missing")
		}

		str.MaintenanceUser = users[0]
		return
	})

	p.First = stage
	res.Error(p.Run())
	if res.Status {

	}
	WriteJson(w, res)
}

func removeMaintenanceHistory(w http.ResponseWriter, req *http.Request) {
	var res Response = Response{Status: true}

	req.ParseForm()
	p := pipeline.NewPipeline()
	stage := pipeline.NewStage(func() (str struct{ Id []int64 }, e error) {
		form := req.PostForm
		ids, ok := form["id"]

		if !ok {
			return str, errors.New("id param is missing")
		}

		for _, id := range ids {
			if id, e := strconv.ParseInt(id, 10, 64); e == nil {
				str.Id = append(str.Id, id)
			}
		}

		return
	})
	p.First = stage
	res.Error(p.Run())
	if res.Status {
		ids := p.GetInt("Id")
		res.Error(model.RemoveMaintenanceHistoriesById(ids...))
	}

	WriteJson(w, res)
}

func HandleMaintenanceHistory(router *mux.Router) {
	log.Println("[Router]", "Handling maintenance history")
	s := router.PathPrefix("/history").Subrouter()
	s.HandleFunc("/", getMaintenanceHistories).Methods("GET")
	s.HandleFunc("/", removeMaintenanceHistory).Methods("DELETE")
}
