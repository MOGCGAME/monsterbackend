package api

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"monster/db"
	"monster/helper"

	"github.com/gorilla/mux"
	"github.com/lonng/nex"
	log "github.com/sirupsen/logrus"
)

func MakeMonsterService() http.Handler { //怪兽
	router := mux.NewRouter()

	router.Handle("/monster/getMonster", nex.Handler(getMonsterHandler)).Methods("POST")
	router.Handle("/monster/getMonsterDetail", nex.Handler(getMonsterDetailHandler)).Methods("POST")
	router.Handle("/monster/showMonsterDetail", nex.Handler(showMonsterDetailHandler)).Methods("POST")
	router.Handle("/monster/updateMonsterLevel", nex.Handler(updateMonsterLevel)).Methods("POST")
	router.Handle("/monster/updateMonsterEnergy", nex.Handler(updateMonsterEnergy)).Methods("POST")
	return router
}

func updateMonsterEnergy(r *http.Request) (int, error) {
	uid, isValid := helper.VerifyJWT(r)
	if !isValid {
		return 0, errors.New("Invalid token")
	}
	reqJSON := helper.ReadParameters(r)
	db.UpdateMonsterEnergy(uid, reqJSON["monsterUId"], reqJSON["monsterEnergy"])

	return 0, nil
}

func updateMonsterLevel(r *http.Request) (int, error) {
	uid, isValid := helper.VerifyJWT(r)
	if !isValid {
		return 0, errors.New("Invalid token")
	}
	reqJSON := helper.ReadParameters(r)
	levelAndExp := db.GetExperienceAndLevel(uid, reqJSON["monsterUId"])
	if levelAndExp != nil {
		level, err := strconv.Atoi(levelAndExp["level"])
		if err != nil {
			log.Println(err)
		}
		if level >= 5 {
			return 5, err
		}
		currExp, err := strconv.ParseFloat(levelAndExp["experience"], 64)
		if err != nil {
			log.Println(err)
		}
		expToGain, err := strconv.ParseFloat(reqJSON["monsterExp"], 64)
		if err != nil {
			log.Println("error 2nd")
			log.Println(err)
		}
		newExperience := currExp + expToGain
		fmt.Println(newExperience)
		levelAndUpperLimitMap := map[int]float64{
			0: 2000, 1: 4000, 2: 8000, 3: 16000, 4: 32000,
		}
		newLevel := level
		if newExperience >= levelAndUpperLimitMap[level] {
			newLevel = level + 1
			newExperience = 0
		}
		newLevelAsString := strconv.Itoa(newLevel)
		newExperienceAsString := strconv.FormatFloat(newExperience, 'f', -1, 64)
		db.UpdateMonsterLevelAndExp(uid, reqJSON["monsterUId"], newLevelAsString, newExperienceAsString)
		return newLevel, err
	}

	return 0, nil
}

func getMonsterHandler(r *http.Request) (map[string]interface{}, error) {
	uid, isValid := helper.VerifyJWT(r)
	if !isValid {
		return nil, errors.New("Invalid token")
	}
	reqJSON := helper.ReadParameters(r)

	monster, err := db.GetMonster(uid, reqJSON["teamId"], reqJSON["teamLength"])
	if err != nil {
		log.Println(err)
	}
	fmt.Println("monster:", monster)
	payload := map[string]interface{}{
		"monster": monster,
	}

	return payload, nil
}

func getMonsterDetailHandler(r *http.Request) (map[string]interface{}, error) {
	uid, isValid := helper.VerifyJWT(r)
	reqJSON := helper.ReadParameters(r)
	if !isValid {
		return nil, errors.New("Invalid token")
	}
	monsterDetail, err := db.GetMonsterDetail(uid, reqJSON["sqlString"])
	if err != nil {
		log.Println(err)
	}
	payload := map[string]interface{}{
		"monsterDetail": monsterDetail,
	}
	return payload, nil
}

func showMonsterDetailHandler(r *http.Request) (map[string]string, error) {
	// uid, isValid := helper.VerifyJWT(r)
	// if !isValid {
	// 	return nil, errors.New("Invalid token")
	// }
	// fmt.Println(uid)
	reqJSON := helper.ReadParameters(r)
	monsterDetail, err := db.GetMonsterDetailByMonsterUid(reqJSON["monsteruid"])
	if err != nil {
		log.Println(err)
	}
	return monsterDetail, nil
}

// func updateMonsterLevel(r *http.Request, req *proto.UpdateMonsterExperienceReq) (int, error) {
// 	_, isValid := helper.VerifyJWT(r)
// 	if !isValid {
// 		return 0, errors.New("Invalid token")
// 	}
// 	var exp, monsteruid, currExp, newExp int
// 	exp = int(req.MonsterExp)
// 	monsteruid = int(req.MonsterUId)

// 	levelAndExp := db.GetMonsterExpAndLevel(monsteruid)
// 	if levelAndExp != nil {
// 		level, err := strconv.Atoi((levelAndExp["lv"]))
// 		if err != nil {
// 			log.Println(err)
// 		}
// 		if level >= 5 {
// 			return 5, err
// 		}
// 		currExp = int(levelAndExp["exp"])
// 		expToGain := exp
// 		newExp = currExp + expToGain
// 		levelAndUpperLimitMap := map[int]int{
// 			0: 100, 1: 400, 3: 800, 4: 1600, 5: 3200,
// 		}
// 		newLevel := level
// 		if newExp >= levelAndUpperLimitMap[level] {
// 			newLevel = level + 1
// 			newExp = 0
// 		}

// 		db.UpdateMonsterLevelAndExp(monsteruid, newLevel, newExp)
// 	}

// 	return 0, nil
// }
