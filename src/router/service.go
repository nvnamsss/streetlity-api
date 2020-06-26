package router

import (
	"log"
	"net/http"
	"streelity/v1/middleware"
	"streelity/v1/model/atm"
	"streelity/v1/model/fuel"
	"streelity/v1/model/maintenance"
	"streelity/v1/model/toilet"
	"streelity/v1/sres"
	"streelity/v1/stages"

	"github.com/golang/geo/r2"
	"github.com/gorilla/mux"
	"github.com/nvnamsss/goinf/pipeline"
)

func ServiceInRange(w http.ResponseWriter, req *http.Request) {
	var res struct {
		sres.Response
		Fuels        []fuel.Fuel
		Atms         []atm.Atm
		Maintenances []maintenance.Maintenance
		Toilets      []toilet.Toilet
	}
	res.Status = true

	p := pipeline.NewPipeline()
	stage := stages.InRangeServiceValidateStage(req)
	p.First = stage
	res.Error(p.Run())

	if res.Status {
		lat := p.GetFloatFirstOrDefault("Lat")
		lon := p.GetFloatFirstOrDefault("Lon")
		max_range := p.GetFloatFirstOrDefault("Range")
		var location r2.Point = r2.Point{X: lat, Y: lon}

		res.Fuels = fuel.ServicesInRange(location, max_range)
		res.Atms = atm.ServicesInRange(location, max_range)
		res.Maintenances = maintenance.ServicesInRange(location, max_range)
		res.Toilets = toilet.ServicesInRange(location, max_range)
	}

	sres.WriteJson(w, res)
}

func HandleService(router *mux.Router) {
	log.Println("[Router]", "Handling service")

	s := router.PathPrefix("/service").Subrouter()
	s.HandleFunc("/range", ServiceInRange).Methods("GET")
	HandleFuel(s)
	HandleAtm(s)
	HandleToilet(s)
	HandleMaintenance(s)

	middleware.Versioning(s, "1.0.0", "2.1.0")
}
