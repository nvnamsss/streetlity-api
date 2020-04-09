package router

import (
	"net/http"

	"github.com/gorilla/mux"
)

func allService(w http.ResponseWriter, res *http.Request) {
	
}

func HandleService(router *mux.Router) {
	s := router.PathPrefix("/service").Subrouter()
	s.HandleFunc("/all", allService).Methods("GET")
}
