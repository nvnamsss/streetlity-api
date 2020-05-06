package router

import (
	"log"

	"github.com/gorilla/mux"
)

func HandleService(router *mux.Router) {
	log.Println("[Router]", "Handling service")

	s := router.PathPrefix("/service").Subrouter()

	HandleFuel(s)
	HandleAtm(s)
	HandleToilet(s)
	HandleMaintainer(s)

	Versioning(s, "1.0.0", "2.1.0")
}
