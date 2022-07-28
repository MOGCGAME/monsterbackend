package db

import (
	"log"
)

func GetEmbattle(uid, length, teamid string) []map[string]string {

	embattle, err := db.QueryString("SELECT * FROM `user_embattle` WHERE user_id = " + uid + " AND length = " + length + " AND team_id = " + teamid)
	if err != nil {
		log.Println(err)
	}

	return embattle
}

func GetCurrentEmbattle(uid, length string) []map[string]string {
	embattle, err := db.QueryString("SELECT * FROM `user_embattle` WHERE user_id = " + uid + " AND length = " + length + " AND current = 1")
	if err != nil {
		log.Println(err)
	}

	return embattle
}

func UpdateEmbattle(uid, length, teamid, monsteruid, monsterid, seqid string) { //根据teamid分配布阵
	embattle, err := db.QueryString("SELECT * FROM `user_embattle` WHERE user_id = " + uid +
		" AND length = " + length + " AND team_id = " + teamid + " AND sequence_id = " + seqid)
	if err != nil {
		log.Println(err)
	}
	if len(embattle) > 0 { //更新站位
		_, err := db.Exec("Update `user_embattle` SET user_monster_uid = " + monsteruid + ", user_monster_id = " + monsterid +
			" WHERE user_id = " + uid + " AND length = " + length + " AND team_id = " + teamid + " AND sequence_id = " + seqid)
		if err != nil {
			log.Println(err)
		}
	} else { //添加新站位
		_, err := db.Exec("INSERT INTO `user_embattle` (`user_id`, `team_id`, `user_monster_uid`, `user_monster_id`, `sequence_id`, `length`) VALUES (?, ?, ?, ?, ?, ?)",
			uid, teamid, monsteruid, monsterid, seqid, length)
		if err != nil {
			log.Println(err)
		}
	}
}

func UseEmbattle(uid, length, teamid string) {
	_, err := db.Exec("Update `user_embattle` SET current = 0 WHERE user_id = " + uid + " AND length = " + length + " AND team_id != " + teamid)
	if err != nil {
		log.Println(err)
	}

	_, err = db.Exec("Update `user_embattle` SET current = 1 WHERE user_id = " + uid + " AND length = " + length + " AND team_id = " + teamid)
	if err != nil {
		log.Println(err)
	}
}
