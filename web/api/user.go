package api

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/lonng/nex"
	log "github.com/sirupsen/logrus"
	"monster/db"
	"monster/db/model"
	"monster/helper"
)

var (
	logger = log.WithFields(log.Fields{"component": "http", "service": "login"})
)

func MakeUserService() http.Handler { //用户
	router := mux.NewRouter()

	router.Handle("/user/getGuest", nex.Handler(getGuestHandler)).Methods("POST")
	router.Handle("/user/getInfo", nex.Handler(getInfoHandler)).Methods("POST")
	router.Handle("/user/getHeadIcon", nex.Handler(getHeadIconHandler)).Methods("POST")
	router.Handle("/user/getHeadFrame", nex.Handler(getHeadFrameHandler)).Methods("POST")
	router.Handle("/user/updateHead", nex.Handler(updateHeadHandler)).Methods("POST")
	router.Handle("/user/updateNickname", nex.Handler(updateNicknameHandler)).Methods("POST")
	return router
}

func updateHeadHandler(r *http.Request) (map[string]string, error) {
	uid, isValid := helper.VerifyJWT(r)
	if !isValid {
		return nil, errors.New("Invalid token")
	}
	reqJSON := helper.ReadParameters(r)
	db.UpdateHead(uid, reqJSON["itemId"], reqJSON["type"])
	seq := db.GetSeqByItem(uid, reqJSON["itemId"], reqJSON["type"])
	payload := map[string]string{
		"code": "success",
		"seq":  strconv.Itoa(seq),
	}

	return payload, nil
}

func updateNicknameHandler(r *http.Request) (map[string]string, error) {
	uid, isValid := helper.VerifyJWT(r)
	if !isValid {
		return nil, errors.New("Invalid token")
	}
	reqJSON := helper.ReadParameters(r)
	minLength := 3
	maxLength := 16
	var code string
	if len(reqJSON["nickname"]) >= minLength && len(reqJSON["nickname"]) <= maxLength{
		db.UpdateNickname(uid, reqJSON["nickname"])
		code = "success"
	}else{
		code = "昵称不符合长度，需在"+strconv.Itoa(minLength)+"和"+strconv.Itoa(maxLength)+"的字符范围内"
	}

	payload := map[string]string{
		"code": code,
		"nickname": reqJSON["nickname"],
	}

	return payload, nil
}


func getHeadFrameHandler(r *http.Request) (map[string]interface{}, error) { // new
	uid, isValid := helper.VerifyJWT(r)
	if !isValid {
		return nil, errors.New("Invalid token")
	}
	frame, err := db.GetHeadFrameById(uid)
	if err != nil {
		log.Println(err)
	}
	payload := map[string]interface{}{
		"frame": frame,
	}

	return payload, nil
}

func getHeadIconHandler(r *http.Request) (map[string]interface{}, error) { // new
	uid, isValid := helper.VerifyJWT(r)
	if !isValid {
		return nil, errors.New("Invalid token")
	}
	icon, err := db.GetHeadIconById(uid)
	if err != nil {
		log.Println(err)
	}
	payload := map[string]interface{}{
		"icon": icon,
	}

	return payload, nil
}

func getInfoHandler(r *http.Request) (map[string]interface{}, error) { // new
	reqJSON := helper.ReadParameters(r)
	info, err := db.GetInfoById(reqJSON["uid"])
	if err != nil {
		log.Println(err)
	}
	payload := map[string]interface{}{
		"info": info,
	}

	return payload, nil
}

func getGuestById(uid string) (*model.User, map[string]interface{}) {
	user, err := db.GetUser(helper.StringToInt(uid))
	if err != nil || user == nil {
		log.Println(err)
		return nil, nil
	}

	payload := buildUserPayload(user)

	return user, payload
}

//前端 提取游客信息 POST请求（暂时使用）
func getGuestHandler(r *http.Request) (map[string]interface{}, error) {
	reqJSON := helper.ReadParameters(r)
	//验证JWT Tokem
	uid, isValid := helper.VerifyJWT(r)
	//检查验证结果有问题
	if !isValid {
		//创造新用户ID
		return createGuestHandler(r, reqJSON)
	}
	//获取用户资讯
	user, payload := getGuestById(uid)
	//检查用户资讯为空
	if user == nil {
		//创造新用户ID
		return createGuestHandler(r, reqJSON)
	}
	return payload, nil
}

//前端 创建游客 POST请求（暂时使用）
func createGuestHandler(r *http.Request, reqJSON map[string]string) (map[string]interface{}, error) {
	//创建游客用户
	user, err := db.CreateGuest(reqJSON)
	//检查err
	if err != nil {
		payload := map[string]interface{}{
			"error": err.Error(),
		}
		return payload, err
	}
	//获取游客用户的User UID
	uid := strconv.Itoa(user.Uid)
	//重新建立User payload
	payload := buildUserPayload(user)
	///产生新的加密token
	jwtString := helper.NewJWT(uid, payload)
	//把加密token加进payload
	payload["jwt"] = jwtString

	return payload, nil
}

func buildUserPayload(user *model.User) map[string]interface{} {
	payload := map[string]interface{}{
		"uid":      user.Uid,
		"nickName": user.NickName,
		"headIcon": user.HeadIcon,
		"frame":    user.Frame,
		"gameCoin": user.GameCoin,
		"strength": user.Strength,
		"rank1":    user.Rank1,
		"rank2":    user.Rank2,
	}
	return payload
}
