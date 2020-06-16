package router

import (
	"errors"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"streelity/v1/middleware"
	"streelity/v1/model"
	"streelity/v1/model/maintenance"
	"streelity/v1/sres"
	"streelity/v1/srpc"
	"streelity/v1/stages"

	"github.com/golang/geo/r2"
	"github.com/gorilla/mux"
	"github.com/nvnamsss/goinf/pipeline"
)

/*AUTH REQUIRED*/
func orderMaintenance(w http.ResponseWriter, req *http.Request) {
	var res sres.Response = sres.Response{Status: true}

	req.ParseForm()
	p := pipeline.NewPipeline()

	vStage := pipeline.NewStage(func() (str struct {
		CommonUser string
		Reason     string
		Note       string
		ServiceId  []int64
	}, e error) {
		form := req.PostForm

		commonUsers, ok := form["commonUser"]
		if !ok {
			return str, errors.New("commonUser param is misisng")
		}

		ids, ok := form["id"]

		if !ok {
			return str, errors.New("id param is missing")
		}

		reasons, ok := form["reason"]
		if !ok {
			return str, errors.New("reason param is missing")
		}

		notes, ok := form["note"]
		if ok {
			str.Note = notes[0]
		}

		for _, id := range ids {
			v, e := strconv.ParseInt(id, 10, 64)
			if e != nil {
				continue
			}

			str.ServiceId = append(str.ServiceId, v)
		}

		str.CommonUser = commonUsers[0]
		str.Reason = reasons[0]

		return
	})
	p.First = vStage
	res.Error(p.Run())

	if res.Status {
		commonUser := p.GetString("CommonUser")[0]
		service_ids := p.GetInt("ServiceId")
		reason := p.GetString("Reason")[0]
		note := p.GetStringFirstOrDefault("Note")

		services := maintenance.MaintenanceByIds(service_ids...)
		maintenanceIds := []string{}
		for _, s := range services {
			maintenanceIds = append(maintenanceIds, s.Owner)
		}

		if len(maintenanceIds) == 0 {
			res.Error(errors.New("cannot find any suitable maintenance user"))
		} else {
			order, e := srpc.RequestOrder(url.Values{
				"maintenanceUser": maintenanceIds,
				"commonUser":      {commonUser},
				"reason":          {reason},
				"note":            {note},
			})
			log.Println(order)
			res.Error(e)
		}

	}

	sres.WriteJson(w, res)
}

func acceptOrderMaintenance(w http.ResponseWriter, req *http.Request) {
	var res sres.Response = sres.Response{Status: true}

	req.ParseForm()
	p := pipeline.NewPipeline()
	vStage := pipeline.NewStage(func() (str struct {
		User      string
		OrderId   int64
		Timestamp int64
	}, e error) {
		form := req.PostForm
		_, ok := form["user"]
		if !ok {
			return str, errors.New("user param is missing")
		}

		timestamps, ok := form["timestamp"]
		if !ok {
			return str, errors.New("timestamp param is missing")
		}

		ids, ok := form["orderId"]
		if !ok {
			return str, errors.New("oderId param is missing")
		}

		id, e := strconv.ParseInt(ids[0], 10, 64)

		if e != nil {
			return str, errors.New("cannot parse orderId to int")
		}

		timestamp, e := strconv.ParseInt(timestamps[0], 10, 64)

		if e != nil {
			return str, errors.New("cannot parse timestamp to int")
		}

		str.User = form["user"][0]
		str.OrderId = id
		str.Timestamp = timestamp

		return
	})
	p.First = vStage

	res.Error(p.Run())

	if res.Status {
		user := p.GetString("User")[0]
		id := p.GetInt("OrderId")[0]
		timestamp := p.GetInt("Timestamp")[0]

		maintenance.UpdateMaintenanceHistory(id, user, timestamp)
	}

	sres.WriteJson(w, res)
}

func updateMaintenance(w http.ResponseWriter, req *http.Request) {
	var res sres.Response = sres.Response{Status: true}

	req.ParseForm()
	p := pipeline.NewPipeline()
	vStage := pipeline.NewStage(func() (str struct{ Id int64 }, e error) {
		form := req.PostForm
		_, ok := form["id"]
		if !ok {
			return str, errors.New("id param is missing")
		}

		id, err := strconv.ParseInt(form["id"][0], 10, 64)
		if err != nil {
			return str, errors.New("cannot parse id to int")
		}
		str.Id = id

		return
	})
	p.First = vStage

	res.Error(p.Run())

	if res.Status {
		id := p.GetInt("Id")[0]
		values := make(map[string]string)
		form := req.PostForm
		for key, value := range form {
			values[key] = value[0]
		}

		maintenance.UpdateMaintenance(id, values)
	}

	sres.WriteJson(w, res)
}

