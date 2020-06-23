package rtoilet

import (
	"net/http"
	"streelity/v1/model/toilet"
	"streelity/v1/sres"
	"streelity/v1/stages"

	"github.com/golang/geo/r2"
	"github.com/gorilla/mux"
	"github.com/nvnamsss/goinf/pipeline"
)

func GetUnconfirmed(w http.ResponseWriter, req *http.Request) {
	var res struct {
		sres.Response
		Service toilet.ToiletUcf
	}
	res.Status = true

	p := pipeline.NewPipeline()
	stage := stages.IdValidateStage(req.URL.Query())
	p.First = stage
	res.Error(p.Run())

	if res.Status {
		id := p.GetInt("Id")[0]
		if service, e := toilet.UcfById(id); e != nil {
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
		Services []toilet.ToiletUcf
	}
	res.Status = true
	res.Services = toilet.AllToiletUcfs()
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
		id := p.GetIntFirstOrDefault("Id")
		t := p.GetString("UpvoteType")[0]

		switch t {
		case "Immediately":
			if e := toilet.UpvoteUcfImmediately(id); e != nil {
				res.Error(e)
			}
			break
		default:
			if e := toilet.UpvoteUcf(id); e != nil {
				res.Error(e)
			}
		}

	}

	sres.WriteJson(w, res)
}

func UnconfirmedInRange(w http.ResponseWriter, req *http.Request) {
	var res struct {
		sres.Response
		Services []toilet.ToiletUcf
	}
	res.Status = true

	p := pipeline.NewPipeline()
	stage := stages.InRangeServiceValidateStage(req)
	p.First = stage
	res.Error(p.Run())

	if res.Status {
		location := r2.Point{X: p.GetFloatFirstOrDefault("Lat"), Y: p.GetFloatFirstOrDefault("Lon")}
		r := p.GetFloatFirstOrDefault("Range")
		res.Services = toilet.UcfInRange(location, r)
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
		if e := toilet.DeleteUcf(id); e != nil {
			res.Error(e)
		}
	}
	sres.WriteJson(w, res)
}

func HandleUnconfirmed(router *mux.Router) *mux.Router {
	s := router.PathPrefix("/toilet_ucf").Subrouter()

	s.HandleFunc("/", GetUnconfirmed).Methods("GET")
	s.HandleFunc("/", GetAllUnconfirmed).Methods("GET")
	s.HandleFunc("/", DeleteUnconfirmed).Methods("DELETE")
	s.HandleFunc("/range", UnconfirmedInRange).Methods("GET")
	s.HandleFunc("/upvote", UpvoteUnconfirmed).Methods("POST")

	return s
}
