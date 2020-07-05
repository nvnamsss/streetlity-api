package ratm

import (
	"bytes"
	"errors"
	"io"
	"log"
	"net/http"
	"strconv"
	"streelity/v1/model/atm"
	"streelity/v1/model/maintenance"
	"streelity/v1/sres"
	"streelity/v1/stages"

	"github.com/golang/geo/r2"
	"github.com/gorilla/mux"
	"github.com/nvnamsss/goinf/pipeline"
)

func GetService(w http.ResponseWriter, req *http.Request) {
	var res struct {
		sres.Response
		Service atm.Atm
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
			if service, e := atm.ServiceById(id); e != nil {
				res.Error(e)
			} else {
				res.Service = service
			}
			break
		case 2:
			lat := p.GetFloat("Lat")[0]
			lon := p.GetFloat("Lon")[0]
			if service, e := atm.ServiceByLocation(lat, lon); e != nil {
				res.Error(e)
			} else {
				res.Service = service
			}
			break
		case 3:
			address := p.GetString("Address")[0]
			if service, e := atm.ServiceByAddress(address); e != nil {
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
		Services []atm.Atm
	}
	res.Status = true

	p := pipeline.NewPipeline()
	stage := stages.QueryServicesValidateStage(req)
	p.First = stage
	res.Error(p.Run())

	if res.Status {
		address := p.GetString("Address")[0]
		if services, e := atm.ServicesByAddress(address); e != nil {
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
		Services []atm.Atm
	}
	res.Status = true

	if services, e := atm.AllServices(); e != nil {
		res.Error(e)
	} else {
		res.Services = services
	}

	sres.WriteJson(w, res)
}

func CreateService(w http.ResponseWriter, req *http.Request) {
	var res struct {
		sres.Response
		Service atm.Atm
	}
	res.Status = true
	p := pipeline.NewPipeline()
	stage := stages.CreateServiceValidate(req)
	bankStage := pipeline.NewStage(func() (str struct {
		BankId int64
	}, e error) {
		form := req.PostForm
		ids, ok := form["bank_id"]
		if !ok {
			return str, errors.New("bank_id param is missing")
		}

		if id, e := strconv.ParseInt(ids[0], 10, 64); e != nil {
			return str, errors.New("bank_id cannot parse to int64")
		} else {
			str.BankId = id
		}

		return
	})
	stage.NextStage(bankStage)
	p.First = stage

	res.Error(p.Run())

	if res.Status {
		lat := p.GetFloatFirstOrDefault("Lat")
		lon := p.GetFloatFirstOrDefault("Lon")
		address := p.GetStringFirstOrDefault("Address")
		note := p.GetStringFirstOrDefault("Note")
		images := p.GetString("Images")
		bank_id := p.GetIntFirstOrDefault("BankId")
		contributor := p.GetStringFirstOrDefault("Contributor")

		var ucf atm.Atm
		ucf.Lat = float32(lat)
		ucf.Lon = float32(lon)
		ucf.Address = address
		ucf.Note = note
		ucf.BankId = bank_id
		ucf.Contributor = contributor
		ucf.SetImages(images...)

		if service, e := atm.CreateService(ucf); e != nil {
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
		Services []atm.Atm
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

		res.Services = atm.ServicesInRange(location, max_range)
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
		maintenance.Import(buf.Bytes(), t)

	}

	sres.WriteJson(w, res)
}

func HandleService(router *mux.Router) *mux.Router {
	s := router.PathPrefix("/atm").Subrouter()

	s.HandleFunc("/", CreateService).Methods("POST")
	s.HandleFunc("/", GetService).Methods("GET")
	s.HandleFunc("/s", GetServices).Methods("GET")
	s.HandleFunc("/all", AllServices).Methods("GET")
	s.HandleFunc("/create", CreateService).Methods("POST")
	s.HandleFunc("/range", ServiceInRange).Methods("GET")
	s.HandleFunc("/import", Import).Methods("POST")

	return s
}
