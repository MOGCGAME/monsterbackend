package db

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"time"

	"monster/helper"
)

func GetRanking(uid string, rankNum string) (int, error) {
	rank, err := db.QueryString("SELECT * FROM `user` WHERE uid = " + uid + " LIMIT 1")
	if err != nil {
		log.Println(err)
	}
	if len(rank) > 0 {
		var ranking int
		if rankNum == "3" {
			ranking = helper.StringToInt(rank[0]["rank1"])

		} else if rankNum == "5" {
			ranking = helper.StringToInt(rank[0]["rank2"])
		}

		return ranking, nil
	}

	return 0, err
}

func GetSelfEmbattle(uid, length string) ([]map[string]string, error) {
	//玩家布阵
	var selfMonster []map[string]string
	selfEmbattle, err := db.QueryString("SELECT * FROM `user_embattle` WHERE user_id = " + uid + " AND current = 1 AND length = " + length)
	for i := 0; i < len(selfEmbattle); i++ {
		monster, err := db.QueryString("SELECT * FROM `monster_info` WHERE uid = " + selfEmbattle[i]["user_monster_uid"])
		monsterseq, err := db.QueryString("SELECT * FROM `user_embattle` WHERE user_monster_uid = " + selfEmbattle[i]["user_monster_uid"])
		monsterInfo := getMonsterInfo(monster, monsterseq)

		if err != nil {
			log.Println(err)
		}
		if monster[0]["energy"] == "0" {
			err1 := errors.New("monster not enough energy")
			var payload []map[string]string
			energy := make(map[string]string)
			energy["monster_id"] = monster[0]["monster_id"]
			payload = append(payload, energy)

			return payload, err1
		}
		selfMonster = append(selfMonster, monsterInfo)
	}
	if err != nil {
		log.Println(err)
	}

	return selfMonster, nil
}

func GetPvPEnemyEmbattle(uid, length, rank1, rank2, match string) []map[string]string {
	//随机匹配同段位玩家
	var enemyUid string
	var enemyMonster []map[string]string
	var enemy int
	var querystring1 string
	var querystring2 string
	if length == "5" { // new
		querystring1 = " AND user.rank2 >= "
		querystring2 = " AND user.rank2 <= "
	} else {
		querystring1 = " AND user.rank1 >= "
		querystring2 = " AND user.rank1 <= "
	}

	var matching []map[string]string

	for stay, timeout := true, time.After(time.Second*30); stay; {
		matching, _ = db.QueryString("SELECT * FROM `user` WHERE uid != " + uid + " AND matching = " + match + querystring1 + rank1 + querystring2 + rank2)
		if len(matching) > 0 {
			stay = false
		}
		select {
		case <-timeout:
			fmt.Println("30 seconds has been reached!")
			stay = false
		default:
		}
	}
	if len(matching) > 0 {
		rand.Seed(time.Now().UnixNano())
		enemy = rand.Intn(len(matching))
		enemyUid = matching[enemy]["uid"]
		enemyEmbattle, err1 := db.QueryString("SELECT * FROM `user_embattle` WHERE user_id = " + enemyUid + " AND current = 1 AND length = " + length)
		if err1 != nil {
			log.Println(err1)
		}
		for i := 0; i < len(enemyEmbattle); i++ {
			monster, err := db.QueryString("SELECT * FROM `monster_info` WHERE uid = " + enemyEmbattle[i]["user_monster_uid"])
			monsterseq, err := db.QueryString("SELECT * FROM `user_embattle` WHERE user_monster_uid = " + enemyEmbattle[i]["user_monster_uid"])
			monsterInfo := getMonsterInfo(monster, monsterseq)

			if err != nil {
				log.Println(err)
			}
			if monster[0]["energy"] == "0" {
				return GetPvPEnemyEmbattle(uid, length, rank1, rank2, match)
			}
			enemyMonster = append(enemyMonster, monsterInfo)
		}
		return enemyMonster
	} else {
		return nil
	}
}

