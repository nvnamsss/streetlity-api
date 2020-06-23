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

func GetUnconfirmed(w http.ResponseWriter, req *http.Request) {
	var res struct {
		sres.Response
		Service maintenance.MaintenanceUcf
	}
	res.Status = true

	p := pipeline.NewPipeline()
	stage := stages.IdValidateStage(req.URL.Query())
	p.First = stage
	res.Error(p.Run())

	if res.Status {
		id := p.GetInt("Id")[0]
		if service, e := maintenance.UcfById(id); e != nil {
			res.Error(e)
		} else {
			res.Service = service
		}
	}

	sres.WriteJson(w, res)
}

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
	type_stage := stages.UpvoteTypeStage(req)
	stage.NextStage(type_stage)
	p.First = stage
	res.Error(p.Run())

	if res.Status {
		if res.Status {
			id := p.GetIntFirstOrDefault("Id")
			t := p.GetString("UpvoteType")[0]

			switch t {
			case "Immediately":
				if e := maintenance.UpvoteUcfImmediately(id); e != nil {
					res.Error(e)
				}
				break
			default:
				if e := maintenance.UpvoteUcf(id); e != nil {
					res.Error(e)
				}
			}

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

func HandleUnconfirmed(router *mux.Router) *mux.Router {
	s := router.PathPrefix("/maintenance_ucf").Subrouter()

	s.HandleFunc("/", GetUnconfirmed).Methods("GET")
	s.HandleFunc("/", UpdateReview).Methods("POST")
	s.HandleFunc("/", DeleteUnconfirmed).Methods("DELETE")
	s.HandleFunc("/all", GetAllUnconfirmed).Methods("GET")
	s.HandleFunc("/range", UnconfirmedInRange).Methods("GET")
	s.HandleFunc("/upvote", UpvoteUnconfirmed).Methods("POST")

	return s
}
