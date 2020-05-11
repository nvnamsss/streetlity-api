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
	s := pipeline.NewStage(func() (str struct {
		Address string
		Note    string
	}, e error) {
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

		_, ok := form["note"]
		if ok {
			str.Note = form["note"][0]
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

func InRangeServiceValidateStage(req *http.Request) *pipeline.Stage {
	validateParamsStage := pipeline.NewStage(func() error {
		query := req.URL.Query()

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
		query := req.URL.Query()

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

	return validateParamsStage
}