func GetPvEEnemyEmbattle(checkpoint, stage string) ([]map[string]string, error) {
	var enemyUid string
	matching, err := db.QueryString("SELECT * FROM `stage_bot` WHERE check_point = " + checkpoint + " AND stage = " + stage)
	enemyUid = matching[0]["bot_id"]
	enemyEmbattle, err := db.QueryString("SELECT * FROM `bot_embattle` WHERE bot_id = " + enemyUid)
	var enemyMonster []map[string]string
	if len(enemyEmbattle) > 0 {
		for i := 0; i < len(enemyEmbattle); i++ {
			monster, err := db.QueryString("SELECT * FROM `bot_monster_info` WHERE uid = " + enemyEmbattle[i]["user_monster_uid"] + " AND user_id = " + enemyUid)
			monsterseq, err := db.QueryString("SELECT * FROM `bot_embattle` WHERE user_monster_uid = " + enemyEmbattle[i]["user_monster_uid"] + " AND bot_id = " + enemyUid)
			monsterInfo := getMonsterInfo(monster, monsterseq)

			if err != nil {
				log.Println(err)
			}
			enemyMonster = append(enemyMonster, monsterInfo)
		}
	}

	if err != nil {
		log.Println(err)
	}

	return enemyMonster, nil
}

func GetPvPEnemy(self, enemy, length string) ([]map[string]string, error) {
	embattle, err := db.QueryString("SELECT * FROM `user_embattle` WHERE user_id = " + enemy + " AND current = 1 AND length = " + length +
		" OR user_id = " + self + " AND current = 1 AND length = " + length)
	var Monster []map[string]string
	if len(embattle) > 0 {
		for i := 0; i < len(embattle); i++ {
			monster, err := db.QueryString("SELECT * FROM `monster_info` WHERE uid = " + embattle[i]["user_monster_uid"])
			monsterseq, err := db.QueryString("SELECT * FROM `user_embattle` WHERE user_monster_uid = " + embattle[i]["user_monster_uid"] + " AND length = " + length + " AND current = 1")
			monsterskill, err := db.QueryString("SELECT * FROM `monster_skill` WHERE monster_id = " + embattle[i]["user_monster_uid"])
			monsterInfo := getMonsterInfo(monster, monsterseq)

			if len(monsterskill) > 0 {
				skill, err := db.QueryString("SELECT * FROM `skill` WHERE skill = " + monsterskill[0]["skill"])
				if err != nil {
					log.Println(err)
				}
				monsterInfo["monster_skill"] = skill[0]["skill"]
				monsterInfo["monster_trigger"] = skill[0]["trigger"]
			} else {
				monsterInfo["monster_skill"] = "0"
				monsterInfo["monster_trigger"] = "0"
			}

			if err != nil {
				log.Println(err)
			}
			Monster = append(Monster, monsterInfo)
		}
	}

	if err != nil {
		log.Println(err)
	}

	return Monster, nil
}

