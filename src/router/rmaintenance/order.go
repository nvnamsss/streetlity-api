package rmaintenance

import (
	"errors"
	"log"
	"net/http"
	"net/url"
	"streelity/v1/model/maintenance"
	"streelity/v1/sres"
	"streelity/v1/srpc"
	"streelity/v1/stages"

	"github.com/gorilla/mux"
	"github.com/nvnamsss/goinf/pipeline"
)

func EmergencyOrder(w http.ResponseWriter, req *http.Request) {
	var res sres.Response = sres.Response{Status: true}

	p := pipeline.NewPipeline()
	stage := stages.EmergencyOrderValidate(req)
	p.First = stage
	res.Error(p.Run())

	if res.Status {
		common_user := p.GetString("CommonUser")[0]
		eusers := p.GetString("EmergencyMaintenance")
		reason := p.GetString("Reason")[0]
		note := p.GetStringFirstOrDefault("Note")
		phone := p.GetString("Phone")[0]
		order_type := "2"
		if order, e := srpc.RequestOrder(url.Values{
			"maintenance_users": eusers,
			"common_user":       {common_user},
			"reason":            {reason},
			"phone":             {phone},
			"note":              {note},
			"type":              {order_type},
		}); e != nil {
			res.Error(e)
		} else {
			log.Println(order)
			res.Status = order.Status
			res.Message = order.Message
		}
	}

	sres.WriteJson(w, res)
}

func CommonOrder(w http.ResponseWriter, req *http.Request) {
	var res sres.Response = sres.Response{Status: true, Message: "Order successfully"}

	p := pipeline.NewPipeline()
	stage := stages.CommonOrderValidate(req)
	p.First = stage
	res.Error(p.Run())

	if res.Status {
		common_user := p.GetString("CommonUser")[0]
		service_ids := p.GetInt("ServiceId")
		reason := p.GetString("Reason")[0]
		note := p.GetStringFirstOrDefault("Note")
		phone := p.GetString("Phone")[0]
		t := "1"
		services := maintenance.ServicesByIds(service_ids...)
		maintenance_users := []string{}
		for _, s := range services {
			if s.Maintainer != "" {
				maintainers := s.GetMaintainers()
				for maintainer, _ := range maintainers {
					maintenance_users = append(maintenance_users, maintainer)
				}
			}
		}

		if len(maintenance_users) == 0 {
			res.Error(errors.New("cannot find any suitable maintenance user"))
		} else {
			if order, e := srpc.RequestOrder(url.Values{
				"maintenance_users": maintenance_users,
				"common_user":       {common_user},
				"reason":            {reason},
				"phone":             {phone},
				"note":              {note},
				"type":              {t},
			}); e != nil {
				res.Error(e)
			} else {
				log.Println(order)
				res.Status = order.Status
				res.Message = order.Message
			}
		}
	}

	sres.WriteJson(w, res)
}

func HandleOrder(router *mux.Router) *mux.Router {
	s := router.PathPrefix("/order").Subrouter()
	s.HandleFunc("/e", EmergencyOrder).Methods("POST")
	return s
}
