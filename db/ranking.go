package db

import (
	"log"
)

func GetPvPRanking(mode string) ([]map[string]string, error) { // new
	ranking, err := db.QueryString("SELECT * FROM `user` u ORDER BY u.rank1 DESC")
	if mode == "5" {
		ranking, err = db.QueryString("SELECT * FROM `user` u ORDER BY u.rank2 DESC")
	}

	if err != nil {
		log.Println(err)
	}

	return ranking, nil
}
