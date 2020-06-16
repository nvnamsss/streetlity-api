package router

import (
	"log"
	"streelity/v1/middleware"

	"github.com/gorilla/mux"
)

func HandleService(router *mux.Router) {
	log.Println("[Router]", "Handling service")

	s := router.PathPrefix("/service").Subrouter()

	HandleFuel(s)
	HandleAtm(s)
	HandleToilet(s)
	HandleMaintenance(s)

	middleware.Versioning(s, "1.0.0", "2.1.0")
}
