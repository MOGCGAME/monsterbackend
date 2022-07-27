package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

func MakeBattleService() http.Handler { //怪兽
	router := mux.NewRouter()

	return router
}
