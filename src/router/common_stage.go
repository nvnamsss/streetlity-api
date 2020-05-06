package router

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/nvnamsss/goinf/pipeline"
)

/*Containing commom stage for pipeline, in order to reduce the effort and lines of code to implement
the pipeline in request handle*/

//ServiceValidateStage create the validated stage for adding a new service
func AddingServiceValidateStage(req *http.Request) *pipeline.Stage {
	s := pipeline.NewStage(func() (str struct{ Address string }, e error) {
		form := req.PostForm
		location, locationOk := form["location"]
		_, addressOk := form["address"]
		if !locationOk {
			return str, errors.New("location param is missing")
		} else {
			if len(location) < 2 {
				return str, errors.New("location param must have 2 values")
			}
		}

		if !addressOk {
			return str, errors.New("address param is missing")
		} else {
			str.Address = form["address"][0]
		}

		return
	})

	return s
}

//AddingServiceParsingStage create the parsing stage for adding a new service
func AddingServiceParsingStage(req *http.Request) *pipeline.Stage {
	s := pipeline.NewStage(func() (str struct {
		Lat float64
		Lon float64
	}, e error) {
		var (
			latErr error
			lonErr error
		)
		form := req.PostForm

		str.Lat, latErr = strconv.ParseFloat(form["location"][0], 64)
		str.Lon, lonErr = strconv.ParseFloat(form["location"][1], 64)
		if latErr != nil {
			return str, errors.New("cannot parse location[0] to float")
		}
		if lonErr != nil {
			return str, errors.New("cannot parse location[1] to float")
		}
		return str, nil
	})

	return s
}
