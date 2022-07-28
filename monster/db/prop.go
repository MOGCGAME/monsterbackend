package db

import (
	"log"
)

func GetProp(uid string) ([]map[string]string, error) {
	var prop []map[string]string
	messages, err := db.QueryString("SELECT * FROM `user_prop` WHERE user_id = " + uid)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	for i := 0; i < len(messages); i++ {
		props, err := db.QueryString("SELECT * FROM `prop_info` WHERE id = " + messages[i]["prop_id"])
		propInfo := make(map[string]string)
		propInfo["id"] = props[0]["id"]
		propInfo["name"] = props[0]["name"]
		propInfo["rarity"] = props[0]["rarity"]
		propInfo["classify"] = props[0]["classify"]
		propInfo["introduce"] = props[0]["introduce"]
		propInfo["amount"] = messages[i]["amount"]
		if err != nil {
			log.Println(err)
		}
		prop = append(prop, propInfo)
	}
	return prop, nil
}
