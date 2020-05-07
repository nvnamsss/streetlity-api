package router

import (
	"fmt"
	"log"
	"net/http"
	"streelity/v1/model"
)

type Adapter func(http.Handler) http.Handler

func Adapt(h http.Handler, adapters ...Adapter) http.Handler {
	for _, adapter := range adapters {
		h = adapter(h)
	}
	return h
}

func Notify(logger *log.Logger) Adapter {

	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("before")
			defer fmt.Println("after")

			h.ServeHTTP(w, r)
		})
	}
}

//Authenticate middleware
//
//Request must have `Auth` header to be passed
func Authenticate(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Auth")
		if auth != "" {
			log.Println("[Authorization]", "Received key", auth)
			err := model.Auth(auth)
			if err != nil {
				var res Response = Response{Status: false, Message: err.Error()}
				w.WriteHeader(http.StatusUnauthorized)
				WriteJson(w, res)
			} else {
				h.ServeHTTP(w, r)
			}
		} else {
			var res Response = Response{Status: false, Message: "Authorization failure."}
			log.Println("[Authorization]", r.URL, "Authorization failure")
			WriteJson(w, res)
		}

	})
}

func Middleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("middleware", r.URL)
		h.ServeHTTP(w, r)
	})
}
