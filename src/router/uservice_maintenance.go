package router

import (
	"errors"
	"net/http"
	"strconv"
	"streelity/v1/model/maintenance"
	"streelity/v1/router/sres"

	"github.com/nvnamsss/goinf/pipeline"
)

func upvoteMaintenance(w http.ResponseWriter, req *http.Request) {
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
		t, ok := req.PostForm["type"]

		if ok && t[0] == "immediately" {
			res.Error(maintenance.UpvoteMaintenanceUcfByIdImmediately(id))
		} else {
			res.Error(maintenance.UpvoteMaintenanceUcfById(id))
		}
	}

	sres.WriteJson(w, res)
}

func getUMaintenance(w http.ResponseWriter, req *http.Request) {
	var res struct {
		sres.Response
		Services []maintenance.MaintenanceUcf
	}

	p := pipeline.NewPipeline()
	stage := pipeline.NewStage(func() (str struct {
		Lat     float64
		Lon     float64
		Range   float64
		Address string
	}, e error) {
		query := req.URL.Query()

		location, ok := query["location"]
		if !ok {
			return str, errors.New("location param is missing")
		}

		ranges, ok := query["range"]
		if !ok {
			return str, errors.New("range param is missing")
		}

		if len(location) < 2 {
			return str, errors.New("location must has at least 2 values")
		}

		lat, e := strconv.ParseFloat(location[0], 64)
		if e != nil {
			return str, errors.New("location[0] cannot parse to float")
		}

		lon, e := strconv.ParseFloat(location[1], 64)
		if e != nil {
			return str, errors.New(("location[1] cannot parse to float"))
		}

		r, e := strconv.ParseFloat(ranges[0], 64)
		if !ok {
			return str, errors.New("range cannot parse to float")
		}

		str.Lat = lat
		str.Lon = lon
		str.Range = r

		return
	})
	p.First = stage

	res.Error(p.Run())

	if res.Status {
	}
}
