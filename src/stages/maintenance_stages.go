package stages

import (
	"errors"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/nvnamsss/goinf/pipeline"
)

func AddMaintainerValidate(req *http.Request) *pipeline.Stage {
	req.ParseForm()
	stage := pipeline.NewStage(func() (str struct {
		ServiceId  int64
		Maintainer string
	}, e error) {
		form := req.PostForm
		log.Println(form)
		service_ids, ok := form["service_id"]
		if !ok {
			return str, errors.New("service_id param is missing")
		}
		maintainers, ok := form["maintainer"]
		if !ok {
			return str, errors.New("maintainer param is missing")
		}

		if service_id, e := strconv.ParseInt(service_ids[0], 10, 64); e != nil {
			return str, errors.New("service_id cannot parse to int64")
		} else {
			str.ServiceId = service_id
		}
		str.Maintainer = maintainers[0]
		return
	})

	return stage
}

func RemoveMaintainerValidate(req *http.Request) *pipeline.Stage {
	stage := pipeline.NewStage(func() (str struct {
		ServiceId  int64
		Maintainer string
	}, e error) {
		query, _ := url.ParseQuery(req.URL.RawQuery)
		log.Println(query)
		service_ids, ok := query["service_id"]
		if !ok {
			return str, errors.New("service_id param is missing")
		}
		maintainers, ok := query["maintainer"]
		if !ok {
			return str, errors.New("maintainer param is missing")
		}

		if service_id, e := strconv.ParseInt(service_ids[0], 10, 64); e != nil {
			return str, errors.New("service_id cannot parse to int64")
		} else {
			str.ServiceId = service_id
		}
		str.Maintainer = maintainers[0]
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

func UpdateMaintenanceValidate(req *http.Request) *pipeline.Stage {
	req.ParseForm()
	stage := pipeline.NewStage(func() (str struct {
		Lat    float32
		Lon    float32
		Note   string
		Images []string
	}, e error) {
		form := req.PostForm
		if images, ok := form["images"]; ok {
			str.Images = images
		}
		if _, ok := form["lat"]; ok {
			if lat, e := strconv.ParseFloat(form["lat"][0], 64); e != nil {
				return str, errors.New("cannot parse lat to float32")
			} else {
				str.Lat = float32(lat)
			}
		}

		if _, ok := form["lon"]; ok {
			if lon, e := strconv.ParseFloat(form["lon"][0], 64); e != nil {
				return str, errors.New("cannot parse lon to float32")
			} else {
				str.Lon = float32(lon)
			}
		}

		if notes, ok := form["note"]; ok {
			str.Note = notes[0]
		}

		return
	})

	return stage
}
