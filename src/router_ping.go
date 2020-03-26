package main

import (
	"log"
	"net/http"
)

func ping(w http.ResponseWriter, req *http.Request) {
	log.Println("Ping")
	w.Write([]byte("Ping"))

}

func init() {
	log.Println("init test")
	s := Router.PathPrefix("/ping").Subrouter()
	s.HandleFunc("/", ping).
		Methods("GET", "POST")
}