func addMaintenance(w http.ResponseWriter, req *http.Request) {
	var res struct {
		sres.Response
		Service maintenance.Maintenance
	}
	res.Status = true

	req.ParseForm()

	var pipe *pipeline.Pipeline = pipeline.NewPipeline()
	validateParamsStage := stages.AddingServiceValidateStage(req)
	pipe.First = validateParamsStage
	res.Error(pipe.Run())

	if res.Status {
		var s maintenance.MaintenanceUcf
		lat := pipe.GetFloat("Lat")[0]
		lon := pipe.GetFloat("Lon")[0]
		note := pipe.GetString("Note")[0]
		address := pipe.GetString("Address")[0]
		images := pipe.GetString("Images")
		name, ok := req.PostForm["name"]
		s.Lat = float32(lat)
		s.Lon = float32(lon)
		s.Note = note
		s.SetImages(images...)
		s.Address = address
		if ok {
			s.Name = name[0]
		}
		err := maintenance.AddMaintenanceUcf(s)

		if err != nil {
			res.Status = false
			res.Message = err.Error()
		} else {
			res.Message = "Create new Maintenance successfully"
		}
	}
	log.Println(res)
	sres.WriteJson(w, res)
}

/*NON-AUTH REQUIRED*/

func getMaintenances(w http.ResponseWriter, req *http.Request) {
	var res struct {
		sres.Response
		Maintenance []maintenance.Maintenance
	}

	res.Status = true

	res.Maintenance = maintenance.AllMaintenances()

	sres.WriteJson(w, res)
}

func getMaintenanceById(w http.ResponseWriter, req *http.Request) {
	var res struct {
		sres.Response
		Maintenance maintenance.Maintenance
	}

	res.Status = true
	query := req.URL.Query()

	pipe := pipeline.NewPipeline()
	validateParams := pipeline.NewStage(func() (str struct{ Id int64 }, e error) {
		ids, ok := query["id"]
		if !ok {
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
		s, e := maintenance.MaintenanceById(id)
		res.Error(e)
		if res.Status {
			res.Maintenance = s
		}
	}

	sres.WriteJson(w, res)
}

//getFuelInRange process the in-range query. the request must provide there
//
// Parameters:
// 	- `location`: X and Y coordinator
// 	- `range` : range to find
//
func getMaintenanceInRange(w http.ResponseWriter, req *http.Request) {
	var res struct {
		sres.Response
		Maintenances []maintenance.Maintenance
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

		res.Maintenances = maintenance.MaintenancesInRange(location, max_range)
	}

	sres.WriteJson(w, res)
}

func getMaintenance(w http.ResponseWriter, req *http.Request) {
	var res struct {
		sres.Response
		Service maintenance.Maintenance
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

		if m, e := maintenance.MaintenanceByService(s); e == nil {
			res.Service = m
		} else {
			res.Error(e)
		}

	}

	sres.WriteJson(w, res)
}

func HandleMaintenance(router *mux.Router) {
	log.Println("[Router]", "Handling fuel")
	s := router.PathPrefix("/maintenance").Subrouter()
	s.HandleFunc("/all", getMaintenances).Methods("GET")
	s.HandleFunc("/update", updateMaintenance).Methods("POST")
	s.HandleFunc("/range", getMaintenanceInRange).Methods("GET")
	s.HandleFunc("/id", getMaintenanceById).Methods("GET")
	s.HandleFunc("/add", addMaintenance).Methods("POST")
	s.HandleFunc("/order", orderMaintenance).Methods("POST")
	s.HandleFunc("/accept", acceptOrderMaintenance).Methods("POST")
	s.HandleFunc("/", getMaintenance).Methods("GET")

	r := s.PathPrefix("/add").Subrouter()
	r.HandleFunc("", addMaintenance).Methods("POST")
	r.Use(middleware.Authenticate)

	r = s.PathPrefix("/upvote").Subrouter()
	r.HandleFunc("", upvoteMaintenance).Methods("POST")
	r.Use(middleware.Authenticate)

}
