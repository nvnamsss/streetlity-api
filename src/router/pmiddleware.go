package router

import (
	"errors"
	"net/http"
	"regexp"

	"github.com/gorilla/mux"
)

//sumOfRunes return the sum of runes by int in a string
func sumOfRunes(s string) (sum int) {
	for _, r := range s {
		sum += int(r)
	}

	return
}

//compareVersion compare the
func compareVersion(currentVersion, minVersion, maxVersion string) (status bool, err error) {
	regFormat, _ := regexp.Compile("\\d*?[.]\\d*")
	reg, _ := regexp.Compile("[.]")
	if !regFormat.Match([]byte(currentVersion)) {
		return false, errors.New("currentVersion format is invalid")
	}

	if !regFormat.Match([]byte(minVersion)) {
		return false, errors.New("currentVersion format is invalid")
	}

	if !regFormat.Match([]byte(maxVersion)) {
		return false, errors.New("currentVersion format is invalid")
	}

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
	minValue := sumOfRunes(min[0])*100 + sumOfRunes(min[1])*10 + sumOfRunes(min[2])
	currentValue := sumOfRunes(current[0])*100 + sumOfRunes(current[1])*10 + sumOfRunes(current[2])
	maxValue := sumOfRunes(max[0])*100 + sumOfRunes(max[1])*10 + sumOfRunes(max[2])

	if currentValue >= minValue && currentValue <= maxValue {
		return true, nil
	}

	return false, errors.New("This version is not supported")
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
