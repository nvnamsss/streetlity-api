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
	stage := stages.QueryServiceValidateStage(req)
	p.First = stage
	res.Error(p.Run())

	if res.Status {
		c := p.GetIntFirstOrDefault("Case")

		switch c {
		case 1:
			id := p.GetInt("Id")[0]
			if service, e := toilet.UcfById(id); e != nil {
				res.Error(e)
			} else {
				res.Service = service
			}
			break
		case 2:
			lat := p.GetFloat("Lat")[0]
			lon := p.GetFloat("Lon")[0]
			if service, e := toilet.UcfByLocation(lat, lon); e != nil {
				res.Error(e)
			} else {
				res.Service = service
			}
			break
		case 3:
			address := p.GetString("Address")[0]
			if service, e := toilet.UcfByAddress(address); e != nil {
				res.Error(e)
			} else {
				res.Service = service
			}
			break
		}
	}

	sres.WriteJson(w, res)
}

func GetUnconfirmeds(w http.ResponseWriter, req *http.Request) {
	var res struct {
		sres.Response
		Services []toilet.ToiletUcf
	}
	res.Status = true

	p := pipeline.NewPipeline()
	stage := stages.QueryServicesValidateStage(req)
	p.First = stage
	res.Error(p.Run())

	if res.Status {
		address := p.GetString("Address")[0]
		if services, e := toilet.UcfsByAddress(address); e != nil {
			res.Error(e)
		} else {
			res.Services = services
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
	stage := stages.UpvoteValidateStage(req)
	p.First = stage
	res.Error(p.Run())

	if res.Status {
		id := p.GetInt("ServiceId")[0]
		t := p.GetStringFirstOrDefault("UpvoteType")

		switch t {
		case "Immediately":
			if e := toilet.UpvoteServiceImmediately(id); e != nil {
				res.Error(e)
			}
			break
		default:
			if e := toilet.UpvoteService(id); e != nil {
				res.Error(e)
			}
		}

	}

	sres.WriteJson(w, res)
}

func DownvoteService(w http.ResponseWriter, req *http.Request) {
	var res sres.Response = sres.Response{Status: true, Message: "Downvote successfully"}

	req.ParseForm()
	p := pipeline.NewPipeline()
	stage := stages.UpvoteValidateStage(req)
	p.First = stage
	res.Error(p.Run())

	if res.Status {
		id := p.GetInt("ServiceId")[0]
		t := p.GetStringFirstOrDefault("UpvoteType")

		switch t {
		default:
			if e := toilet.DownvoteService(id); e != nil {
				res.Error(e)
			}
		}
	}

	sres.WriteJson(w, res)
}

func UnconfirmedInRange(w http.ResponseWriter, req *http.Request) {
	var res struct {
		sres.Response
		Services []toilet.Toilet
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
	s.HandleFunc("/s", GetUnconfirmeds).Methods("GET")
	s.HandleFunc("/range", UnconfirmedInRange).Methods("GET")
	s.HandleFunc("/upvote", UpvoteUnconfirmed).Methods("POST")
	s.HandleFunc("/downvote", DownvoteService).Methods("POST")
	return s
}
