package api

import (
	"errors"
	"net/http"

	"monster/db"
	"monster/helper"

	"github.com/gorilla/mux"
	"github.com/lonng/nex"
	log "github.com/sirupsen/logrus"
)

func MakeFriendService() http.Handler { //背包
	router := mux.NewRouter()
	router.Handle("/friend/searchFriend", nex.Handler(searchFriendHandler)).Methods("POST")
	router.Handle("/friend/getFriendList", nex.Handler(getFriendListHandler)).Methods("POST")
	router.Handle("/friend/addFriend", nex.Handler(addFriendHandler)).Methods("POST")
	router.Handle("/friend/deleteFriend", nex.Handler(deleteFriendHandler)).Methods("POST")
	router.Handle("/friend/getFriendRequest", nex.Handler(getFriendRequestHandler)).Methods("POST")
	router.Handle("/friend/acceptFriend", nex.Handler(acceptFriendHandler)).Methods("POST")
	router.Handle("/friend/rejectFriend", nex.Handler(rejectFriendHandler)).Methods("POST")
	router.Handle("/friend/getMessage", nex.Handler(getMessageHandler)).Methods("POST")
	router.Handle("/friend/getMessageList", nex.Handler(getMessageListHandler)).Methods("POST")
	router.Handle("/friend/sendMessage", nex.Handler(sendMessageHandler)).Methods("POST")
	return router
}

func searchFriendHandler(r *http.Request) (map[string]interface{}, error) {
	reqJSON := helper.ReadParameters(r)
	strangerList, err := db.GetFriendById(reqJSON["friendid"])
	if err != nil {
		log.Println(err)
	}
	payload := map[string]interface{}{
		"strangerList": strangerList,
	}

	return payload, nil
}

func getFriendListHandler(r *http.Request) (map[string]interface{}, error) {
	uid, isValid := helper.VerifyJWT(r)
	if !isValid {
		return nil, errors.New("Invalid token")
	}
	// friend list limit 60
	friendList, err := db.GetFriendList(uid)
	if err != nil {
		log.Println(err)
	}
	payload := map[string]interface{}{
		"friendList": friendList,
	}

	return payload, nil
}

func addFriendHandler(r *http.Request) (map[string]interface{}, error) {
	uid, isValid := helper.VerifyJWT(r)
	if !isValid {
		return nil, errors.New("Invalid token")
	}
	reqJSON := helper.ReadParameters(r)
	request := db.AddFriend(uid, reqJSON["friendid"])
	var code string
	if request == 1 {
		code = "1"
	} else if request == 2 {
		code = "2"
	} else if request == 4 {
		code = "4"
	} else {
		code = "3"
	}
	payload := map[string]interface{}{
		"code": code,
	}

	return payload, nil
}

func deleteFriendHandler(r *http.Request) (map[string]interface{}, error) { // new
	uid, isValid := helper.VerifyJWT(r)
	if !isValid {
		return nil, errors.New("Invalid token")
	}
	reqJSON := helper.ReadParameters(r)
	db.DeleteFriend(uid, reqJSON["friendid"])
	payload := map[string]interface{}{
		"code": "success",
	}

	return payload, nil
}

func getFriendRequestHandler(r *http.Request) (map[string]interface{}, error) {
	uid, isValid := helper.VerifyJWT(r)
	if !isValid {
		return nil, errors.New("Invalid token")
	}
	requestList, err := db.GetRequestFriend(uid)
	if err != nil {
		log.Println(err)
	}
	payload := map[string]interface{}{
		"requestList": requestList,
	}

	return payload, nil
}

func acceptFriendHandler(r *http.Request) (map[string]interface{}, error) {
	uid, isValid := helper.VerifyJWT(r)
	if !isValid {
		return nil, errors.New("Invalid token")
	}
	reqJSON := helper.ReadParameters(r)
	db.AcceptFriend(uid, reqJSON["friendid"])
	payload := map[string]interface{}{
		"code": "success",
	}

	return payload, nil
}

func rejectFriendHandler(r *http.Request) (map[string]interface{}, error) { // new
	uid, isValid := helper.VerifyJWT(r)
	if !isValid {
		return nil, errors.New("Invalid token")
	}
	reqJSON := helper.ReadParameters(r)
	db.RejectFriend(uid, reqJSON["friendid"])
	payload := map[string]interface{}{
		"code": "success",
	}

	return payload, nil
}

func energyGivenHandler(r *http.Request) (map[string]interface{}, error) { // new
	uid, isValid := helper.VerifyJWT(r)
	if !isValid {
		return nil, errors.New("Invalid token")
	}
	var code string
	reqJSON := helper.ReadParameters(r)
	checkLimit := db.CheckEnergyLimitPerDay(uid)
	if checkLimit == 0 {
		db.EnergyGiven(uid, reqJSON["friendid"])
		code = "success"
	} else {
		code = "limit"
	}

	payload := map[string]interface{}{
		"code": code,
	}

	return payload, nil
}

func getMessageHandler(r *http.Request) (map[string]interface{}, error) {
	uid, isValid := helper.VerifyJWT(r)
	if !isValid {
		return nil, errors.New("Invalid token")
	}
	reqJSON := helper.ReadParameters(r)
	message, err := db.GetMessage(uid, reqJSON["friendid"])
	if err != nil {
		log.Println(err)
	}
	payload := map[string]interface{}{
		"message": message,
	}

	return payload, nil
}

func getMessageListHandler(r *http.Request) (map[string]interface{}, error) { // new
	uid, isValid := helper.VerifyJWT(r)
	if !isValid {
		return nil, errors.New("Invalid token")
	}
	messageList, err := db.GetMessageById(uid)
	if err != nil {
		log.Println(err)
	}
	payload := map[string]interface{}{
		"messageList": messageList,
	}

	return payload, nil
}

func sendMessageHandler(r *http.Request) (map[string]interface{}, error) {
	uid, isValid := helper.VerifyJWT(r)
	if !isValid {
		return nil, errors.New("Invalid token")
	}
	reqJSON := helper.ReadParameters(r)
	db.SendMessage(reqJSON["senderid"], reqJSON["receiverid"], reqJSON["message"], reqJSON["time"])
	message, err := db.GetMessage(uid, reqJSON["receiverid"])
	if err != nil {
		log.Println(err)
	}
	payload := map[string]interface{}{
		"message": message,
	}

	return payload, nil
}
