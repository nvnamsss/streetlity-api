package router

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"streelity/v1/middleware"
	"streelity/v1/model"
	"streelity/v1/model/toilet"
	"streelity/v1/router/rtoilet"
	"streelity/v1/sres"
	"streelity/v1/stages"

	"github.com/golang/geo/r2"
	"github.com/gorilla/mux"
	"github.com/nvnamsss/goinf/pipeline"
)

/*AUTH REQUIRED*/

func addToilet(w http.ResponseWriter, req *http.Request) {
	var res sres.Response = sres.Response{Status: true}
	req.ParseForm()

	var pipe *pipeline.Pipeline = pipeline.NewPipeline()
	validateParamsStage := stages.AddingServiceValidateStage(req)
	pipe.First = validateParamsStage
	res.Error(pipe.Run())

	if res.Status {
		var s toilet.ToiletUcf
		lat := pipe.GetFloat("Lat")[0]
		lon := pipe.GetFloat("Lon")[0]
		note := pipe.GetString("Note")[0]
		address := pipe.GetString("Address")[0]
		images := pipe.GetString("Images")
		s.Lat = float32(lat)
		s.Lon = float32(lon)
		s.Note = note
		s.SetImages(images...)
		s.Address = address
		err := toilet.CreateUcf(s)

		if err != nil {
			res.Status = false
			res.Message = err.Error()
		} else {
			res.Message = "Create new fuel is succeed"
		}
	}

	sres.WriteJson(w, res)
}

func upvoteToilet(w http.ResponseWriter, req *http.Request) {
	var res sres.Response = sres.Response{Status: true}

	req.ParseForm()
	p := pipeline.NewPipeline()
	vStage := pipeline.NewStage(func() (str struct{ Id int64 }, e error) {
		form := req.PostForm
		_, ok := form["id"]
		if !ok {
			e = errors.New("id params is missing")
			return
		}

		str.Id, e = strconv.ParseInt(form["id"][0], 10, 64)
		return
	})
	p.First = vStage
	res.Error(p.Run())

	if res.Status {
		var id int64 = p.GetInt("Id")[0]
		res.Error(toilet.UpvoteUcf(id))
	}

	sres.WriteJson(w, res)
}

/*NON-AUTH REQUIRED*/

func getAllToilets(w http.ResponseWriter, req *http.Request) {
	var res struct {
		sres.Response
		Toilets []toilet.Toilet
	}
	res.Status = true

	if res.Status {
		res.Toilets = toilet.AllToilets()
	}

	sres.WriteJson(w, res)
}

func updateToilet(w http.ResponseWriter, req *http.Request) {

}

func getToiletInRange(w http.ResponseWriter, req *http.Request) {
	var res struct {
		sres.Response
		Toilets []toilet.Toilet
	}

	res.Status = true
	query := req.URL.Query()
	var pipe *pipeline.Pipeline = pipeline.NewPipeline()

	validateParamsStage := pipeline.NewStage(func() error {
		location, locationOk := query["location"]
		if !locationOk {
			return errors.New("location param is missing")
		} else {
			if len(location) < 2 {
				return errors.New("location param must have 2 values")
			}
		}

		_, rangeOk := query["range"]
		if !rangeOk {
			return errors.New("range param is missing")
		}

		return nil
	})

	parseValueStage := pipeline.NewStage(func() (str struct {
		Lat   float64
		Lon   float64
		Range float64
	}, e error) {
		var (
			latErr   error
			lonErr   error
			rangeErr error
		)

		str.Lat, latErr = strconv.ParseFloat(query["location"][0], 64)
		str.Lon, lonErr = strconv.ParseFloat(query["location"][1], 64)
		str.Range, rangeErr = strconv.ParseFloat(query["range"][0], 64)
		if latErr != nil {
			return str, errors.New("cannot parse location[0] to float")
		}
		if lonErr != nil {
			return str, errors.New("cannot parse location[1] to float")
		}
		if rangeErr != nil {
			return str, errors.New("cannot parse range to float")
		}
		return str, nil
	})

	validateParamsStage.NextStage(parseValueStage)
	pipe.First = validateParamsStage

	res.Error(pipe.Run())

	if res.Status {
		lat := pipe.GetFloat("Lat")[0]
		lon := pipe.GetFloat("Lon")[0]
		max_range := pipe.GetFloat("Range")[0]
		var location r2.Point = r2.Point{X: lat, Y: lon}

		res.Toilets = toilet.ToiletsInRange(location, max_range)
	}

	sres.WriteJson(w, res)
}

func getToilet(w http.ResponseWriter, req *http.Request) {
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
		s := model.Service{}

		switch c {
		case 0:
			break
		case 1:
			s.Id = p.GetInt("Id")[0]
			break
		case 2:
			s.Lat = float32(p.GetFloat("Lat")[0])
			s.Lon = float32(p.GetFloat("Lat")[0])
			break
		case 3:
			s.Address = p.GetString("Address")[0]
			break
		}

		if m, e := toilet.ToiletByService(s); e == nil {
			res.Service = m
		} else {
			res.Error(e)
		}

	}

	sres.WriteJson(w, res)
}

func HandleToilet(router *mux.Router) {
	log.Println("[Router]", "Handling Toilet")
	s := router.PathPrefix("/toilet").Subrouter()
	s.HandleFunc("/all", getAllToilets).Methods("GET")
	s.HandleFunc("/update", updateToilet).Methods("POST")
	s.HandleFunc("/range", getToiletInRange).Methods("GET")
	s.HandleFunc("/", getToilet).Methods("GET")

	r := s.PathPrefix("/add").Subrouter()
	r.HandleFunc("", addToilet).Methods("POST")
	r.Use(middleware.Authenticate)

	rtoilet.HandleUnconfirmed(router)
	// r = s.PathPrefix("/upvote").Subrouter()
	// r.HandleFunc("", updateToilet).Methods("POST")
	// r.Use(middleware.Authenticate)
}
