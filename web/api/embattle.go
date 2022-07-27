package api

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/lonng/nex"
	"monster/db"
	"monster/helper"
)

func MakeEmbattleService() http.Handler { //布阵
	router := mux.NewRouter()

	router.Handle("/embattle/getEmbattle", nex.Handler(getEmbattleHandler)).Methods("POST")
	router.Handle("/embattle/getCurrentEmbattle", nex.Handler(getCurrentEmbattleHandler)).Methods("POST")
	router.Handle("/embattle/updateEmbattle", nex.Handler(updateEmbattleHandler)).Methods("POST")
	router.Handle("/embattle/useEmbattle", nex.Handler(useEmbattleHandler)).Methods("POST")
	return router
}

func getEmbattleHandler(r *http.Request) (map[string]interface{}, error) {
	uid, isValid := helper.VerifyJWT(r)
	if !isValid {
		return nil, errors.New("Invalid token")
	}
	reqJSON := helper.ReadParameters(r)
	embattleInfo := db.GetEmbattle(uid, reqJSON["length"], reqJSON["teamId"])

	payload := map[string]interface{}{
		"embattleInfo": embattleInfo,
	}
	return payload, nil
}

func getCurrentEmbattleHandler(r *http.Request) (map[string]interface{}, error) {
	uid, isValid := helper.VerifyJWT(r)
	if !isValid {
		return nil, errors.New("Invalid token")
	}
	reqJSON := helper.ReadParameters(r)
	embattleInfo := db.GetCurrentEmbattle(uid, reqJSON["length"])

	payload := map[string]interface{}{
		"embattleInfo": embattleInfo,
	}
	return payload, nil
}

func updateEmbattleHandler(r *http.Request) (map[string]string, error) {
	uid, isValid := helper.VerifyJWT(r)
	if !isValid {
		return nil, errors.New("Invalid token")
	}
	reqJSON := helper.ReadParameters(r)
	r1 := strings.Replace(reqJSON["monsterId"], "[", "", -1)
	r2 := strings.Replace(r1, "]", "", -1)
	result := strings.Split(r2, " ")

	r3 := strings.Replace(reqJSON["monsterUid"], "[", "", -1)
	r4 := strings.Replace(r3, "]", "", -1)
	result1 := strings.Split(r4, " ")
	for i := 0; i < len(result); i++ {
		db.UpdateEmbattle(uid, reqJSON["length"], reqJSON["teamid"], result[i], result1[i], strconv.Itoa(i+1))
	}

	payload := map[string]string{
		"code": "success",
	}

	return payload, nil
}

func useEmbattleHandler(r *http.Request) (map[string]string, error) {
	uid, isValid := helper.VerifyJWT(r)
	if !isValid {
		return nil, errors.New("Invalid token")
	}
	reqJSON := helper.ReadParameters(r)
	db.UseEmbattle(uid, reqJSON["length"], reqJSON["teamid"])

	payload := map[string]string{
		"code": "success",
	}

	return payload, nil
}
