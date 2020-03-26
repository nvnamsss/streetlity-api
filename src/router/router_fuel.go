package router

import (
	"log"
	"net/http"

	"example.com/m/v2/model"

	"github.com/gorilla/mux"
)

func getFuels(w http.ResponseWriter, req *http.Request) {
	var f model.Fuel

	model.Db.Where(&model.Fuel{Id: 1}).First(&f)

	// model.Db.First(&f, 1)
	log.Println(f)
	w.Write([]byte(string(f.Id)))
}

func HandleFuel(router *mux.Router) {
	log.Println("[Router]", "Handling fuel")
	s := router.PathPrefix("/fuel").Subrouter()
	s.HandleFunc("/getFuels", getFuels).
		Methods("GET")
}
