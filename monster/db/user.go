package db

import (
	"log"
	"math/rand"
	"strconv"
	"time"

	"monster/db/model"
)

func GetUser(uid int) (*model.User, error) {
	u := model.User{Uid: uid}
	has, err := db.Get(&u)
	if err != nil {
		logger.Println(err)
	}

	if !has {
		return nil, err
	}

	return &u, err
}

func CreateGuest(reqJSON map[string]string) (*model.User, error) {
	var randGuest int
	var guest string
	//产生种子码
	rand.Seed(time.Now().UnixNano())
	//随机产生ID号码
	randGuest = rand.Intn(8000000) + 1000000
	//产生ID用户名
	guest = "guest" + strconv.Itoa(randGuest)
	//产生用户初始资讯 (table: User)
	u := &model.User{
		Uid:         randGuest, //nft的账号id
		NickName:    guest,     //初始为nft账号id
		HeadIcon:    1,
		Frame:       1, // new
		GameCoin:    0,
		Strength:    0,
		Rank1:       1000, // new
		Rank2:       1000, // new
		EnergyLimit: 10,   // new
		CheckPoint:  1,
		Stage:       1,
	}
	//加入数据库里
	_, err := db.Insert(u)
	//检查err
	if err != nil {
		logger.Println("insert error:", err)
		//再来一次
		u, err = CreateGuest(reqJSON)
		return u, nil
		// CreateGuest(reqJSON)
		// return nil, err
	}

	return u, nil
}

func GetInfoById(uid string) (map[string]string, error) { // new
	var m31, m32, m33, m51, m52, m53, m54, m55 string
	var r31, r32, r33, r51, r52, r53, r54, r55 string
	user, err := db.QueryString("SELECT * FROM `user` WHERE uid = " + uid)
	monster3, err := db.QueryString("SELECT * FROM `user_embattle` inner join monster_info ON user_monster_uid = uid WHERE user_embattle.user_id = " + uid + " AND length = 3 AND current = 1")
	monster5, err := db.QueryString("SELECT * FROM `user_embattle` inner join monster_info ON user_monster_uid = uid WHERE user_embattle.user_id = " + uid + " AND length = 5 AND current = 1")
	if err != nil {
		log.Println(err)
	}

	m31 = ""
	m32 = ""
	m33 = ""
	m51 = ""
	m52 = ""
	m53 = ""
	m54 = ""
	m55 = ""
	r31 = ""
	r32 = ""
	r33 = ""
	r51 = ""
	r52 = ""
	r53 = ""
	r54 = ""
	r55 = ""


	if len(monster3) >= 3 {
		m31 = monster3[0]["user_monster_id"]
		r31 = monster3[0]["rarity"]
		m32 = monster3[1]["user_monster_id"]
		r31 = monster3[1]["rarity"]
		m33 = monster3[2]["user_monster_id"]
		r31 = monster3[2]["rarity"]
	}

	if len(monster5) >= 5 {
		m51 = monster5[0]["user_monster_id"]
		r51 = monster5[0]["rarity"]
		m52 = monster5[1]["user_monster_id"]
		r52 = monster5[1]["rarity"]
		m53 = monster5[2]["user_monster_id"]
		r53 = monster5[2]["rarity"]
		m54 = monster5[3]["user_monster_id"]
		r54 = monster5[3]["rarity"]
		m55 = monster5[4]["user_monster_id"]
		r55 = monster5[4]["rarity"]
	}

	userDetail := map[string]string{
		"uid":      user[0]["uid"],
		"nickname": user[0]["nick_name"],
		"headicon": user[0]["head_icon"],
		"frame":    user[0]["frame"],
		"rank1":    user[0]["rank1"],
		"rank2":    user[0]["rank2"],
		"m31":      m31, "r31": r31,
		"m32":      m32, "r32": r32,
		"m33":      m33, "r33": r33,
		"m51":      m51, "r51": r51,
		"m52":      m52, "r52": r52,
		"m53":      m53, "r53": r53,
		"m54":      m54, "r54": r54,
		"m55":      m55, "r55": r55,
	}

	return userDetail, nil
}

func GetHeadIconById(uid string) ([]map[string]string, error) {
	var icon []map[string]string
	selfIcon, err := db.QueryString("SELECT * FROM `user_item` WHERE user_id = '" + uid + "' AND type = 'icon'")
	if err != nil {
		log.Println(err)
	}
	for i := 0; i < len(selfIcon); i++ {
		iconInfo := make(map[string]string)
		iconInfo["item_id"] = selfIcon[i]["item_id"]
		iconInfo["item_name"] = selfIcon[i]["item_name"]
		iconInfo["type"] = selfIcon[i]["type"]
		iconInfo["rarity"] = selfIcon[i]["rarity"]
		iconInfo["introduce"] = selfIcon[i]["introduce"]
		iconInfo["used"] = selfIcon[i]["used"]
		icon = append(icon, iconInfo)
	}

	return icon, nil
}

func GetHeadFrameById(uid string) ([]map[string]string, error) {
	var frame []map[string]string
	selfFrame, err := db.QueryString("SELECT * FROM `user_item` WHERE user_id = '" + uid + "' AND type = 'frame'")
	if err != nil {
		log.Println(err)
	}
	for i := 0; i < len(selfFrame); i++ {
		frameInfo := make(map[string]string)
		frameInfo["item_id"] = selfFrame[i]["item_id"]
		frameInfo["item_name"] = selfFrame[i]["item_name"]
		frameInfo["type"] = selfFrame[i]["type"]
		frameInfo["rarity"] = selfFrame[i]["rarity"]
		frameInfo["introduce"] = selfFrame[i]["introduce"]
		frameInfo["used"] = selfFrame[i]["used"]
		frame = append(frame, frameInfo)
	}
	return frame, nil
}

func UpdateHead(uid, itemId, types string) {
	// remove all used icon
	_, err := db.Exec("Update `user_item` SET used = 0 WHERE user_id = " + uid + " AND type = '" + types + "'")
	// update new used icon
	_, err = db.Exec("Update `user_item` SET used = 1 WHERE user_id = " + uid + " AND type = '" + types + "' AND item_id = " + itemId)
	// update user head icon
	if types == "icon" {
		_, err = db.Exec("Update `user` SET head_icon = " + itemId + " WHERE uid = " + uid)
	} else {
		_, err = db.Exec("Update `user` SET frame = " + itemId + " WHERE uid = " + uid)
	}

	if err != nil {
		log.Println(err)
	}
}

func UpdateNickname(uid, nickName string){
	minLength := 3
	maxLength := 16
	if len(nickName) >= minLength && len(nickName) <= maxLength{
		_, err := db.Exec("UPDATE user SET nick_name = '" + nickName + "' WHERE uid = " + uid )
		if err != nil {
			log.Println(err)
		}
	}
}

func GetSeqByItem(uid, itemId, types string) int {
	item, err := db.QueryString("SELECT * FROM `user_item` WHERE user_id = '" + uid + "' AND type = '" + types + "'")
	if err != nil {
		log.Println(err)
	}
	if len(item) > 0 {
		var num int
		for i := 0; i < len(item); i++ {
			if item[i]["item_id"] == itemId {
				num = i
			}
		}
		return num
	}

	return -1
}
