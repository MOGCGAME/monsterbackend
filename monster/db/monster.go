package db

import (
	"log"
	"strconv"
)

func UpdateMonsterEnergy(uid, monsteruid, energy string) {
	_, err := db.Exec("UPDATE `monster_info` SET `energy` = `energy` + " + energy + " WHERE `user_id` = " + uid + " AND `uid` =  " + monsteruid)
	if err != nil {
		log.Println(err)
	}
}

func UpdateMonsterLevelAndExp(uid, monsteruid, lv, exp string) {
	_, err := db.Exec("Update `monster_info` set `lv` = '" + lv + "', `exp` = '" + exp + "' WHERE `uid` = '" + monsteruid + "' AND `user_id` = '" + uid + "'")
	if err != nil {
		log.Println(err)
	}

}

func GetExperienceAndLevel(uid, monsteruid string) map[string]string {
	messages, err := db.QueryString("select * from `monster_info` where `user_id` = '" + uid + "' AND `uid` = '" + monsteruid + "' Limit 1")
	if err != nil {
		log.Println(err)
	}
	if len(messages) > 0 {
		payload := map[string]string{
			"level":      messages[0]["lv"],
			"experience": messages[0]["exp"],
		}
		return payload
	} else {
		return nil
	}
}

func GetMonster(uid, teamId string, teamLength string) ([]map[string]string, error) {
	var monster []map[string]string
	messages, err := db.QueryString("SELECT * FROM `monster_info` WHERE user_id = " + uid)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	for i := 0; i < len(messages); i++ {
		monsterseq, err := db.QueryString("SELECT * FROM `user_embattle` WHERE user_monster_uid = " + messages[i]["uid"] + " AND team_id = " + teamId + " AND length = " + teamLength)
		monsterInfo := make(map[string]string)
		monsterInfo["monster_id"] = messages[i]["monster_id"]
		monsterInfo["monster_uid"] = messages[i]["uid"]
		monsterInfo["monster_name"] = messages[i]["name"]
		monsterInfo["monster_rarity"] = messages[i]["rarity"]
		monsterInfo["monster_element"] = messages[i]["element"]
		monsterInfo["monster_energy"] = messages[i]["energy"]
		if len(monsterseq) > 0 {
			monsterInfo["monster_sequence"] = monsterseq[0]["sequence_id"]
		}
		if err != nil {
			log.Println(err)
		}
		monster = append(monster, monsterInfo)
	}

	return monster, nil
}

func GetMonsterDetail(uid string, sqlString string) ([]map[string]string, error) {
	var monsterDetail []map[string]string

	monsterlist, err := db.QueryString("SELECT * FROM `monster_info` WHERE user_id = " + uid)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	var monsters = monsterlist
	if sqlString != "" {
		monsterlist, err := db.QueryString("SELECT * FROM `monster_info` WHERE user_id = " + uid + " " + sqlString)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		monsters = monsterlist
	}

	for i := 0; i < len(monsters); i++ {
		monsterInfo := make(map[string]string)
		monsterskill, err := db.QueryString("SELECT * FROM `monster_skill` WHERE monster_id = " + monsters[i]["uid"])
		monsterInfo["monster_id"] = monsters[i]["monster_id"]
		monsterInfo["monster_uid"] = monsters[i]["uid"]
		monsterInfo["monster_name"] = monsters[i]["name"]
		monsterInfo["monster_rarity"] = monsters[i]["rarity"]
		monsterInfo["monster_element"] = monsters[i]["element"]
		monsterInfo["monster_hp"] = monsters[i]["hp"]
		monsterInfo["monster_attack"] = monsters[i]["attack"]
		monsterInfo["monster_defend"] = monsters[i]["defend"]
		monsterInfo["monster_speed"] = monsters[i]["speed"]
		monsterInfo["monster_lv"] = monsters[i]["lv"]
		monsterInfo["monster_energy"] = monsters[i]["energy"]
		if len(monsterskill) > 0 {
			monsterInfo["monster_skill"] = monsterskill[0]["skill"]
		} else {
			monsterInfo["monster_skill"] = ""
		}

		if err != nil {
			log.Println(err)
		}
		monsterDetail = append(monsterDetail, monsterInfo)
	}
	return monsterDetail, nil
}

func GetMonsterDetailByMonsterUid(monsteruid string) (map[string]string, error) {
	var skill string
	var skill_introduce string
	var introduce string
	monster, err := db.QueryString("SELECT * FROM `monster_info` WHERE uid =" + monsteruid)
	if err != nil {
		log.Println(err)
	}
	monsterskill, err := db.QueryString("SELECT * FROM `monster_skill` INNER JOIN skill ON monster_skill.skill = skill.skill WHERE monster_id = " + monsteruid)
	if err != nil {
		log.Println(err)
	}
	monster_introduce, err := db.QueryString("SELECT * FROM monster_data WHERE uid = " + monsteruid)
	if err != nil {
		log.Println(err)
	}
	if len(monsterskill) > 0 {
		skill = monsterskill[0]["skill"]
		skill_introduce = monsterskill[0]["introduce"]
	}
	if len(monster_introduce) > 0 {
		introduce = monster_introduce[0]["introduce"]
	}
	monsterDetail := map[string]string{
		"monster_id":              monster[0]["monster_id"],
		"monster_uid":             monster[0]["uid"],
		"monster_name":            monster[0]["name"],
		"monster_rarity":          monster[0]["rarity"],
		"monster_element":         monster[0]["element"],
		"monster_hp":              monster[0]["hp"],
		"monster_attack":          monster[0]["attack"],
		"monster_defend":          monster[0]["defend"],
		"monster_speed":           monster[0]["speed"],
		"monster_lv":              monster[0]["lv"],
		"monster_energy":          monster[0]["energy"],
		"monster_skill":           skill,
		"monster_skill_introduce": skill_introduce,
		"monster_desciption":      introduce,
	}
	return monsterDetail, nil
}

func GetMonsterExpAndLevel(monsteruid int) map[string]string {
	messages, err := db.QueryString("SELECT * FROM `monster_info` WHERE uid =" + strconv.Itoa(monsteruid) + "LIMIT 1")
	if err != nil {
		log.Println(err)
	}
	if len(messages) > 0 {
		payload := map[string]string{
			"lv":  messages[0]["lv"],
			"exp": messages[0]["exp"],
		}
		return payload
	} else {
		return nil
	}
}
