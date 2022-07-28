package db

import (
	"errors"
	"log"

	"monster/helper"
)

func AcceptFriend(uid, friendid string) {
	_, err := db.Exec("DELETE FROM `request_friend_list` WHERE friend_id = " + uid + " AND user_id = " + friendid)
	_, err = db.Exec("INSERT INTO `user_friend_list` (`user_id`, `friend_id`) VALUES (?, ?)",
		uid, friendid)
	_, err = db.Exec("INSERT INTO `user_friend_list` (`user_id`, `friend_id`) VALUES (?, ?)",
		friendid, uid)
	if err != nil {
		log.Println(err)
	}
}

func RejectFriend(uid, friendid string) { // new
	_, err := db.Exec("DELETE FROM `request_friend_list` WHERE friend_id = " + uid + " AND user_id = " + friendid)
	if err != nil {
		log.Println(err)
	}
}

func AddFriend(uid, friendid string) int {
	friend, err := db.QueryString("SELECT * FROM `user_friend_list` WHERE user_id = " + uid + " AND friend_id = " + friendid)
	requested, err := db.QueryString("SELECT * FROM `request_friend_list` WHERE user_id = " + uid + " AND friend_id = " + friendid)
	requested1, err := db.QueryString("SELECT * FROM `request_friend_list` WHERE user_id = " + friendid + " AND friend_id = " + uid)
	if err != nil {
		log.Println(err)
	}
	if uid == friendid {
		return 4
	}
	if len(friend) > 0 {
		return 2
	} else if len(requested) > 0 || len(requested1) > 0 {
		return 1
	} else {
		_, err := db.Exec("INSERT INTO `request_friend_list` (`user_id`, `friend_id`) VALUES (?, ?)",
			uid, friendid)
		if err != nil {
			log.Println(err)
		}
		return 0
	}
}

func DeleteFriend(uid, friendid string) { // new
	_, err := db.Exec("DELETE FROM `user_friend_list` WHERE friend_id = " + uid + " AND user_id = " + friendid)
	_, err = db.Exec("DELETE FROM `user_friend_list` WHERE friend_id = " + friendid + " AND user_id = " + uid)
	if err != nil {
		log.Println(err)
	}
}

func GetFriendById(friendid string) ([]map[string]string, error) {
	friend, err := db.QueryString("SELECT * FROM `user` WHERE uid = " + friendid)
	if err != nil {
		log.Println(err)
	}

	return friend, nil
}

func GetFriendList(uid string) ([]map[string]string, error) {
	friendlist, err := db.QueryString("SELECT * FROM user_friend_list INNER JOIN user ON user.uid = user_friend_list.friend_id WHERE user_id = " + uid)
	if err != nil {
		log.Println(err)
	}
	// id, user_id, friend_id, id, uid, nick_name, head_icon, game_coin, strength, rank1, rank2, check_point, stage, frame, energy_limit, award, matching
	var friendslist []map[string]string
	if len(friendlist) > 0 {
		for i := 0; i < len(friendlist); i++ {
			friendInfo := make(map[string]string)
			friendInfo["uid"] = friendlist[i]["uid"]
			friendInfo["nickname"] = friendlist[i]["nick_name"]
			friendInfo["avatar"] = friendlist[i]["head_icon"]
			friendInfo["frame"] = friendlist[i]["frame"]
			friendslist = append(friendslist, friendInfo)
		}
	}
	/*
		// friend, err := db.QueryString("SELECT * FROM `user_friend_list` WHERE user_id = " + uid)
		// if err != nil {
		// 	log.Println(err)
		// }
		// var friendslist []map[string]string
		// if len(friend) > 0 {
		// 	for i := 0; i < len(friend); i++ {
		// 		friendlist, err := db.QueryString("SELECT * FROM `user` WHERE uid = " + friend[i]["friend_id"])
		// 		friendInfo := make(map[string]string)
		// 		friendInfo["uid"] = friendlist[0]["uid"]
		// 		friendInfo["nickname"] = friendlist[0]["nick_name"]
		// 		friendInfo["avatar"] = friendlist[0]["head_icon"]

		// 		if err != nil {
		// 			log.Println(err)
		// 		}
		// 		friendslist = append(friendslist, friendInfo)
		// 	}
		// }
	*/
	return friendslist, nil
}

func GetRequestFriend(uid string) ([]map[string]string, error) {
	requested, err := db.QueryString("SELECT * FROM `request_friend_list` WHERE friend_id = " + uid)
	if err != nil {
		log.Println(err)
	}
	var requestslist []map[string]string
	if len(requested) > 0 {
		for i := 0; i < len(requested); i++ {
			requestlist, err := db.QueryString("SELECT * FROM `user` WHERE uid = " + requested[i]["user_id"])
			requestInfo := make(map[string]string)
			requestInfo["uid"] = requestlist[0]["uid"]
			requestInfo["nickname"] = requestlist[0]["nick_name"]
			requestInfo["avatar"] = requestlist[0]["head_icon"]
			requestInfo["frame"] = requestlist[0]["frame"]

			if err != nil {
				log.Println(err)
			}
			requestslist = append(requestslist, requestInfo)
		}
	}
	return requestslist, nil
}

