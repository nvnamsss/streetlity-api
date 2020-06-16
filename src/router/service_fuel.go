package router

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"streelity/v1/middleware"
	"streelity/v1/model"
	"streelity/v1/model/fuel"
	"streelity/v1/router/rfuel"
	"streelity/v1/sres"
	"streelity/v1/stages"

	"github.com/golang/geo/r2"
	"github.com/gorilla/mux"
	"github.com/nvnamsss/goinf/pipeline"
)

/*AUTH REQUIRED*/
func updateFuel(w http.ResponseWriter, req *http.Request) {
	var res sres.Response
	res.Status = true

	req.ParseForm()
	form := req.PostForm

	var pipe *pipeline.Pipeline = pipeline.NewPipeline()
	validateParamsStage := pipeline.NewStage(func() error {
		_, idOk := form["id"]
		location, locationOk := form["location"]

		if !idOk {
			return errors.New("id param is missing")
		}

		if !locationOk {
			return errors.New("location param is missing")
		} else {
			if len(location) < 2 {
				return errors.New("location param must have 2 values")
			}
		}

		return nil
	})

	parseValueStage := pipeline.NewStage(func() error {
		_, idErr := strconv.ParseInt(form["id"][0], 10, 64)
		_, latErr := strconv.ParseFloat(form["location"][0], 64)
		_, lonErr := strconv.ParseFloat(form["location"][1], 64)

		if idErr != nil {
			return errors.New("cannot parse id to int")
		}

		if latErr != nil {
			return errors.New("cannot parse location[0] to float")
		}

		if lonErr != nil {
			return errors.New("cannot parse location[1] to float")
		}

		return nil
	})

	validateParamsStage.NextStage(parseValueStage)
	pipe.First = validateParamsStage
	res.Error(pipe.Run())

	if res.Status {
		var f fuel.Fuel
		id, _ := strconv.ParseInt(form["id"][0], 10, 64)
		if err := model.Db.Where(&fuel.Fuel{Service: model.Service{Id: id}}).First(&f).Error; err != nil {
			res.Status = false
			res.Message = err.Error()
		}

	}

	sres.WriteJson(w, res)
}

func addFuel(w http.ResponseWriter, req *http.Request) {
	var res struct {
		sres.Response
		Service fuel.FuelUcf
	}
	res.Status = true

	req.ParseForm()

	var pipe *pipeline.Pipeline = pipeline.NewPipeline()
	validateParamsStage := stages.AddingServiceValidateStage(req)
	pipe.First = validateParamsStage
	res.Error(pipe.Run())

	if res.Status {
		var s fuel.FuelUcf
		lat := pipe.GetFloat("Lat")[0]
		lon := pipe.GetFloat("Lon")[0]
		note := pipe.GetString("Note")[0]
		address := pipe.GetString("Address")[0]
		images := pipe.GetString("Images")
		s.Lat = float32(lat)
		s.Lon = float32(lon)
		s.Note = note
		s.Address = address
		s.SetImages(images...)
		err := fuel.AddFuelUcf(s)

		if err != nil {
			res.Status = false
			res.Message = err.Error()
		} else {
			res.Message = "Create new fuel is succeed"
		}
	}

	sres.WriteJson(w, res)
}

func getFuelReview(w http.ResponseWriter, req *http.Request) {
	var res struct {
		sres.Response
		Review []fuel.Review
	}
	res.Status = true
	p := pipeline.NewPipeline()
	stage := pipeline.NewStage(func() (str struct {
		ReviewId int64
		Order    int64
	}, e error) {
		query := req.URL.Query()

		_, ok := query["review_id"]
		if !ok {
			return str, errors.New("review_id param is missing")
		}

		_, ok = query["order"]
		if !ok {
			return str, errors.New("order param is missing")
		}

		review_id, e := strconv.ParseInt(query["review_id"][0], 10, 64)
		if e != nil {
			return str, errors.New("review_id cannot parse to int64")
		}

		order, e := strconv.ParseInt(query["order"][0], 10, 64)
		if e != nil {
			return str, errors.New("order cannot parse to int64")
		}

		str.ReviewId = review_id
		str.Order = order
		return
	})

	p.First = stage
	res.Error(p.Run())

	if res.Status {
		review_id := p.GetIntFirstOrDefault("ReviewId")
		order := p.GetIntFirstOrDefault("Order")

		fuel.ReviewByService(review_id, order, 5)
	}
}

/*NON-AUTH REQUIRED*/

func getFuels(w http.ResponseWriter, req *http.Request) {
	var res struct {
		sres.Response
		Fuels []fuel.Fuel
	}

	res.Status = true

	res.Fuels = fuel.AllFuels()

	log.Println("[GetFuels]", res.Fuels)

	sres.WriteJson(w, res)
}

//getFuelInRange process the in-range query. the request must provide there
//
// Parameters:
// 	- `location`: X and Y coordinator
// 	- `range` : range to find
//
func getFuelInRange(w http.ResponseWriter, req *http.Request) {
	var res struct {
		sres.Response
		Fuels []fuel.Fuel
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

		res.Fuels = fuel.FuelsInRange(location, max_range)
	}

	sres.WriteJson(w, res)
}

func upvoteFuel(w http.ResponseWriter, req *http.Request) {
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
		res.Error(fuel.UpvoteFuelUcf(id))
	}

	sres.WriteJson(w, res)
}

func getFuel(w http.ResponseWriter, req *http.Request) {
	var res struct {
		sres.Response
		Service fuel.Fuel
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

		if m, e := fuel.FuelByService(s); e == nil {
			res.Service = m
		} else {
			res.Error(e)
		}

	}

	sres.WriteJson(w, res)
}

func HandleFuel(router *mux.Router) {
	log.Println("[Router]", "Handling fuel")
	s := router.PathPrefix("/fuel").Subrouter()
	s.HandleFunc("/all", getFuels).Methods("GET")
	s.HandleFunc("/update", updateFuel).Methods("POST")
	s.HandleFunc("/range", getFuelInRange).Methods("GET")
	s.HandleFunc("/", getFuel).Methods("GET")

	rfuel.HandleReview(s)
	// s.HandleFunc("/review", addFuelReview).Methods("POST")

	r := s.PathPrefix("/add").Subrouter()
	r.HandleFunc("", addFuel).Methods("POST")
	r.Use(middleware.Authenticate)

	r = s.PathPrefix("/upvote").Subrouter()
	r.HandleFunc("", upvoteFuel).Methods("POST")
	r.Use(middleware.Authenticate)
}
