package router

import (
	"errors"
	"net/http"
	"regexp"

	"github.com/gorilla/mux"
)

func compareVersion(currentVersion, minVersion, maxVersion string) (status bool, err error) {
	reg, _ := regexp.Compile(".")
	current := reg.Split(currentVersion, -1)
	min := reg.Split(minVersion, -1)
	max := reg.Split(maxVersion, -1)

	if len(current) < 3 {
		return false, errors.New("currentVersion format is invalid")
	}

	if len(min) < 3 {
		return false, errors.New("minVersion format is invalid")
	}

	if len(max) < 3 {
		return false, errors.New("maxVersion format is invalid")
	}

	switch {
	case (min[0] < current[0] && max[0] > current[0]):
		return true, nil
	case min[0] > current[0] || max[0] < current[1]:
		return false, errors.New("This version is not supported")
	case min[1] < current[0] && max[1] > current[1]:
		return true, nil
	case min[1] > current[1] || max[1] < current[1]:
		return false, errors.New("This version is not supported")
	case min[2] < current[2] && max[2] > current[2]:
		return true, nil
	case min[2] > current[2] || max[2] < current[2]:
		return false, errors.New("This version is not supported")
	}

	return false, errors.New("I do not intent to be here " + currentVersion + " " + minVersion + " " + maxVersion)
}

//Versioning middleware
//
//Only accept the request with in the version range
func Versioning(router *mux.Router, minVersion string, maxVersion string) {
	middleware := func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			version := r.Header.Get("Version")

			if version != "" {
				var res Response
				var err error
				res.Status, err = compareVersion(version, minVersion, maxVersion)

				if err != nil {
					res.Error(err)
					WriteJson(w, res)
				} else {
					h.ServeHTTP(w, r)
				}
			} else {
				var res Response = Response{Status: false, Message: "Version is missing"}
				WriteJson(w, res)
			}
		})
	}

	router.Use(middleware)
}
