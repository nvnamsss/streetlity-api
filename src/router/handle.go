package router

import (
	"github.com/gorilla/mux"
)

func Handle(router *mux.Router) {

	HandleService(router)
	HandlePing(router)
}
