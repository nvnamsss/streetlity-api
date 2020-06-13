package middleware

import (
	"log"
	"net/http"
	"streelity/v1/model"
	"streelity/v1/router/sres"
)

//Authenticate middleware
//
//Request must have `Auth` header to be passed
func Authenticate(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Auth")
		if auth != "" {
			log.Println("[Authorization]", "Received key", auth)
			err := model.Authenticate(auth)
			if err != nil {
				var res sres.Response = sres.Response{Status: false, Message: err.Error()}
				w.WriteHeader(http.StatusUnauthorized)
				sres.WriteJson(w, res)
			} else {
				h.ServeHTTP(w, r)
			}
		} else {
			var res sres.Response = sres.Response{Status: false, Message: "Authorization failure."}
			log.Println("[Authorization]", r.URL, "Authorization failure")
			sres.WriteJson(w, res)
		}

	})
}
