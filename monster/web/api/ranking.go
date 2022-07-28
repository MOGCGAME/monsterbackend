package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/lonng/nex"
	log "github.com/sirupsen/logrus"
	"monster/db"
	"monster/helper"
)

func MakeRankingService() http.Handler { //排行榜
	router := mux.NewRouter()
	router.Handle("/ranking/getPvPRanking", nex.Handler(getPvPRankingHandler)).Methods("POST")
	return router
}

func getPvPRankingHandler(r *http.Request) (map[string]interface{}, error) { // new
	reqJSON := helper.ReadParameters(r)
	ranking, err := db.GetPvPRanking(reqJSON["mode"])
	if err != nil {
		log.Println(err)
	}
	payload := map[string]interface{}{
		"ranking": ranking,
	}

	return payload, nil
}