func GetMessage(uid, friendid string) ([]map[string]string, error) {
	messages, err := db.QueryString("SELECT * FROM `user_chat_list` WHERE sender_id = " + uid + " AND receiver_id = " + friendid + " OR sender_id = " + friendid + " AND receiver_id = " + uid)
	// 进入message后，把read变为1，表示已读 // new
	_, err = db.Exec("UPDATE `user_chat_list` SET user_chat_list.read = 1 WHERE sender_id = " + uid + " AND receiver_id = " + friendid + " OR sender_id = " + friendid + " AND receiver_id = " + uid)
	if err != nil {
		log.Println(err)
	}

	if len(messages) > 0 {
		var messagelist []map[string]string
		for i := 0; i < len(messages); i++ {
			senderInfo, err := db.QueryString("SELECT * FROM `user` WHERE uid = " + messages[i]["sender_id"])
			messageInfo := make(map[string]string)
			messageInfo["senderid"] = messages[i]["sender_id"]
			messageInfo["sendername"] = senderInfo[0]["nick_name"]
			messageInfo["message"] = messages[i]["msg"]
			messageInfo["avatar"] = senderInfo[0]["head_icon"]
			messageInfo["time"] = messages[i]["send_time"]

			if err != nil {
				log.Println(err)
			}
			messagelist = append(messagelist, messageInfo)
		}
		return messagelist, nil
	} else {
		return nil, errors.New("No Message")
	}
}

func SendMessage(senderid, receiverid, message, time string) {
	_, err := db.Exec("INSERT INTO `user_chat_list` (`sender_id`, `receiver_id`, `msg`, `send_time`) VALUES (?, ?, ?, ?)",
		senderid, receiverid, message, time)
	if err != nil {
		log.Println(err)
	}
}

func GetMessageById(uid string) ([]map[string]string, error) { // new
	friends, err := db.QueryString("SELECT * FROM user_friend_list INNER JOIN user ON user.uid = user_friend_list.friend_id WHERE user_id = " + uid)
	if err != nil {
		log.Println(err)
	}
	// 	SELECT * FROM
	// 	(SELECT * FROM monster.user_chat_list
	// 	WHERE send_time IN (
	// 		SELECT MAX(send_time) FROM monster.user_chat_list WHERE receiver_id = 4643485 AND sender_id = 5457837 OR sender_id = 4643485 AND receiver_id = 5457837
	// 		GROUP BY CONCAT(IF (sender_id > receiver_id, sender_id, receiver_id),IF (sender_id < receiver_id, sender_id, receiver_id))
	//         )
	// 	AND (receiver_id = 4643485 AND sender_id = 5457837 OR sender_id = 4643485 AND receiver_id = 5457837)
	// ORDER BY id DESC) c ORDER BY c.send_time DESC

	var messages []map[string]string
	for y := 0; y < len(friends); y++ {
		message, err := db.QueryString(
			"SELECT * FROM (SELECT * FROM `user_chat_list` WHERE send_time IN (" +
				" SELECT MAX(send_time) FROM `user_chat_list` WHERE receiver_id = " + uid + " AND sender_id = " + friends[y]["uid"] + " OR sender_id = " + uid + " AND receiver_id = " + friends[y]["uid"] +
				" GROUP BY CONCAT(IF (sender_id > receiver_id, sender_id, receiver_id)," +
				" IF (sender_id < receiver_id, sender_id, receiver_id)))" +
				" AND (receiver_id = " + uid + " AND sender_id = " + friends[y]["uid"] + " OR sender_id = " + uid + " AND receiver_id = " + friends[y]["uid"] +
				" ) ORDER BY id DESC) c ORDER BY c.send_time DESC")
		if err != nil {
			log.Println(err)
		}
		for z := 0; z < len(message); z++ {
			messages = append(messages, message[z])
		}
	}

	var messageslist []map[string]string
	var friend string
	var friendList []string
	if len(messages) > 0 {
		for i := 0; i < len(messages); i++ {
			if messages[i]["sender_id"] != uid {
				friend = messages[i]["sender_id"]
			} else if messages[i]["receiver_id"] != uid {
				friend = messages[i]["receiver_id"]
			}
			var inFriendList bool = false
			for _, v := range friendList {
				if v == friend {
					inFriendList = true
				}
			}
			if inFriendList {
				continue
			} else {
				friendList = append(friendList, friend)
			}
			messagelist, err := db.QueryString("SELECT * FROM `user` WHERE uid = " + friend)
			messageInfo := make(map[string]string)
			messageInfo["uid"] = messagelist[0]["uid"]
			messageInfo["nickname"] = messagelist[0]["nick_name"]
			messageInfo["avatar"] = messagelist[0]["head_icon"]
			messageInfo["frame"] = messagelist[0]["frame"]
			messageInfo["read"] = messages[i]["read"]
			messageInfo["sendtime"] = messages[i]["send_time"]
			if err != nil {
				log.Println(err)
			}
			messageslist = append(messageslist, messageInfo)
		}
	}

	return messageslist, nil
}

func CheckEnergyLimitPerDay(uid string) int { // new
	user, err := db.QueryString("SELECT * FROM `user` WHERE uid = " + uid)
	if err != nil {
		log.Println(err)
	}
	energy := helper.StringToInt(user[0]["energy_limit"])
	var amount int
	if energy > 0 { // still have gift
		amount = 0
	} else { // limited
		amount = 1
	}

	return amount
}

func EnergyGiven(uid, friendid string) { // new
	// give energy
	_, err := db.Exec("Update `user_prop` SET amount = amount + 1 WHERE prop_id = 1001 AND user_id = " + friendid)
	if err != nil {
		log.Println(err)
	}
	// reduce number of gift
	_, err = db.Exec("Update `user` SET energy_limit = energy_limit - 1 WHERE uid = " + friendid)
	if err != nil {
		log.Println(err)
	}
}

func RenewEnergyLimit() { // new
	_, err := db.Exec("Update `user` SET energy_limit = 10")
	if err != nil {
		log.Println(err)
	}
}
