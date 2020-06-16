package rfuel

import (
	"net/http"
	"streelity/v1/model"
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
		if service, e := fuel.FuelById(id); e != nil {
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
		ucf := fuel.FuelUcf{model.ServiceUcf{Lat: float32(lat), Lon: float32(lon), Address: address, Note: note, Images: ""}}
		for _, image := range images {
			ucf.Images += image + ";"
		}

		if e := fuel.AddFuelUcf(ucf); e != nil {
			res.Error(e)
		}
	}

	sres.WriteJson(w, res)
}

func HandleFuel(router *mux.Router) {
	s := router.PathPrefix("/fuel").Subrouter()

	s.HandleFunc("/", CreateService).Methods("POST")
	s.HandleFunc("/", GetService).Methods("GET")
}
