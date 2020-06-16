package rmaintenance

import (
	"net/http"
	"streelity/v1/model/maintenance"
	"streelity/v1/sres"
	"streelity/v1/stages"

	"github.com/golang/geo/r2"
	"github.com/gorilla/mux"
	"github.com/nvnamsss/goinf/pipeline"
)

func GetAllUnconfirmed(w http.ResponseWriter, req *http.Request) {
	var res struct {
		sres.Response
		Services []maintenance.MaintenanceUcf
	}
	res.Status = true

	res.Services = maintenance.AllUcfs()
	sres.WriteJson(w, res)
}

func UpvoteUnconfirmed(w http.ResponseWriter, req *http.Request) {
	var res sres.Response = sres.Response{Status: true}
	req.ParseForm()
	p := pipeline.NewPipeline()
	stage := stages.IdValidateStage(req.PostForm)
	p.First = stage
	res.Error(p.Run())

	if res.Status {
		id := p.GetIntFirstOrDefault("Id")
		if e := maintenance.UpvoteUcf(id); e != nil {
			res.Error(e)
		}
	}

	sres.WriteJson(w, res)
}

func UnconfirmedInRange(w http.ResponseWriter, req *http.Request) {
	var res struct {
		sres.Response
		Services []maintenance.MaintenanceUcf
	}
	res.Status = true

	p := pipeline.NewPipeline()
	stage := stages.InRangeServiceValidateStage(req)
	p.First = stage
	res.Error(p.Run())

	if res.Status {
		location := r2.Point{X: p.GetFloatFirstOrDefault("Lat"), Y: p.GetFloatFirstOrDefault("Lon")}
		r := p.GetFloatFirstOrDefault("Range")
		res.Services = maintenance.UcfInRange(location, r)
	}

	sres.WriteJson(w, res)
}

func DeleteUnconfirmed(w http.ResponseWriter, req *http.Request) {
	var res sres.Response = sres.Response{Status: true}

	req.ParseForm()
	p := pipeline.NewPipeline()
	stage := stages.IdValidateStage(req.PostForm)
	p.First = stage
	res.Error(p.Run())

	if res.Status {
		id := p.GetIntFirstOrDefault("Id")
		if e := maintenance.DeleteUcf(id); e != nil {
			res.Error(e)
		}
	}
	sres.WriteJson(w, res)
}

func HandleUnconfirmed(router *mux.Router) {
	s := router.PathPrefix("/maintenance_ucf").Subrouter()

	s.HandleFunc("/", GetAllUnconfirmed).Methods("GET")
	s.HandleFunc("/", DeleteUnconfirmed).Methods("DELETE")
	s.HandleFunc("/range", UnconfirmedInRange).Methods("GET")
	s.HandleFunc("/upvote", UpvoteUnconfirmed).Methods("POST")
}
