package router

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"streelity/v1/model"

	"github.com/golang/geo/r2"
	"github.com/gorilla/mux"
	"github.com/nvnamsss/goinf/pipeline"
)

/*AUTH REQUIRED*/
func orderMaintenance(w http.ResponseWriter, req *http.Request) {
	var res Response = Response{Status: true}

	WriteJson(w, res)
}

func acceptOrderMaintenance(w http.ResponseWriter, req *http.Request) {
	var res Response = Response{Status: true}

	WriteJson(w, res)
}

func updateMaintenance(w http.ResponseWriter, req *http.Request) {
	var res Response = Response{Status: true}

	req.ParseForm()

	var pipe *pipeline.Pipeline = pipeline.NewPipeline()
	validateParamsStage := pipeline.NewStage(func() error {
		form := req.PostForm
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
		form := req.PostForm
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
		var m model.Maintenance
		id := pipe.GetInt("Id")[0]
		if err := model.Db.Where(&model.Maintenance{Service: model.Service{Id: id}}).First(&m).Error; err != nil {
			res.Status = false
			res.Message = err.Error()
		}

	}

	WriteJson(w, res)
}

func addMaintenance(w http.ResponseWriter, req *http.Request) {
	var res Response = Response{Status: true}

	req.ParseForm()

	var pipe *pipeline.Pipeline = pipeline.NewPipeline()
	validateParamsStage := AddingServiceValidateStage(req)
	parseValueStage := AddingServiceParsingStage(req)

	validateParamsStage.NextStage(parseValueStage)
	pipe.First = validateParamsStage
	res.Error(pipe.Run())

	if res.Status {
		var m model.MaintenanceUcf
		lat := pipe.GetFloat("Lat")[0]
		lon := pipe.GetFloat("Lon")[0]
		m.Lat = float32(lat)
		m.Lon = float32(lon)

		err := model.AddMaintenanceUcf(m)

		if err != nil {
			res.Status = false
			res.Message = err.Error()
		} else {
			res.Message = "Create new Maintenance successfully"
		}
	}

	WriteJson(w, res)
}

func upvoteMaintenance(w http.ResponseWriter, req *http.Request) {
	var res Response = Response{Status: true}

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
		res.Error(model.UpvoteMaintenanceUcf(id))
	}

	WriteJson(w, res)
}

/*NON-AUTH REQUIRED*/

func getMaintenances(w http.ResponseWriter, req *http.Request) {
	var res struct {
		Response
		Maintenance []model.Maintenance
	}

	res.Status = true

	res.Maintenance = model.AllMaintenances()

	WriteJson(w, res)
}

func getMaintenance(w http.ResponseWriter, req *http.Request) {
	var res struct {
		Response
		Maintenance model.Maintenance
	}

	res.Status = true
	query := req.URL.Query()

	pipe := pipeline.NewPipeline()
	validateParams := pipeline.NewStage(func() (str struct{ Id int64 }, e error) {
		ids, ok := query["id"]
		if ok {
			return str, errors.New("id param is missing")
		}
		var err error
		str.Id, err = strconv.ParseInt(ids[0], 10, 64)

		if err != nil {
			return str, errors.New("cannot parse id to int")
		}

		return str, nil
	})

	pipe.First = validateParams

	res.Error(pipe.Run())

	if res.Status {
		id := pipe.GetInt("Id")[0]
		res.Maintenance = model.MaintenanceById(id)
	}

	WriteJson(w, res)
}

//getFuelInRange process the in-range query. the request must provide there
//
// Parameters:
// 	- `location`: X and Y coordinator
// 	- `range` : range to find
//
func getMaintenanceInRange(w http.ResponseWriter, req *http.Request) {
	var res struct {
		Response
		Maintenances []model.Maintenance
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

		res.Maintenances = model.MaintenancesInRange(location, max_range)
	}

	WriteJson(w, res)
}

func HandleMaintenance(router *mux.Router) {
	log.Println("[Router]", "Handling fuel")
	s := router.PathPrefix("/maintain").Subrouter()
	s.HandleFunc("/all", getMaintenances).Methods("GET")
	s.HandleFunc("/update", updateMaintenance).Methods("POST")
	s.HandleFunc("/id", getMaintenance).Methods("GET")
	s.HandleFunc("/range", getMaintenanceInRange).Methods("GET")
	s.HandleFunc("/add", addMaintenance).Methods("POST")

	r := s.PathPrefix("/add").Subrouter()
	r.HandleFunc("", addMaintenance).Methods("POST")
	r.Use(Authenticate)

	r = s.PathPrefix("/upvote").Subrouter()
	r.HandleFunc("", upvoteMaintenance).Methods("POST")
	r.Use(Authenticate)

}
