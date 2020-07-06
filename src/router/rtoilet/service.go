package rtoilet

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"streelity/v1/model/toilet"
	"streelity/v1/sres"
	"streelity/v1/stages"

	"github.com/golang/geo/r2"
	"github.com/gorilla/mux"
	"github.com/nvnamsss/goinf/pipeline"
)

func GetService(w http.ResponseWriter, req *http.Request) {
	var res struct {
		sres.Response
		Service toilet.Toilet
	}
	res.Status = true
	p := pipeline.NewPipeline()
	stage := stages.QueryServiceValidateStage(req)
	p.First = stage
	res.Error(p.Run())
	if res.Status {
		c := p.GetInt("Case")[0]
		switch c {
		case 1:
			id := p.GetInt("Id")[0]
			if service, e := toilet.ServiceById(id); e != nil {
				res.Error(e)
			} else {
				res.Service = service
			}
			break
		case 2:
			lat := p.GetFloat("Lat")[0]
			lon := p.GetFloat("Lon")[0]
			if service, e := toilet.ServiceByLocation(lat, lon); e != nil {
				res.Error(e)
			} else {
				res.Service = service
			}
			break
		case 3:
			address := p.GetString("Address")[0]
			if service, e := toilet.ServiceByAddress(address); e != nil {
				res.Error(e)
			} else {
				res.Service = service
			}
			break
		}
	}
	sres.WriteJson(w, res)
}

func GetServices(w http.ResponseWriter, req *http.Request) {
	var res struct {
		sres.Response
		Services []toilet.Toilet
	}
	res.Status = true

	p := pipeline.NewPipeline()
	stage := stages.QueryServicesValidateStage(req)
	p.First = stage
	res.Error(p.Run())

	if res.Status {
		address := p.GetString("Address")[0]
		if services, e := toilet.ServicesByAddress(address); e != nil {
			res.Error(e)
		} else {
			res.Services = services
		}
	}

	sres.WriteJson(w, res)
}

func AllServices(w http.ResponseWriter, req *http.Request) {
	var res struct {
		sres.Response
		Services []toilet.Toilet
	}
	res.Status = true

	if services, e := toilet.AllServices(); e != nil {
		res.Error(e)
	} else {
		res.Services = services
	}

	sres.WriteJson(w, res)
}

func CreateService(w http.ResponseWriter, req *http.Request) {
	var res struct {
		sres.Response
		Service toilet.ToiletUcf
	}
	res.Status = true
	p := pipeline.NewPipeline()
	stage := stages.CreateServiceValidate(req)
	p.First = stage

	res.Error(p.Run())

	if res.Status {
		lat := p.GetFloatFirstOrDefault("Lat")
		lon := p.GetFloatFirstOrDefault("Lon")
		address := p.GetStringFirstOrDefault("Address")
		note := p.GetStringFirstOrDefault("Note")
		images := p.GetString("Images")
		contributor := p.GetString("Contributor")[0]

		var ucf toilet.ToiletUcf
		ucf.Lat = float32(lat)
		ucf.Lon = float32(lon)
		ucf.Address = address
		ucf.Note = note
		ucf.Contributor = contributor
		ucf.SetImages(images...)

		if service, e := toilet.CreateUcf(ucf); e != nil {
			res.Error(e)
		} else {
			res.Service = service
		}
	}

	sres.WriteJson(w, res)
}

func ServiceInRange(w http.ResponseWriter, req *http.Request) {
	var res struct {
		sres.Response
		Services []toilet.Toilet
	}
	res.Status = true
	pipe := pipeline.NewPipeline()
	stage := stages.InRangeServiceValidateStage(req)
	pipe.First = stage

	res.Error(pipe.Run())

	if res.Status {
		lat := pipe.GetFloatFirstOrDefault("Lat")
		lon := pipe.GetFloatFirstOrDefault("Lon")
		max_range := pipe.GetFloatFirstOrDefault("Range")
		var location r2.Point = r2.Point{X: lat, Y: lon}

		res.Services = toilet.ServicesInRange(location, max_range)
	}

	sres.WriteJson(w, res)
}

func Import(w http.ResponseWriter, req *http.Request) {
	var res sres.Response = sres.Response{Status: true}

	req.ParseForm()
	p := pipeline.NewPipeline()
	stage := stages.ImportValidate(req.URL.Query())
	p.First = stage
	res.Error(p.Run())

	if res.Status {
		t := p.GetString("Type")[0]
		req.ParseMultipartForm(32 << 20) // limit your max input length!
		file, _, e := req.FormFile("f")

		if e != nil {
			log.Println("[Upload]", "cannot find", "f param", "in the form")
		}

		defer file.Close()

		var buf bytes.Buffer
		io.Copy(&buf, file)
		toilet.Import(buf.Bytes(), t)

	}

	sres.WriteJson(w, res)
}

func HandleService(router *mux.Router) *mux.Router {
	s := router.PathPrefix("/toilet").Subrouter()

	s.HandleFunc("/", CreateService).Methods("POST")
	s.HandleFunc("/", GetService).Methods("GET")
	s.HandleFunc("/s", GetServices).Methods("GET")
	s.HandleFunc("/all", AllServices).Methods("GET")
	s.HandleFunc("/range", ServiceInRange).Methods("GET")
	s.HandleFunc("/create", CreateService).Methods("POST")
	s.HandleFunc("/import", Import).Methods("POST")

	return s
}
