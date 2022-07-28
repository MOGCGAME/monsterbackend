package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

func MakeEquipmentService() http.Handler { //装备
	router := mux.NewRouter()
	return router
}
