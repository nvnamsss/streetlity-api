package stages

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/nvnamsss/goinf/pipeline"
)

func SetOwnerValidate(req *http.Request) *pipeline.Stage {
	req.ParseForm()
	stage := pipeline.NewStage(func() (str struct {
		ServiceId int64
		Owner     string
	}, e error) {
		form := req.PostForm
		service_ids, ok := form["service_id"]
		if !ok {
			return str, errors.New("service_id param is missing")
		}
		owners, ok := form["owner"]
		if !ok {
			return str, errors.New("owner param is missing")
		}

		if service_id, e := strconv.ParseInt(service_ids[0], 10, 64); e != nil {
			return str, errors.New("service_id cannot parse to int64")
		} else {
			str.ServiceId = service_id
		}
		str.Owner = owners[0]
		return
	})

	return stage
}

func QueryMaintenanceValidate(req *http.Request) *pipeline.Stage {
	stage := pipeline.NewStage(func() (str struct {
		Case    int
		Id      int64
		Lat     float64
		Lon     float64
		Address string
	}, e error) {
		query := req.URL.Query()
		cases, ok := query["case"]
		if !ok {
			return str, errors.New("case param is missing")
		}
		if c, e := strconv.Atoi(cases[0]); e != nil {
			return str, errors.New("case param cannot parse to int")
		} else {
			str.Case = c
		}

		switch str.Case {
		case 1:
			ids, ok := query["id"]
			if !ok {
				return str, errors.New("id param is missing")
			}
			if id, e := strconv.ParseInt(ids[0], 10, 64); e != nil {
				return str, errors.New("id param cannot parse to int64")
			} else {
				str.Id = id
			}
			break
		case 2:
			lats, ok := query["lat"]
			if !ok {
				return str, errors.New("lat param is missing")
			}
			lons, ok := query["lon"]
			if !ok {
				return str, errors.New("lon param is missing")
			}

			if lat, e := strconv.ParseFloat(lats[0], 64); e != nil {
				return str, errors.New("cannot parse lat param to float64")
			} else {
				str.Lat = lat
			}
			if lon, e := strconv.ParseFloat(lons[0], 64); e != nil {
				return str, errors.New("cannot parse lon param to float64")
			} else {
				str.Lon = lon
			}
			break
		case 3:
			addresses, ok := query["address"]
			if !ok {
				return str, errors.New("address param is missing")
			}
			str.Address = addresses[0]
			break
		}

		return
	})
	return stage
}
