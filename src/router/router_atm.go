package router

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func getAtms(w http.ResponseWriter, req *http.Request) {

}

func updateAtm(w http.ResponseWriter, req *http.Request) {

}

func getAtmById(w http.ResponseWriter, req *http.Request) {

}

func HandleAtm(router *mux.Router) {
	log.Println("[Router]", "Handling Atm")
	s := router.PathPrefix("/atm").Subrouter()
	s.HandleFunc("/all", getAtms).Methods("GET")
	s.HandleFunc("/update", updateAtm).Methods("POST")
	s.HandleFunc("/id", getFuel).Methods("GET")
}
