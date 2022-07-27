package api

import (
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/lonng/nex"
	log "github.com/sirupsen/logrus"
	"monster/db"
	"monster/helper"
)

func MakePropService() http.Handler { //背包
	router := mux.NewRouter()
	router.Handle("/prop/getProp", nex.Handler(getPropHandler)).Methods("POST")
	router.Handle("/prop/useProp", nex.Handler(usePropHandler)).Methods("POST")
	return router
}

func getPropHandler(r *http.Request) (map[string]interface{}, error) {
	uid, isValid := helper.VerifyJWT(r)
	if !isValid {
		return nil, errors.New("Invalid token")
	}
	propInfo, err := db.GetProp(uid)
	if err != nil {
		log.Println(err)
	}

	payload := map[string]interface{}{
		"propInfo": propInfo,
	}

	return payload, nil
}

func usePropHandler(r *http.Request) (map[string]interface{}, error) {

	payload := map[string]interface{}{
		"info": 1,
	}
	return payload, nil
}