func GetPvEEnemy(self, enemy, length string) ([]map[string]string, error) {
	embattle, err := db.QueryString("SELECT * FROM `user_embattle` WHERE user_id = " + self + " AND current = 1 AND length = " + length)
	botEmbattle, err := db.QueryString("SELECT * FROm `bot_embattle` WHERE bot_id = " + enemy)
	var Monster []map[string]string

	if len(embattle) > 0 && len(botEmbattle) > 0 {
		for i := 0; i < len(embattle); i++ {
			monster, err := db.QueryString("SELECT * FROM `monster_info` WHERE uid = " + embattle[i]["user_monster_uid"])
			monsterseq, err := db.QueryString("SELECT * FROM `user_embattle` WHERE user_monster_uid = " + embattle[i]["user_monster_uid"] + " AND length = " + length + " AND current = 1")
			monsterskill, err := db.QueryString("SELECT * FROM `monster_skill` WHERE monster_id = " + embattle[i]["user_monster_uid"])
			if err != nil {
				log.Println(err)
			}
			monsterInfo := getMonsterInfo(monster, monsterseq)

			if len(monsterskill) > 0 {
				skill, err := db.QueryString("SELECT * FROM `skill` WHERE skill = " + monsterskill[0]["skill"])
				if err != nil {
					log.Println(err)
				}
				monsterInfo["monster_skill"] = skill[0]["skill"]
				monsterInfo["monster_trigger"] = skill[0]["trigger"]
			} else {
				monsterInfo["monster_skill"] = "0"
				monsterInfo["monster_trigger"] = "0"
			}
			Monster = append(Monster, monsterInfo)
		}
		for j := 0; j < len(botEmbattle); j++ {
			monster, err := db.QueryString("SELECT * FROM `bot_monster_info` WHERE uid = " + botEmbattle[j]["user_monster_uid"] + " AND user_id = " + enemy)
			monsterseq, err := db.QueryString("SELECT * FROM `bot_embattle` WHERE user_monster_uid = " + botEmbattle[j]["user_monster_uid"] + " AND bot_id = " + enemy)
			monsterskill, err := db.QueryString("SELECT * FROM `monster_skill` WHERE monster_id = " + botEmbattle[j]["user_monster_uid"])
			if err != nil {
				log.Println(err)
			}
			monsterInfo := getMonsterInfo(monster, monsterseq)

			if len(monsterskill) > 0 {
				skill, err := db.QueryString("SELECT * FROM `skill` WHERE skill = " + monsterskill[0]["skill"])
				if err != nil {
					log.Println(err)
				}
				monsterInfo["monster_skill"] = skill[0]["skill"]
				monsterInfo["monster_trigger"] = skill[0]["trigger"]
			} else {
				monsterInfo["monster_skill"] = "0"
				monsterInfo["monster_trigger"] = "0"
			}
			Monster = append(Monster, monsterInfo)
		}
	}
	if err != nil {
		log.Println(err)
	}
	return Monster, nil
}

func GetMatch() (map[string]interface{}, error) {

	payload := map[string]interface{}{}

	return payload, nil
}

func GetPlayerCheckPoint(uid string, checkpoint int) (int, error) {

	user, err := db.QueryString("SELECT * FROM `user` WHERE uid = " + uid + " LIMIT 1")
	if err != nil {
		log.Println(err)
	}

	if len(user) > 0 {
		if helper.StringToInt(user[0]["check_point"]) >= checkpoint {
			return checkpoint, nil
		}
	}
	return 0, err
}

func GetPlayerStage(uid string, stage int) (int, error) {

	user, err := db.QueryString("SELECT * FROM `user` WHERE uid = " + uid + " LIMIT 1")
	if err != nil {
		log.Println(err)
	}

	if len(user) > 0 {
		if helper.StringToInt(user[0]["stage"]) >= stage {
			return stage, nil
		}
	}
	return 0, err
}

func GetStageLen(checkpoint string, stage string) (int, error) {
	stage_bot_query, err := db.QueryString("SELECT * FROM stage_bot WHERE check_point = " + checkpoint + " AND stage = " + stage)
	if err != nil {
		log.Println(err)
	}
	stage_bot := stage_bot_query[0]["bot_id"]
	bot_embattle_query, err := db.QueryString("SELECT * FROM bot_embattle WHERE bot_id = " + stage_bot)
	if len(bot_embattle_query) > 0 {
		return len(bot_embattle_query), nil
	} else {
		return 0, err
	}
}

func GetAwardByStage(checkpoint, stage string) (string, error) {
	stageaward, err := db.QueryString("SELECT * FROM `stage_award` WHERE check_point = " + checkpoint + " AND stage = " + stage)

	if err != nil {
		log.Println(err)
	}

	return stageaward[0]["award"], nil
}

func GetExpByStage(checkpoint, stage string) (string, error) {
	stageaward, err := db.QueryString("SELECT * FROM `stage_award` WHERE check_point = " + checkpoint + " AND stage = " + stage)

	if err != nil {
		log.Println(err)
	}

	return stageaward[0]["exp"], nil
}

