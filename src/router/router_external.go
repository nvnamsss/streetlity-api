package router

import (
	"net/http"

	"github.com/gorilla/mux"
)

///1 lat = 111km
func external(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("External"))
}

func Handle(router *mux.Router) {
	s := router.PathPrefix("/external").Subrouter()
	s.HandleFunc("/", external).
		Methods("GET", "POST")

	HandleAuth(router)
	HandleFuel(router)
	HandlePing(router)
}
