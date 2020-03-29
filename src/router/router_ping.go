package router

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func ping(w http.ResponseWriter, req *http.Request) {
	log.Println("Ping")
	w.Write([]byte("Ping"))

}

func HandlePing(router *mux.Router) {
	log.Println("[Router]", "Handling ping")
	s := router.PathPrefix("/ping").Subrouter()
	s.HandleFunc("/", ping).
		Methods("GET", "POST")
}