func GetCurrentStage(uid string) int {
	user, err := db.QueryString("SELECT * FROM `user` WHERE uid = " + uid + " LIMIT 1")
	if err != nil {
		log.Println(err)
	}

	return helper.StringToInt(user[0]["stage"])
}

func RenewMatching(uid, matching string) {
	_, err := db.Exec("UPDATE `user` SET `matching` = ? WHERE `uid` = ?",
		matching, uid)
	if err != nil {
		log.Println(err)
	}
}

func UpdateMatching(uid, matching string) {
	_, err := db.Exec("UPDATE `user` SET `matching` = ? WHERE `uid` = ?",
		matching, uid)
	if err != nil {
		log.Println(err)
	}
}

func UpdateAward(award, uid string) {
	_, err := db.Exec("UPDATE `user` SET `award` = ? WHERE `uid` = ?",
		award, uid)
	if err != nil {
		log.Println(err)
	}
}

func UpdateRank(rank, uid string, rankNum string) { // new
	switch rankNum {
	case "3":
		_, err := db.Exec("UPDATE `user` SET `rank1` = ? WHERE `uid` = ?",
			rank, uid)
		if err != nil {
			log.Println(err)
		}
	case "5":
		_, err := db.Exec("UPDATE `user` SET `rank2` = ? WHERE `uid` = ?",
			rank, uid)
		if err != nil {
			log.Println(err)
		}
	}

}

func UpdateRank1(rank, uid string) { // new
	_, err := db.Exec("UPDATE `user` SET `rank1` = `rank1` + ? WHERE `uid` = ?",
		rank, uid)
	if err != nil {
		log.Println(err)
	}
}

func UpdateRank2(rank, uid string) { // new
	_, err := db.Exec("UPDATE `user` SET `rank2` = `rank2` + ? WHERE `uid` = ?",
		rank, uid)
	if err != nil {
		log.Println(err)
	}
}

func UpdateCheckPoint(uid string) {
	_, err := db.Exec("UPDATE `user` SET `checkpoint` = `checkpoint` + 1 WHERE uid = " + uid)
	if err != nil {
		log.Println(err)
	}
}

func UpdateStage(uid string, stage int) {
	if stage == 1 {
		_, err := db.Exec("UPDATE `user` SET `stage` = 1 WHERE uid = " + uid)
		if err != nil {
			log.Println(err)
		}
	} else {
		_, err := db.Exec("UPDATE `user` SET `stage` = `stage` + 1 WHERE uid = " + uid)
		if err != nil {
			log.Println(err)
		}
	}

}

func getMonsterInfo(monster []map[string]string, monsterseq []map[string]string) map[string]string {
	monsterInfo := make(map[string]string)
	monsterInfo["user_id"] = monster[0]["user_id"]
	monsterInfo["monster_id"] = monster[0]["monster_id"]
	monsterInfo["monster_uid"] = monster[0]["uid"]
	monsterInfo["monster_name"] = monster[0]["name"]
	monsterInfo["monster_rarity"] = monster[0]["rarity"]
	monsterInfo["monster_element"] = monster[0]["element"]
	monsterInfo["monster_hp"] = monster[0]["hp"]
	monsterInfo["monster_max_hp"] = monster[0]["hp"]
	monsterInfo["monster_attack"] = monster[0]["attack"]
	monsterInfo["monster_defend"] = monster[0]["defend"]
	monsterInfo["monster_speed"] = monster[0]["speed"]
	monsterInfo["monster_hit"] = "100"
	monsterInfo["monster_miss"] = "20"
	monsterInfo["monster_skill_rate"] = "20"
	monsterInfo["monster_energy"] = monster[0]["energy"]
	monsterInfo["monster_sequence"] = monsterseq[0]["sequence_id"]
	monsterInfo["monster_initpos"] = ""
	monsterInfo["monster_initneg"] = ""
	monsterInfo["monster_positive"] = ""
	monsterInfo["monster_negative"] = ""
	return monsterInfo
}
