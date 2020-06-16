package rfuel

import (
	"net/http"
	"streelity/v1/model/fuel"
	"streelity/v1/sres"
	"streelity/v1/stages"

	"github.com/gorilla/mux"
	"github.com/nvnamsss/goinf/pipeline"
)

func GetService(w http.ResponseWriter, req *http.Request) {
	var res struct {
		sres.Response
		Service fuel.Fuel
	}

	p := pipeline.NewPipeline()
	stage := stages.IdValidateStage(req.URL.Query())
	p.First = stage

	res.Error(p.Run())

	if res.Status {
		id := p.GetIntFirstOrDefault("Id")
		if service, e := fuel.ServiceById(id); e != nil {
			res.Error(e)
		} else {
			res.Service = service
		}
	}

	sres.WriteJson(w, res)
}

func CreateService(w http.ResponseWriter, req *http.Request) {
	var res struct {
		sres.Response
		Service fuel.FuelUcf
	}

	p := pipeline.NewPipeline()
	stage := stages.AddingServiceValidateStage(req)
	p.First = stage

	res.Error(p.Run())

	if res.Status {
		lat := p.GetFloatFirstOrDefault("Lat")
		lon := p.GetFloatFirstOrDefault("Lon")
		address := p.GetStringFirstOrDefault("Address")
		note := p.GetStringFirstOrDefault("Note")
		images := p.GetString("Images")
		var ucf fuel.FuelUcf
		ucf.Lat = float32(lat)
		ucf.Lon = float32(lon)
		ucf.Address = address
		ucf.Note = note
		ucf.SetImages(images...)

		if e := fuel.CreateUcf(ucf); e != nil {
			res.Error(e)
		}
	}

	sres.WriteJson(w, res)
}

func HandleService(router *mux.Router) {
	s := router.PathPrefix("/fuel").Subrouter()

	s.HandleFunc("/", CreateService).Methods("POST")
	s.HandleFunc("/", GetService).Methods("GET")
}
