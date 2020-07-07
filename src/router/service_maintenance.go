package router

import (
	"log"
	"streelity/v1/router/rmaintenance"

	"github.com/gorilla/mux"
)

func HandleMaintenance(router *mux.Router) {
	log.Println("[Router]", "Handling fuel")
	s := rmaintenance.HandleService(router)
	rmaintenance.HandleReview(s)
	rmaintenance.HandleUnconfirmed(router)
	rmaintenance.HandleHistory(s)
	rmaintenance.HandleOrder(s)
	// s.HandleFunc("/order", orderMaintenance).Methods("POST")
	// s.HandleFunc("/accept", acceptOrderMaintenance).Methods("POST")
}
