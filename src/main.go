package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	_ "example.com/m/v2/himompkg"
	"example.com/m/v2/router"
	"github.com/gorilla/mux"
)

var Router *mux.Router = mux.NewRouter()
var Server http.Server

func init() {
}
func himom(w http.ResponseWriter, req *http.Request) {
	log.Println("Hi mom")
	w.Write([]byte("Hi mom"))

}
func main() {
	var wait time.Duration

	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()

	Router.HandleFunc("/himom", himom)
	router.Handle(Router)

	Server := &http.Server{
		Addr:         "0.0.0.0:9000",
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 30,
		IdleTimeout:  time.Second * 60,
		Handler:      Router,
	}

	go func() {
		if err := Server.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	<-c
	ctx, cancel := context.WithTimeout(context.Background(), wait)

	defer cancel()

	Server.Shutdown(ctx)
	log.Println("shutting down")

	os.Exit(0)
}
