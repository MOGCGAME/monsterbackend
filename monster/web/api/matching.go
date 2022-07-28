package api

import (
	"errors"
	"math"
	"math/rand"
	"monster/db"
	"monster/helper"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/lonng/nex"
)

func MakeMatchingService() http.Handler { //布阵
	router := mux.NewRouter()
	router.Handle("/matching/setRanking", nex.Handler(setRankingHandler)).Methods("POST")
	router.Handle("/matching/getRanking", nex.Handler(getRankingHandler)).Methods("POST")
	router.Handle("/matching/getPvERecord", nex.Handler(getPvERecordHandler)).Methods("POST")
	router.Handle("/matching/getPvPMatching", nex.Handler(getPvPMatchingHandler)).Methods("POST")
	router.Handle("/matching/getPvEMatching", nex.Handler(getPvEMatchingHandler)).Methods("POST")
	router.Handle("/matching/getSelfInfo", nex.Handler(getSelfInfoHandler)).Methods("POST")
	router.Handle("/matching/renewMatching", nex.Handler(renewMatchingHandler)).Methods("POST")
	router.Handle("/matching/getCheckPoint", nex.Handler(getCheckPoint)).Methods("POST")
	router.Handle("/matching/getStage", nex.Handler(getStage)).Methods("POST")
	router.Handle("/matching/updatePVPResult", nex.Handler(updatePVPResultHandler)).Methods("POST")
	router.Handle("/matching/updatePVEResult", nex.Handler(updatePVEResultHandler)).Methods("POST")
	return router
}

func updatePVEResultHandler(r *http.Request) (int, error) {
	uid, isValid := helper.VerifyJWT(r)
	if !isValid {
		return 0, errors.New("Invalid token")
	}
	reqJSON := helper.ReadParameters(r)

	if reqJSON["result"] == "0" {
		currStage := db.GetCurrentStage(uid)

		if helper.StringToInt(reqJSON["stage"]) >= currStage {
			if helper.StringToInt(reqJSON["stage"]) == 10 {
				db.UpdateCheckPoint(uid)
				db.UpdateStage(uid, 1) // 更新成下一关的一阶段
			} else {
				db.UpdateStage(uid, 0)
			}
		}
	}

	db.UpdateAward(reqJSON["award"], uid)

	return 0, nil
}

func updatePVPResultHandler(r *http.Request) (int, error) {
	uid, isValid := helper.VerifyJWT(r)
	if !isValid {
		return 0, errors.New("Invalid token")
	}
	reqJSON := helper.ReadParameters(r)

	if reqJSON["length"] == "5" { // new
		db.UpdateRank2(reqJSON["rank"], uid)
	} else {
		db.UpdateRank1(reqJSON["rank"], uid)
	}

	db.UpdateAward(reqJSON["award"], uid)

	return 0, nil
}

func getCheckPoint(r *http.Request) (int, error) {
	uid, isValid := helper.VerifyJWT(r)
	if !isValid {
		return 0, errors.New("Invalid token")
	}
	reqJSON := helper.ReadParameters(r)
	checkpoint, err := db.GetPlayerCheckPoint(uid, helper.StringToInt(reqJSON["checkpoint"]))
	if err != nil {
		return 0, errors.New("Invalid token")
	}
	return checkpoint, nil
}

func getStage(r *http.Request) (map[string]interface{}, error) {
	uid, isValid := helper.VerifyJWT(r)
	if !isValid {
		return nil, errors.New("Invalid token")
	}
	reqJSON := helper.ReadParameters(r)
	checkpoint, err := db.GetPlayerCheckPoint(uid, helper.StringToInt(reqJSON["checkpoint"]))
	if err != nil {
		return nil, errors.New("Invalid token")
	}
	stage, err := db.GetPlayerStage(uid, helper.StringToInt(reqJSON["stage"]))
	if err != nil {
		return nil, errors.New("Invalid token")
	}
	length, err := db.GetStageLen(reqJSON["checkpoint"], reqJSON["stage"])
	payload := map[string]interface{}{
		"checkpoint": checkpoint,
		"stage":      stage,
		"length":     length,
	}
	return payload, nil
}

func renewMatchingHandler(r *http.Request) (map[string]string, error) {
	uid, isValid := helper.VerifyJWT(r)
	if !isValid {
		return nil, errors.New("Invalid token")
	}
	reqJSON := helper.ReadParameters(r)

	db.RenewMatching(uid, reqJSON["matching"])

	payload := map[string]string{
		"code": "success",
	}

	return payload, nil
}

func getSelfInfoHandler(r *http.Request) (map[string]interface{}, error) {
	uid, isValid := helper.VerifyJWT(r)
	if !isValid {
		return nil, errors.New("Invalid token")
	}
	reqJSON := helper.ReadParameters(r)

	selfEmbattle, err := db.GetSelfEmbattle(uid, reqJSON["length"])
	if err != nil {
		return nil, errors.New("Monster Not Enough Energy")
	}

	payload := map[string]interface{}{
		"mode": "pvp",
		"self": selfEmbattle,
	}

	return payload, nil
}

func getPvERecordHandler(r *http.Request) (map[string]interface{}, error) {
	payload := map[string]interface{}{}

	return payload, nil
}

func getPvPMatchingHandler(r *http.Request) (map[string]interface{}, error) {
	uid, isValid := helper.VerifyJWT(r)
	if !isValid {
		return nil, errors.New("Invalid token")
	}
	reqJSON := helper.ReadParameters(r)

	db.UpdateMatching(uid, reqJSON["matching"])

	selfEmbattle, err := db.GetSelfEmbattle(uid, reqJSON["length"])
	if err != nil {
		payload := map[string]interface{}{"code": "monster not enough energy"}
		return payload, nil
	}
	rank, err := db.GetRanking(uid, reqJSON["length"])
	if err != nil {
		payload := map[string]interface{}{"code": "wrong ranking"}
		return payload, nil
	}
	rank1 := rank - 100
	rank2 := rank + 100
	enemyEmbattle := db.GetPvPEnemyEmbattle(uid, reqJSON["length"], strconv.Itoa(rank1), strconv.Itoa(rank2), reqJSON["matching"])
	if enemyEmbattle == nil {
		payload := map[string]interface{}{"code": "waiting opponent"}
		return payload, nil
	} else {
		var enemyId string = enemyEmbattle[0]["user_id"]
		// for k, _ := range enemyEmbattle {
		// 	if k == 0 {
		// 		enemyId = enemyEmbattle[k]["user_id"]
		// 	}
		// }

		var selfEmbattleLen, enemyEmbattleLen int
		selfEmbattleLen = len(selfEmbattle)
		enemyEmbattleLen = len(enemyEmbattle)

		matchInfo, _ := db.GetPvPEnemy(uid, enemyId, reqJSON["length"])
		matchInfo1, err := db.GetPvPEnemy(uid, enemyId, reqJSON["length"])
		if err != nil {
			payload := map[string]interface{}{"code": "match info error"}
			return payload, nil
		}
		initBattle = matchInfo1
		battleInfo = matchInfo
		buffMatchInfo, buff := getBuff(initBattle, helper.StringToInt(uid))
		battleInfo = buffMatchInfo
		//update selfEmbattle and enemyEmbattle
		for b := 0; b < len(battleInfo); b++ {
			if battleInfo[b]["user_id"] == uid {
				for s := 0; s < len(selfEmbattle); s++ {
					if selfEmbattle[s]["monster_uid"] == battleInfo[b]["monster_uid"] {
						selfEmbattle[s]["monster_hp"] = battleInfo[b]["monster_hp"]
						selfEmbattle[s]["monster_attack"] = battleInfo[b]["monster_attack"]
						selfEmbattle[s]["monster_defend"] = battleInfo[b]["monster_defend"]
						selfEmbattle[s]["monster_max_hp"] = battleInfo[b]["monster_max_hp"]
						selfEmbattle[s]["monster_speed"] = battleInfo[b]["monster_speed"]
						selfEmbattle[s]["monster_miss"] = battleInfo[b]["monster_miss"]
						selfEmbattle[s]["monster_hit"] = battleInfo[b]["monster_hit"]
						selfEmbattle[s]["monster_skill_rate"] = battleInfo[b]["monster_skill_rate"]
						selfEmbattle[s]["monster_positive"] = battleInfo[b]["monster_positive"]
						selfEmbattle[s]["monster_negative"] = battleInfo[b]["monster_negative"]
						break
					}
				}
			} else {
				for e := 0; e < len(enemyEmbattle); e++ {
					if enemyEmbattle[e]["monster_uid"] == battleInfo[b]["monster_uid"] {
						enemyEmbattle[e]["monster_hp"] = battleInfo[b]["monster_hp"]
						enemyEmbattle[e]["monster_attack"] = battleInfo[b]["monster_attack"]
						enemyEmbattle[e]["monster_defend"] = battleInfo[b]["monster_defend"]
						enemyEmbattle[e]["monster_max_hp"] = battleInfo[b]["monster_max_hp"]
						enemyEmbattle[e]["monster_speed"] = battleInfo[b]["monster_speed"]
						enemyEmbattle[e]["monster_miss"] = battleInfo[b]["monster_miss"]
						enemyEmbattle[e]["monster_hit"] = battleInfo[b]["monster_hit"]
						enemyEmbattle[e]["monster_skill_rate"] = battleInfo[b]["monster_skill_rate"]
						enemyEmbattle[e]["monster_positive"] = battleInfo[b]["monster_positive"]
						enemyEmbattle[e]["monster_negative"] = battleInfo[b]["monster_negative"]
						break
					}
				}
			}
		}
		match := getFullMatch(helper.StringToInt(uid), selfEmbattleLen, enemyEmbattleLen)
		ranking := "20" //排行榜积分
		exp := "160"    //经验值
		payload := map[string]interface{}{
			"code":          "success",
			"mode":          "pvp",
			"self":          selfEmbattle,
			"enemy":         enemyEmbattle,
			"buff":          buff,
			"matchInfo":     initBattle,
			"speed":         speedCount,
			"match":         match,
			"rank":          ranking,
			"exp":           exp,
			"room":          reqJSON["length"],
			"buffMatchInfo": buffMatchInfo,
		}

		//更新数据库，双方已退出PVP
		if payload["code"] == "success" {
			db.RenewMatching(uid, "0")
			db.RenewMatching(enemyId, "0")
		}

		return payload, nil
	}
}

func getPvEMatchingHandler(r *http.Request) (map[string]interface{}, error) {
	uid, isValid := helper.VerifyJWT(r)
	if !isValid {
		return nil, errors.New("Invalid token")
	}
	reqJSON := helper.ReadParameters(r)
	selfEmbattle, err := db.GetSelfEmbattle(uid, reqJSON["length"])
	if err != nil {
		return nil, errors.New("Monster Not Enough Energy")
	}
	enemyEmbattle, err := db.GetPvEEnemyEmbattle(reqJSON["checkpoint"], reqJSON["stage"])
	if err != nil {
		payload := map[string]interface{}{"error": "Invalid token"}
		return payload, errors.New("Invalid token")
	}

	award, err1 := db.GetAwardByStage(reqJSON["checkpoint"], reqJSON["stage"])
	exp, err1 := db.GetExpByStage(reqJSON["checkpoint"], reqJSON["stage"])
	if err1 != nil {
		payload := map[string]interface{}{"error": "Invalid token"}
		return payload, errors.New("Invalid token")
	}

	var enemyId string
	for k, _ := range enemyEmbattle {
		if k == 0 {
			enemyId = enemyEmbattle[k]["user_id"]
		}
	}

	matchInfo, err := db.GetPvEEnemy(uid, enemyId, reqJSON["length"])
	if err != nil {
		payload := map[string]interface{}{"error": "Invalid token"}
		return payload, errors.New("Invalid token")
	}

	prior := getPriorHand()
	buffMatchInfo, buff := getBuff(matchInfo, helper.StringToInt(uid))
	selfEmbattle, enemyEmbattle = giveStartBuff(uid, buffMatchInfo, selfEmbattle, enemyEmbattle)
	payload := map[string]interface{}{
		"mode":       "pve",
		"self":       selfEmbattle,
		"enemy":      enemyEmbattle,
		"buff":       buff,
		"matchInfo":  buffMatchInfo,
		"initInfo":   matchInfo,
		"prior":      prior,
		"checkpoint": reqJSON["checkpoint"],
		"stage":      reqJSON["stage"],
		"award":      award,
		"exp":        exp,
	}

	return payload, nil
}

var battleInfo []map[string]string
var initBattle []map[string]string
var speedCount []map[string]string

type (
	BuffInfo struct {
		Self  []int
		Enemy []int
	}

	RoundInfo struct {
		Round     int
		PriorHand int
		SelfLeft  int
		EnemyLeft int
		Count     int
		TurnInfo  []TurnInfo
	}

	TurnInfo struct {
		AtkMonsterUid         int                 `json:"atkMonsterUid"`
		AtkMonsterId          int                 `json:"atkMonsterId"`
		AtkMonsterSide        bool                `json:"atkMonsterSide"`
		AtkMonsterSeq         int                 `json:"atkMonsterSeq"`
		AtkMonsterHp          int                 `json:"atkMonsterHp"`
		AtkMonsterDamage      int                 `json:"atkMonsterDamage"`
		AtkMonsterElement     int                 `json:"atkMonsterElement"`
		AtkMonsterSkill       int                 `json:"atkMonsterSkill"`
		AtkMonsterHeal        int                 `json:"atkMonsterHeal"`
		AtkMonsterNegative    string              `json:"atkMonsterNegative"`
		AtkMonsterPositive    string              `json:"atkMonsterPositive"`
		AtkMonsterTrigger     int                 `json:"atkMonsterTrigger"`
		AtkMonsterStun        bool                `json:"atkMonsterStun"`
		AtkMultiple           []map[string]string `json:"atkmultiple"`
		AtkMonsterTriggered   bool                `json:"atkMonsterTriggered"`
		AtkMonsterMiss        bool                `json:"atkMonsterMiss"`
		TargetMonsterId       int                 `json:"targetMonsterId"`
		TargetMonsterSeq      int                 `json:"targetMonsterSeq"`
		TargetMonsterHp       int                 `json:"targetMonsterHp"`
		TargetMonsterAttack   int                 `json:"targetMonsterAttack"`
		TargetMonsterDefense  int                 `json:"targetMonsterDefense"`
		TargetMonsterSpeed    int                 `json:"targetMonsterSpeed"`
		TargetMonsterElement  int                 `json:"targetMonsterElement"`
		TargetMonsterNegative string              `json:"targetMonsterNegative"`
	}
	// TurnInfo struct {
	// 	AtkMonsterUid         int                 `json:"atkMonsterUid"`
	// 	AtkMonsterId          int                 `json:"atkMonsterId"`
	// 	AtkMonsterSide        int               `json:"atkMonsterSide"`
	// 	AtkMonsterSeq         int                 `json:"atkMonsterSeq"`
	// 	AtkMonsterHp          int                 `json:"atkMonsterHp"`
	// 	AtkMonsterDamage      int                 `json:"atkMonsterDamage"`
	// 	AtkMonsterElement     int                 `json:"atkMonsterElement"`
	// 	AtkMonsterSkill       int                 `json:"atkMonsterSkill"`
	// 	AtkMonsterHeal        float64                `json:"atkMonsterHeal"`
	// 	AtkMonsterNegative    string              `json:"atkMonsterNegative"`
	// 	AtkMonsterPositive    string              `json:"atkMonsterPositive"`
	// 	AtkMonsterTrigger     int                 `json:"atkMonsterTrigger"`
	// 	AtkMultiple           []map[string]string `json:"atkmultiple"`
	// 	AtkMonsterTriggered   int                 `json:"atkMonsterTriggered"`
	// 	AtkMonsterMiss        int                 `json:"atkMonsterMiss"`
	// 	TargetMonsterId       int                 `json:"targetMonsterId"`
	// 	TargetMonsterSeq      int                 `json:"targetMonsterSeq"`
	// 	TargetMonsterHp       int                 `json:"targetMonsterHp"`
	// 	TargetMonsterElement  int                 `json:"targetMonsterElement"`
	// 	TargetMonsterNegative string              `json:"targetMonsterNegative"`
	// }
)

//开局状态
func getBuff(matchInfo []map[string]string, selfUid int) ([]map[string]string, BuffInfo) {
	var selfpositive, selfnegative, enemypositive, enemynegative, self, enemy []int
	var skill, a, b int
	var c float64
	var selfIndex, enemyIndex []int
	for a := 0; a < len(matchInfo); a++ {
		if helper.StringToInt(matchInfo[a]["user_id"]) == selfUid {
			selfIndex = append(selfIndex, a)
		} else {
			enemyIndex = append(enemyIndex, a)
		}
	}
	//根据每个matchInfo
	for i := 0; i < len(matchInfo); i++ {
		//从matchInfo获取怪兽技能ID
		skill = helper.StringToInt(matchInfo[i]["monster_skill"])
		//根据怪兽技能ID
		switch skill {
		//怪兽我方全体上升类型 （正面buff）
		case 9001, 9002, 9003, 9004, 9005:
			//检查怪兽的user id 是我方user id
			if helper.StringToInt(matchInfo[i]["user_id"]) == selfUid {
				//把技能ID加进selfpositive和self里
				selfpositive = append(selfpositive, skill)
				self = append(self, skill)
			} else {
				//把技能ID加进enemypositive和enemy里
				enemypositive = append(enemypositive, skill)
				enemy = append(enemy, skill)
			}
		//怪兽敌方全体下降类型 （负面buff）
		case 9011, 9012, 9013, 9014, 9015:
			//检查怪兽的user id 是我方user id
			if helper.StringToInt(matchInfo[i]["user_id"]) == selfUid {
				//把技能ID加进enemynegative和enemy里
				enemynegative = append(enemynegative, skill)
				enemy = append(enemy, skill)
			} else {
				//把技能ID加进selfnegative和self里
				selfnegative = append(selfnegative, skill)
				self = append(self, skill)
			}
		case 9100:
			if helper.StringToInt(matchInfo[i]["user_id"]) == selfUid {
				var randomIndex = rand.Intn(len(enemyIndex))
				var matchIndex = enemyIndex[randomIndex]
				matchInfo[matchIndex]["monster_skill_rate"] = strconv.Itoa(helper.StringToInt(matchInfo[matchIndex]["monster_skill_rate"]) - 5)
			} else {
				var randomIndex = rand.Intn(len(selfIndex))
				var matchIndex = selfIndex[randomIndex]
				matchInfo[matchIndex]["monster_skill_rate"] = strconv.Itoa(helper.StringToInt(matchInfo[matchIndex]["monster_skill_rate"]) - 5)
			}
		}
	}
	//根据每个我方的正面buff
	for i := 0; i < len(selfpositive); i++ {
		//根据每个matchInfo(每只怪兽)
		for j := 0; j < len(matchInfo); j++ {
			//检查怪兽的user id 是我方user id
			if selfUid == helper.StringToInt(matchInfo[j]["user_id"]) {
				//a为目前的正面buff
				a = selfpositive[i]
				//检查a的正面buff数值
				switch a {
				//初始被动:全体攻击UP 2%
				case 9001:
					//攻击 * 1.02
					c = float64(helper.StringToInt(matchInfo[j]["monster_attack"])) * 1.02
					//进位(四舍五入)
					b = round(c)
					//存入matchinfo
					matchInfo[j]["monster_attack"] = strconv.Itoa(b)
				//初始被动:全体防御UP 2%
				case 9002:
					//防御 * 1.02
					c = float64(helper.StringToInt(matchInfo[j]["monster_defend"])) * 1.02
					//进位(四舍五入)
					b = round(c)
					//存入matchinfo
					matchInfo[j]["monster_defend"] = strconv.Itoa(b)
				//初始被动:全体速度UP 2%
				case 9003:
					//生命 * 1.02
					c = float64(helper.StringToInt(matchInfo[j]["monster_speed"])) * 1.02
					//进位(四舍五入)
					b = round(c)
					//存入matchinfo
					matchInfo[j]["monster_speed"] = strconv.Itoa(b)
				//初始被动:全体生命UP 2%
				case 9004:
					//速度 * 1.02
					c = float64(helper.StringToInt(matchInfo[j]["monster_hp"])) * 1.02
					//进位(四舍五入)
					b = round(c)
					//存入matchinfo
					matchInfo[j]["monster_hp"] = strconv.Itoa(b)
					matchInfo[j]["monster_max_hp"] = strconv.Itoa(b)
				//初始被动:全体闪避UP 2%
				case 9005:
					b = helper.StringToInt(matchInfo[j]["monster_hit"]) + 20
					//存入matchinfo
					matchInfo[j]["monster_hit"] = strconv.Itoa(b)
				}
			}
		}
	}
	//根据每个我方的负面buff
	for i := 0; i < len(selfnegative); i++ {
		//根据每个matchInfo(每只怪兽)
		for j := 0; j < len(matchInfo); j++ {
			//检查怪兽的user id 是我方user id
			if selfUid == helper.StringToInt(matchInfo[j]["user_id"]) {
				//获取目前的负面buff
				a = selfnegative[i]
				//检查a的负面buff数值
				switch a {
				//初始被动:全体攻击DOWN 2%
				case 9011:
					//攻击 * 0.98
					c = float64(helper.StringToInt(matchInfo[j]["monster_attack"])) * 0.98
					//进位(四舍五入)
					b = round(c)
					//存入matchinfo
					matchInfo[j]["monster_attack"] = strconv.Itoa(b)
				//初始被动:全体防御DOWN 2%
				case 9012:
					//防守 * 0.98
					c = float64(helper.StringToInt(matchInfo[j]["monster_defend"])) * 0.98
					//进位(四舍五入)
					b = round(c)
					//存入matchinfo
					matchInfo[j]["monster_defend"] = strconv.Itoa(b)
				//初始被动:全体速度DOWN 2%
				case 9013:
					//速度 * 0.98
					c = float64(helper.StringToInt(matchInfo[j]["monster_hp"])) * 0.98
					//进位(四舍五入)
					b = round(c)
					//存入matchinfo
					matchInfo[j]["monster_hp"] = strconv.Itoa(b)
					matchInfo[j]["monster_max_hp"] = strconv.Itoa(b)
				//初始被动:全体生命DOWN 2%
				case 9014:
					//生命 * 0.98
					c = float64(helper.StringToInt(matchInfo[j]["monster_speed"])) * 0.98
					//进位(四舍五入)
					b = round(c)
					//存入matchinfo
					matchInfo[j]["monster_speed"] = strconv.Itoa(b)
				//初始被动:全体闪避DOWN 2%
				case 9015:
					b = helper.StringToInt(matchInfo[j]["monster_hit"]) - 10
					//存入matchinfo
					matchInfo[j]["monster_hit"] = strconv.Itoa(b)
				}
			}
		}
	}

	//根据每个敌方的正面buff
	for i := 0; i < len(enemypositive); i++ {
		//根据每个matchInfo(每只怪兽)
		for j := 0; j < len(matchInfo); j++ {
			//检查怪兽的user id 不是我方user id
			if selfUid != helper.StringToInt(matchInfo[j]["user_id"]) {
				//获取敌人的正面Buff
				a = enemypositive[i]
				switch a {
				//初始被动:全体攻击UP 2%
				case 9001:
					//攻击 * 1.02
					c = float64(helper.StringToInt(matchInfo[j]["monster_attack"])) * 1.02
					//进位(四舍五入)
					b = round(c)
					//存入matchinfo
					matchInfo[j]["monster_attack"] = strconv.Itoa(b)
				//初始被动:全体防御UP 2%
				case 9002:
					//防御 * 1.02
					c = float64(helper.StringToInt(matchInfo[j]["monster_defend"])) * 1.02
					//进位(四舍五入)
					b = round(c)
					//存入matchinfo
					matchInfo[j]["monster_defend"] = strconv.Itoa(b)
				//初始被动:全体速度+ 1
				case 9003:
					//速度 +1
					c = float64(helper.StringToInt(matchInfo[j]["monster_speed"])) + 1
					//进位(四舍五入)
					b = round(c)
					//存入matchinfo
					matchInfo[j]["monster_speed"] = strconv.Itoa(b)
				//初始被动:全体生命UP 2%
				case 9004:
					//生命 * 1.02
					c = float64(helper.StringToInt(matchInfo[j]["monster_hp"])) * 1.02
					//进位(四舍五入)
					b = round(c)
					//存入matchinfo
					matchInfo[j]["monster_hp"] = strconv.Itoa(b)
					matchInfo[j]["monster_max_hp"] = strconv.Itoa(b)
				//初始被动:全体闪避UP 2%
				case 9005:
					b = helper.StringToInt(matchInfo[j]["monster_hit"]) + 20
					//存入matchinfo
					matchInfo[j]["monster_hit"] = strconv.Itoa(b)
				}
			}
		}
	}
	//根据每个敌方的负面buff
	for i := 0; i < len(enemynegative); i++ {
		//根据每个matchInfo(每只怪兽)
		for j := 0; j < len(matchInfo); j++ {
			//检查怪兽的user id 不是我方user id
			if selfUid != helper.StringToInt(matchInfo[j]["user_id"]) {
				//获取敌方负面buff
				a = enemynegative[i]
				//根据敌方的负面buff
				switch a {
				//初始被动:全体攻击DOWN 2%
				case 9011:
					//攻击 * 0.98
					c = float64(helper.StringToInt(matchInfo[j]["monster_attack"])) * 0.98
					//进位(四舍五入)
					b = round(c)
					//存入matchinfo
					matchInfo[j]["monster_attack"] = strconv.Itoa(b)
				//初始被动:全体防御DOWN 2%
				case 9012:
					//防守 * 0.98
					c = float64(helper.StringToInt(matchInfo[j]["monster_defend"])) * 0.98
					//进位(四舍五入)
					b = round(c)
					//存入matchinfo
					matchInfo[j]["monster_defend"] = strconv.Itoa(b)
				//初始被动:全体速度DOWN 2%
				case 9013:
					//速度 * 0.98
					c = float64(helper.StringToInt(matchInfo[j]["monster_speed"])) - 1
					//进位(四舍五入)
					b = round(c)
					//存入matchinfo
					if b > 0 {
						matchInfo[j]["monster_speed"] = strconv.Itoa(b)
					} else {
						matchInfo[j]["monster_speed"] = strconv.Itoa(1)
					}

				//初始被动:全体生命DOWN 2%
				case 9014:
					//生命 * 0.98
					c = float64(helper.StringToInt(matchInfo[j]["monster_hp"])) * 0.98
					//进位(四舍五入)
					b = round(c)
					//存入matchinfo
					matchInfo[j]["monster_hp"] = strconv.Itoa(b)
					matchInfo[j]["monster_max_hp"] = strconv.Itoa(b)
				//初始被动:全体闪避DOWN 2%
				case 9015:
					b = helper.StringToInt(matchInfo[j]["monster_hit"]) - 10
					//存入matchinfo
					matchInfo[j]["monster_hit"] = strconv.Itoa(b)
				}
			}
		}
	}

	return matchInfo, BuffInfo{
		Self:  self,
		Enemy: enemy,
	}
}

func getFullMatch(selfUid, selfEmbattleStartLen, enemyEmbattleStartLen int) []RoundInfo {
	//获取先后手
	prior := getPriorHand()
	x := 1
	var selfEmbattleCount, enemyEmbattleCount, count int
	var round RoundInfo
	var roundInfo []RoundInfo
	var turn TurnInfo
	var turnInfo []TurnInfo
	for {

		//我方怪兽初始数量
		selfEmbattleCount = selfEmbattleStartLen
		//敌方怪兽初始数量
		enemyEmbattleCount = enemyEmbattleStartLen
		//先后手交换
		if prior == 0 {
			prior = 1
		} else {
			prior = 0
		}
		//Array of monster with speed ascending order 怪兽根据速度顺序
		monstersWithSpeedOrder := getMonsterSpeed(battleInfo, prior, selfUid)
		turnInfo = []TurnInfo{}
		//根据speed生成turn， 从最后一个battleInfo逆序看起
		for i := len(battleInfo) - 1; i >= 0; i-- {
			//根据每个battleInfo
			for j := 0; j < len(battleInfo); j++ {
				//当前怪兽(monstersWithSpeedOrder)i, 等于当前怪兽的UID(battleInfo)
				if monstersWithSpeedOrder[i] == helper.StringToInt(battleInfo[j]["monster_uid"]) {
					//获取该回合的伤害记录
					turn = calculateDamage(monstersWithSpeedOrder[i], helper.StringToInt(battleInfo[j]["user_id"]), selfUid)
				}
			}
			//检查攻击怪兽是否存活（HP>0)
			if turn.AtkMonsterHp > 0 {
				//存活数量+1
				count += 1
				//加入turnInfo
				turnInfo = append(turnInfo, turn)
			}
		}
		//根据每个battleInfo
		for i := 0; i < len(battleInfo); i++ {
			//检查当前我方怪兽血量是否为0以下
			if helper.StringToInt(battleInfo[i]["monster_hp"]) <= 0 && helper.StringToInt(battleInfo[i]["user_id"]) == selfUid {
				//我方怪兽存活数量 - 1
				selfEmbattleCount -= 1
			}
			//检查当前敌方怪兽血量是否为0以下
			if helper.StringToInt(battleInfo[i]["monster_hp"]) <= 0 && helper.StringToInt(battleInfo[i]["user_id"]) != selfUid {
				//敌方怪兽存活数量 - 1
				enemyEmbattleCount -= 1
			}
		}
		//记录回合数，先后手，我方怪兽存活数量,敌方怪兽存活数量，还有回合资讯
		round = RoundInfo{
			Round:     x,
			PriorHand: prior,
			SelfLeft:  selfEmbattleCount,
			EnemyLeft: enemyEmbattleCount,
			Count:     count,
			TurnInfo:  turnInfo,
		}
		roundInfo = append(roundInfo, round)
		//回合数+1
		x++
		//检查任意1方怪兽存活数量为0
		if selfEmbattleCount == 0 || enemyEmbattleCount == 0 {
			//产生match结果
			return roundInfo
		}
	}
}

func getPriorHand() int {
	//根据现在随机产生种子码
	rand.Seed(time.Now().UnixNano())
	//随机return 0 或 1
	priorHand := rand.Intn(2)
	return priorHand
}

func getMonsterSpeed(matchInfo []map[string]string, prior, selfUid int) []int {
	var speed1, monster1, speed2, monster2, result []int
	//根据每个matchInfo
	for i := 0; i < len(matchInfo); i++ {
		//检查我方的user uid 是当前怪兽的user id
		if selfUid == helper.StringToInt(matchInfo[i]["user_id"]) {
			//检查先手
			if prior == 0 {
				//速度与怪兽uid放入1
				speed1 = append(speed1, helper.StringToInt(matchInfo[i]["monster_speed"]))
				monster1 = append(monster1, helper.StringToInt(matchInfo[i]["monster_uid"]))
			} else {
				//速度与怪兽uid放入2
				speed2 = append(speed2, helper.StringToInt(matchInfo[i]["monster_speed"]))
				monster2 = append(monster2, helper.StringToInt(matchInfo[i]["monster_uid"]))
			}
		}
	}
	//根据每个matchInfo
	for i := 0; i < len(matchInfo); i++ {
		//检查我方的user uid 不是当前怪兽的user id
		if selfUid != helper.StringToInt(matchInfo[i]["user_id"]) {
			//检查先手
			if prior == 0 {
				//速度与怪兽uid放入2
				speed2 = append(speed2, helper.StringToInt(matchInfo[i]["monster_speed"]))
				monster2 = append(monster2, helper.StringToInt(matchInfo[i]["monster_uid"]))
			} else {
				//速度与怪兽uid放入1
				speed1 = append(speed1, helper.StringToInt(matchInfo[i]["monster_speed"]))
				monster1 = append(monster1, helper.StringToInt(matchInfo[i]["monster_uid"]))
			}
		}
	}
	//把speed2 放入 speed1的后面
	speed1 = append(speed1, speed2...)
	//把monster2 放入 monster1的后面
	monster1 = append(monster1, monster2...)
	//排列顺序
	result = Sort(speed1, monster1)
	return result
}

func getFullSpeed(speed, side int) []map[string]string {
	speedInfo := map[string]string{
		"speed": strconv.Itoa(speed),
		"side":  strconv.Itoa(side),
	}
	speedCount = append(speedCount, speedInfo)
	return speedCount
}

func calculateDamage(atkmonster, uid, selfUid int) TurnInfo {
	var mulInfo []map[string]string
	var mulInfoArray []map[string]string
	var atkMonsterATK, atkMonsterDEF, atkMonsterHP, atkMonsterMAXHP, atkMonsterHIT, atkMonsterSPEED, atkMonsterSKILL, atkMonsterSkillRate, atkMonsterELEMENT int
	var atkMonsterMISS int
	var atkMonsterPositiveBuff, atkMonsterNegativeBuff []int
	var atkMonsterId, atkMonsterSequence int
	var heal = 0
	var damage = 0
	var aoedamage = 0
	var aoeattack = false
	var ignoreDefense = false
	var selfhealApply = false
	var down = 0.95
	var up = 1.05
	var blockDefense = 0
	var stun = false
	var atkmonsterInfo map[string]string
	var targetAliveMonsters []map[string]string
	var targetInfo map[string]string
	for a := 0; a < len(battleInfo); a++ {
		if helper.StringToInt(battleInfo[a]["monster_uid"]) == atkmonster {
			atkmonsterInfo = battleInfo[a]

			atkMonsterId = helper.StringToInt(atkmonsterInfo["monster_id"])
			atkMonsterATK = helper.StringToInt(atkmonsterInfo["monster_attack"])           //攻击
			atkMonsterDEF = helper.StringToInt(atkmonsterInfo["monster_defend"])           //防御
			atkMonsterHP = helper.StringToInt(atkmonsterInfo["monster_hp"])                //生命值
			atkMonsterMAXHP = helper.StringToInt(atkmonsterInfo["monster_max_hp"])         //生命值上限
			atkMonsterHIT = helper.StringToInt(atkmonsterInfo["monster_hit"])              //命中
			atkMonsterMISS = helper.StringToInt(atkmonsterInfo["monster_miss"])            //命中
			atkMonsterSkillRate = helper.StringToInt(atkmonsterInfo["monster_skill_rate"]) //技能触发概率
			atkMonsterSPEED = helper.StringToInt(atkmonsterInfo["monster_speed"])
			atkMonsterSKILL = helper.StringToInt(atkmonsterInfo["monster_skill"])
			atkMonsterELEMENT = helper.StringToInt(atkmonsterInfo["monster_element"])
			atkMonsterSequence = helper.StringToInt(atkmonsterInfo["monster_sequence"])
			atkMonsterPositiveBuff = convertBuffArray(atkmonsterInfo["monster_positive"])
			atkMonsterNegativeBuff = convertBuffArray(atkmonsterInfo["monster_negative"])
			damage = atkMonsterATK
			aoedamage = atkMonsterATK
		} else if helper.StringToInt(battleInfo[a]["user_id"]) != uid && helper.StringToInt(battleInfo[a]["monster_hp"]) > 0 {
			targetAliveMonsters = append(targetAliveMonsters, battleInfo[a])
		}
	}
	var selfSide = (selfUid == helper.StringToInt(atkmonsterInfo["user_id"]))
	var atkMonsterTriggerSkill = 0
	var atkMonsterOwnedTriggerSkill = false
	var attackSuccess = false
	var triggerSuccess = false
	var atkMonsterAlive = true
	atkMonsterAlive = atkMonsterHP > 0
	if atkMonsterAlive {
		//speedCount = getFullSpeed(atkmonster, atkside)
		//atk = 1
		var stunSkillNum = 9023
		if helper.ArrayContainsInt(atkMonsterNegativeBuff, stunSkillNum) {
			atkMonsterNegativeBuff = helper.ArrayRemoveInt(atkMonsterNegativeBuff, helper.ArrayIntIndex(atkMonsterNegativeBuff, stunSkillNum))
			damage = 0
			stun = true
		} else {
			if len(targetAliveMonsters) > 0 {
				targetAliveMonsters = sortTargetAliveMonsterBasedSequence(targetAliveMonsters)
				targetInfo = selectRandomTarget(targetAliveMonsters, len(battleInfo)/2)
			}
			if targetInfo != nil {
				attackSuccess = calculateAttackSuccess(atkMonsterHIT, helper.StringToInt(targetInfo["monster_miss"]))
				if attackSuccess {
					switch atkMonsterSKILL {
					case 9021, 9022, 9023, 9024, 9031, 9032, 9033, 9034, 9035,
						9051, 9052, 9053, 9054, 9061, 9062, 9063, 9064, 9065:
						atkMonsterOwnedTriggerSkill = true
					default:
						atkMonsterOwnedTriggerSkill = false
					}
					if atkMonsterOwnedTriggerSkill {
						triggerSuccess = calculateTriggerSuccess(atkMonsterSkillRate)
						// triggerSuccess = true
						if triggerSuccess {
							atkMonsterTriggerSkill = atkMonsterSKILL
							switch atkMonsterSKILL {
							case 9021:
								damage = int(float64(damage) * 1.5)
							case 9022:
								aoeattack = true
							case 9023:
								target_monster_negative := convertBuffArray(targetInfo["monster_negative"])
								target_monster_negative = append(target_monster_negative, 9023)
								targetInfo["monster_negative"] = convertBuffString(target_monster_negative)
							case 9024:
								selfhealApply = true
							case 9031:
								target_monster_attack := helper.StringToInt(targetInfo["monster_attack"])
								target_monster_attack = int(math.Ceil(float64(target_monster_attack) * down))
								targetInfo["monster_attack"] = strconv.Itoa(target_monster_attack)
							case 9032:
								target_monster_defend := helper.StringToInt(targetInfo["monster_defend"])
								target_monster_defend = int(math.Ceil(float64(target_monster_defend) * down))
								targetInfo["monster_defend"] = strconv.Itoa(target_monster_defend)
							case 9033:
								target_monster_speed := helper.StringToInt(targetInfo["monster_speed"]) - 1
								if target_monster_speed <= 0 {
									target_monster_speed = 1
								}
								targetInfo["monster_speed"] = strconv.Itoa(target_monster_speed)
							case 9034:
								target_monster_max_hp := helper.StringToInt(targetInfo["monster_max_hp"])
								target_monster_max_hp = int(float64(target_monster_max_hp) * down)
								if helper.StringToInt(targetInfo["monster_hp"]) > target_monster_max_hp {
									targetInfo["monster_hp"] = strconv.Itoa(target_monster_max_hp)
								}
							case 9035:
								target_monster_hit := helper.StringToInt(targetInfo["monster_hit"])
								target_monster_hit = int(math.Ceil(float64(target_monster_hit) * down))
								targetInfo["monster_hit"] = strconv.Itoa(target_monster_hit)
							case 9051:
								heal = int(float64(atkMonsterMAXHP) * 0.1)
								originalHP := atkMonsterHP
								originalHP += heal
								if originalHP > atkMonsterMAXHP {
									heal = atkMonsterMAXHP - atkMonsterHP
									atkMonsterHP = atkMonsterMAXHP
								} else {
									atkMonsterHP = originalHP
								}
							case 9052:
								for b := 0; b < len(battleInfo); b++ {
									if helper.StringToInt(battleInfo[b]["user_id"]) == uid {
										var groupHeal = int(math.Ceil(float64(helper.StringToInt(battleInfo[b]["monster_max_hp"])) * 0.05))
										var originalHP = helper.StringToInt(battleInfo[b]["monster_hp"])
										originalHP += groupHeal
										if originalHP > helper.StringToInt(battleInfo[b]["monster_max_hp"]) {
											groupHeal = helper.StringToInt(battleInfo[b]["monster_max_hp"]) - helper.StringToInt(battleInfo[b]["monster_hp"])
											battleInfo[b]["monster_hp"] = battleInfo[b]["monster_max_hp"]
										} else {
											battleInfo[b]["monster_hp"] = strconv.Itoa(originalHP)
										}
										mulInfo = append(mulInfo, map[string]string{
											"id":   battleInfo[b]["monster_uid"],
											"heal": strconv.Itoa(groupHeal),
											"hp":   battleInfo[b]["monster_hp"],
											"seq":  battleInfo[b]["monster_sequence"],
										})
									}
								}
							case 9053:
								if !helper.ArrayContainsInt(atkMonsterPositiveBuff, 9053) {
									atkMonsterPositiveBuff = append(atkMonsterPositiveBuff, 9053)
								}
							case 9054:
								ignoreDefense = true
							case 9061:
								atkMonsterATK = int(math.Ceil(float64(atkMonsterATK) * up))
							case 9062:
								atkMonsterDEF = int(math.Ceil(float64(atkMonsterDEF) * up))
							case 9063:
								atkMonsterSPEED += 1
							case 9064:
								atkMonsterMAXHP = int(math.Ceil(float64(atkMonsterMAXHP) * up))
							case 9065:
								atkMonsterHIT += 5
							}
						}
					}
					targetInfo, damage, blockDefense = blockDefenseOrIgnoreDefense(targetInfo, damage, ignoreDefense)
					if damage != 0 {
						damage = calcElementDmg(damage, atkMonsterELEMENT, helper.StringToInt(targetInfo["monster_element"]))
						damage -= blockDefense
						if damage <= 0 {
							damage = 1
						}
					}
					targetInfo["monster_hp"] = strconv.Itoa(helper.StringToInt(targetInfo["monster_hp"]) - damage)
					if helper.StringToInt(targetInfo["monster_hp"]) < 0 {
						targetInfo["monster_hp"] = "0"
					}
					if selfhealApply {
						heal = int(math.Ceil(float64(damage) * 0.5))
						originalHP := atkMonsterHP
						originalHP += heal
						if originalHP > atkMonsterMAXHP {
							heal = atkMonsterMAXHP - atkMonsterHP
							atkMonsterHP = atkMonsterMAXHP
						} else {
							atkMonsterHP = originalHP
						}
					}
					if aoeattack {
						for d := 0; d < len(battleInfo); d++ {
							if helper.StringToInt(battleInfo[d]["user_id"]) != uid {
								if battleInfo[d]["monster_uid"] != targetInfo["monster_uid"] {
									if helper.StringToInt(battleInfo[d]["monster_hp"]) > 0 {
										var thisaoeDamage = aoedamage
										var aoeblockDefense = helper.StringToInt(battleInfo[d]["monster_defend"])
										battleInfo[d], thisaoeDamage, aoeblockDefense = blockDefenseOrIgnoreDefense(battleInfo[d], thisaoeDamage, ignoreDefense)
										if thisaoeDamage != 0 {
											thisaoeDamage -= aoeblockDefense
											if thisaoeDamage <= 0 {
												thisaoeDamage = 1
											}
										}
										battleInfo[d]["monster_hp"] = strconv.Itoa(helper.StringToInt(battleInfo[d]["monster_hp"]) - thisaoeDamage)
										if helper.StringToInt(battleInfo[d]["monster_hp"]) < 0 {
											battleInfo[d]["monster_hp"] = "0"
										}
										mulInfo = append(mulInfo, map[string]string{
											"id":     battleInfo[d]["monster_uid"],
											"dmg":    strconv.Itoa(thisaoeDamage),
											"hpdown": battleInfo[d]["monster_hp"],
											"seq":    battleInfo[d]["monster_sequence"],
										})
									}
								} else {
									mulInfo = append(mulInfo, map[string]string{
										"id":     targetInfo["monster_uid"],
										"dmg":    strconv.Itoa(damage),
										"hpdown": targetInfo["monster_hp"],
										"seq":    targetInfo["monster_sequence"],
									})
								}
							}
						}
					}
				} else {
					damage = 0
				}
			}
		}
	}
	if len(mulInfo) > 0 {
		mulInfoArray = mulInfo
	}
	for c := 0; c < len(battleInfo); c++ {
		if targetInfo["monster_uid"] == battleInfo[c]["monster_uid"] {
			battleInfo[c]["monster_hp"] = targetInfo["monster_hp"]
			battleInfo[c]["monster_attack"] = targetInfo["monster_attack"]
			battleInfo[c]["monster_defend"] = targetInfo["monster_defend"]
			battleInfo[c]["monster_max_hp"] = targetInfo["monster_max_hp"]
			battleInfo[c]["monster_speed"] = targetInfo["monster_speed"]
			battleInfo[c]["monster_miss"] = targetInfo["monster_miss"]
			battleInfo[c]["monster_hit"] = targetInfo["monster_hit"]
			battleInfo[c]["monster_skill_rate"] = targetInfo["monster_skill_rate"]
			battleInfo[c]["monster_positive"] = targetInfo["monster_positive"]
			battleInfo[c]["monster_negative"] = targetInfo["monster_negative"]
		} else if atkmonster == helper.StringToInt(battleInfo[c]["monster_uid"]) {
			battleInfo[c]["monster_hp"] = strconv.Itoa(atkMonsterHP)
			battleInfo[c]["monster_attack"] = strconv.Itoa(atkMonsterATK)
			battleInfo[c]["monster_defend"] = strconv.Itoa(atkMonsterDEF)
			battleInfo[c]["monster_max_hp"] = strconv.Itoa(atkMonsterMAXHP)
			battleInfo[c]["monster_speed"] = strconv.Itoa(atkMonsterSPEED)
			battleInfo[c]["monster_miss"] = strconv.Itoa(atkMonsterMISS)
			battleInfo[c]["monster_hit"] = strconv.Itoa(atkMonsterHIT)
			battleInfo[c]["monster_skill_rate"] = strconv.Itoa(atkMonsterSkillRate)
			battleInfo[c]["monster_positive"] = convertBuffString(atkMonsterPositiveBuff)
			battleInfo[c]["monster_negative"] = convertBuffString(atkMonsterNegativeBuff)
		}
	}
	if targetInfo != nil {
		return TurnInfo{
			AtkMonsterUid:         atkMonsterId,
			AtkMonsterId:          atkmonster,
			AtkMonsterSide:        selfSide,
			AtkMonsterSeq:         atkMonsterSequence,
			AtkMonsterHp:          atkMonsterHP,
			AtkMonsterDamage:      damage,
			AtkMonsterElement:     atkMonsterELEMENT,
			AtkMonsterSkill:       atkMonsterSKILL,
			AtkMonsterHeal:        heal,
			AtkMonsterNegative:    convertBuffString(atkMonsterNegativeBuff),
			AtkMonsterPositive:    convertBuffString(atkMonsterPositiveBuff),
			AtkMonsterTrigger:     atkMonsterTriggerSkill,
			AtkMonsterStun:        stun,
			AtkMultiple:           mulInfoArray,
			AtkMonsterTriggered:   triggerSuccess,
			AtkMonsterMiss:        attackSuccess,
			TargetMonsterId:       helper.StringToInt(targetInfo["monster_uid"]),
			TargetMonsterSeq:      helper.StringToInt(targetInfo["monster_sequence"]),
			TargetMonsterHp:       helper.StringToInt(targetInfo["monster_hp"]),
			TargetMonsterAttack:   helper.StringToInt(targetInfo["monster_attack"]),
			TargetMonsterDefense:  helper.StringToInt(targetInfo["monster_defend"]),
			TargetMonsterSpeed:    helper.StringToInt(targetInfo["monster_speed"]),
			TargetMonsterNegative: targetInfo["monster_negative"],
		}
	}
	return TurnInfo{
		AtkMonsterUid:         atkMonsterId,
		AtkMonsterId:          atkmonster,
		AtkMonsterSide:        selfSide,
		AtkMonsterSeq:         atkMonsterSequence,
		AtkMonsterHp:          atkMonsterHP,
		AtkMonsterDamage:      damage,
		AtkMonsterElement:     atkMonsterELEMENT,
		AtkMonsterSkill:       atkMonsterSKILL,
		AtkMonsterHeal:        heal,
		AtkMonsterNegative:    convertBuffString(atkMonsterNegativeBuff),
		AtkMonsterPositive:    convertBuffString(atkMonsterPositiveBuff),
		AtkMonsterTrigger:     atkMonsterTriggerSkill,
		AtkMonsterStun:        stun,
		AtkMultiple:           mulInfoArray,
		AtkMonsterTriggered:   triggerSuccess,
		AtkMonsterMiss:        attackSuccess,
		TargetMonsterId:       -1,
		TargetMonsterSeq:      -1,
		TargetMonsterHp:       -1,
		TargetMonsterAttack:   -1,
		TargetMonsterDefense:  -1,
		TargetMonsterSpeed:    -1,
		TargetMonsterNegative: targetInfo["monster_negative"],
	}
}

// func calculateDamage(atkmonster, uid, selfUid int) TurnInfo {
// 	var atkside, atkseq, atkhp, atkdmg, atkarmor, atkspeed, atkmiss, atkhit, atkelement, atkskill, atktrigger int
// 	var targetid, targetseq, targethp, targetdmg, targetarmor, targetspeed, targetmiss, targetelement int
// 	var totaldmg, atk, trigger, attack, effect, target, skill, atkuid, aoeatk, ignore, stun, atkuuid int
// 	var atkheal, enemydown, selfup, down, up, healing, allhealing, currenthp, heal float64
// 	var atkpositive, atknegative, targetpositive, targetnegative string
// 	var atkmultiple []map[string]string
// 	var mulInfo map[string]string

// 	//检查每个battleInfo
// 	for i := 0; i < len(battleInfo); i++ {
// 		//检查攻击的怪兽等于为battleInfo的怪兽 （用monster uid来辨认）
// 		if atkmonster == helper.StringToInt(battleInfo[i]["monster_uid"]) {
// 			//从battleInfo获取攻击怪兽的图片id
// 			atkuid = helper.StringToInt(battleInfo[i]["monster_id"])
// 			//从battleInfo获取攻击怪兽攻击数值
// 			atkdmg = helper.StringToInt(battleInfo[i]["monster_attack"])
// 			//从battleInfo获取攻击怪兽防御数值
// 			atkarmor = helper.StringToInt(battleInfo[i]["monster_defend"])
// 			//从battleInfo获取攻击怪兽生命数值
// 			atkhp = helper.StringToInt(battleInfo[i]["monster_hp"])
// 			//从battleInfo获取攻击怪兽速度数值
// 			atkspeed = helper.StringToInt(battleInfo[i]["monster_speed"])
// 			//从battleInfo获取攻击怪兽命中率数值
// 			atkhit = helper.StringToInt(battleInfo[i]["monster_hit"])
// 			//从battleInfo获取攻击怪兽闪避率数值
// 			atkmiss = helper.StringToInt(battleInfo[i]["monster_miss"])
// 			//从battleInfo获取攻击怪兽元素属性
// 			atkelement = helper.StringToInt(battleInfo[i]["monster_element"])
// 			//从battleInfo获取攻击怪兽技能ID
// 			atkskill = helper.StringToInt(battleInfo[i]["monster_skill"])
// 			//从battleInfo获取攻击怪兽技能触发ID
// 			trigger = helper.StringToInt(battleInfo[i]["monster_trigger"])
// 			//从battleInfo获取攻击怪兽的user UID
// 			atkuuid = helper.StringToInt(battleInfo[i]["user_id"])
// 			//从battleInfo获取攻击怪兽所拥有的正面buff
// 			atkpositive = battleInfo[i]["monster_positive"]
// 			//从battleInfo获取攻击怪兽所拥有的负面buff
// 			atknegative = battleInfo[i]["monster_negative"]
// 			//检查攻击怪兽的user UID 等于 玩家的UID
// 			if atkuuid == selfUid {
// 				atkside = 0
// 			} else {
// 				atkside = 1
// 			}
// 			//从battleInfo获取攻击怪兽的位置顺序
// 			atkseq = helper.StringToInt(battleInfo[i]["monster_sequence"])
// 			//检查攻击怪兽还活着
// 			if atkhp > 0 {
// 				//?
// 				speedCount = getFullSpeed(atkmonster, atkside)
// 				//?
// 				atk = 1
// 				//根据目标的获取攻击命中成功或者失败， 0为命中成功， 1为敌方闪避成功
// 				attack = triggerEffect1(atkhit)
// 				//检查攻击怪兽有正面buff （不等于空）
// 				if atkpositive != "" {
// 					//9061 = 触发被动:20%单体攻击UP 5%
// 					if atkpositive[:4] == "9061" {
// 						battleInfo[i]["monster_attack"] = initBattle[i]["monster_attack"]
// 					}
// 					//9062 = 触发被动:20%单体防御UP 5%
// 					if atkpositive[:4] == "9062" {
// 						battleInfo[i]["monster_armor"] = initBattle[i]["monster_armor"]
// 					}
// 					//9063 = 触发被动:20%单体速度UP 5%
// 					if atkpositive[:4] == "9063" {
// 						battleInfo[i]["monster_speed"] = initBattle[i]["monster_speed"]
// 					}
// 					//9064 = 触发被动:20%单体生命UP 5%
// 					if atkpositive[:4] == "9064" {
// 						battleInfo[i]["monster_hp"] = initBattle[i]["monster_hp"]
// 					}
// 					//9065 = 触发被动:20%单体闪避UP 5%
// 					if atkpositive[:4] == "9065" {
// 						battleInfo[i]["monster_hit"] = initBattle[i]["monster_hit"]
// 					}
// 					//把前4位和逗号(共5个char)去除，获取剩余部分
// 					atkpositive = atkpositive[5:]
// 				}
// 				//检查命中成功
// 				if attack == 0 {
// 					//检查攻击方所被给予的负面buff不为9023（'触发被动:20%晕眩攻击'） 或者 等于无
// 					if atknegative == "" || atknegative[:4] != "9023" {
// 						//检查攻击方被给予的负面buff不为无
// 						if atknegative != "" {
// 							//9031 - 触发被动:20%单体攻击DOWN 5%
// 							if atknegative[:4] == "9031" {
// 								battleInfo[i]["monster_attack"] = initBattle[i]["monster_attack"]
// 							}
// 							//9032 - 触发被动:20%单体防御DOWN 5%
// 							if atknegative[:4] == "9032" {
// 								battleInfo[i]["monster_armor"] = initBattle[i]["monster_armor"]
// 							}
// 							//9033 - 触发被动:20%单体速度DOWN 5%
// 							if atknegative[:4] == "9033" {
// 								battleInfo[i]["monster_speed"] = initBattle[i]["monster_speed"]
// 							}
// 							//9034 - 触发被动:20%单体生命DOWN 5%
// 							if atknegative[:4] == "9034" {
// 								battleInfo[i]["monster_hp"] = initBattle[i]["monster_hp"]
// 							}
// 							//9035 - 触发被动:20%单体闪避DOWN 5%
// 							if atknegative[:4] == "9035" {
// 								battleInfo[i]["monster_hit"] = initBattle[i]["monster_hit"]
// 							}
// 							//把前4位和逗号(共5个char)去除，获取剩余部分
// 							atknegative = atknegative[5:]
// 						}
// 						//检查攻击方怪兽拥有技能ID
// 						switch atkskill {
// 						//触发被动技能
// 						case 9021, 9022, 9023, 9024, 9031, 9032, 9033, 9034, 9035, 9051, 9052, 9053, 9054, 9061, 9062, 9063, 9064, 9065:
// 							//有触发被动技能 = 1
// 							skill = 1
// 						}
// 						//
// 						atkheal = float64(0)
// 						enemydown = 1
// 						selfup = 1
// 						aoeatk = 0
// 						stun = 0
// 						atktrigger = 0
// 						up = float64(105) / float64(100)  // 105%
// 						down = float64(95) / float64(100) // 95%
// 						//检查攻击怪兽有被动触发技能
// 						if skill == 1 {
// 							//随机20%成功触发， 1为成功触发， 0 为失败触发
// 							effect = triggerEffect()
// 							if effect == 1 { //触发
// 								if trigger == 3 { // 对对手触发
// 									switch atkskill { //根据技能ID
// 									case 9021: //触发被动:20%暴击 150%
// 										atkdmg = int(float64(atkdmg) * float64(150) / float64(100)) //攻击伤害150%
// 										atktrigger = 9021
// 									case 9022: //触发被动:20%范围攻击
// 										aoeatk = 1 ///变成范围攻击
// 										atktrigger = 9022
// 									case 9023: //触发被动:20%晕眩攻击
// 										stun = 1 //变成晕眩攻击
// 										atktrigger = 9023
// 									case 9024: //触发被动:20%吸血攻击
// 										atkheal = float64(atkdmg) * float64(50) / float64(100) //恢复50%攻击伤害的血量
// 										atktrigger = 9024
// 									case 9031: //触发被动:20%单体攻击DOWN 5%
// 										enemydown = down //95%
// 										atktrigger = 9031
// 									case 9032: //触发被动:20%单体防御DOWN 5%
// 										enemydown = down //95%
// 										atktrigger = 9032
// 									case 9033: //触发被动:20%单体速度DOWN 5%
// 										enemydown = down //95%
// 										atktrigger = 9033
// 									case 9034: //触发被动:20%单体生命DOWN 5%
// 										enemydown = down //95%
// 										atktrigger = 9034
// 									case 9035: //触发被动:20%单体闪避DOWN 5%
// 										enemydown = down //95%
// 										atktrigger = 9035
// 									}
// 								} else if trigger == 4 { // 对自己触发
// 									switch atkskill { //根据技能ID
// 									case 9051: //触发被动:20%单体回复UP 10%
// 										atktrigger = 9051
// 										healing = float64(110) / float64(100) //110%
// 									case 9052: //触发被动:20%全体回复UP 5%
// 										atktrigger = 9052
// 										allhealing = float64(5) / float64(100) //5%
// 									case 9053: //触发被动:20%单体护盾(抵挡一次伤害)
// 										atktrigger = 9053
// 										atkpositive = atkpositive + strconv.Itoa(atktrigger) + "," //给与护盾
// 									case 9054: //触发被动:20%无视防御
// 										ignore = 1 //变为无视防御攻击
// 										atktrigger = 9054
// 									case 9061: //触发被动:20%单体攻击UP 5%
// 										selfup = up //105%
// 										atktrigger = 9061
// 									case 9062: //触发被动:20%单体防御UP 5%
// 										selfup = up //105%
// 										atktrigger = 9062
// 									case 9063: //触发被动:20%单体速度UP 5%
// 										selfup = up //105%
// 										atktrigger = 9063
// 									case 9064: //触发被动:20%单体生命UP 5%
// 										selfup = up //105%
// 										atktrigger = 9064
// 									case 9065: //触发被动:20%单体闪避UP 5%
// 										selfup = up //105%
// 										atktrigger = 9065
// 									}
// 								}
// 							}
// 						}
// 					} else {
// 						//把前4位和逗号(共5个char)去除，获取剩余部分
// 						atknegative = atknegative[5:]
// 						//攻击伤害为0
// 						atkdmg = 0
// 					}
// 				} else {
// 					//攻击伤害为0
// 					atkdmg = 0
// 				}
// 				//触发被动:20%吸血攻击
// 				if atktrigger == 9024 {
// 					//获取恢复后的血量（可能超出）
// 					currenthp = float64(atkhp) * atkheal
// 					//检查恢复血量有没有超出初始血量上限
// 					if currenthp > float64(helper.StringToInt(initBattle[i]["monster_hp"])) {
// 						//血量变为初始血量上限（超出）
// 						atkhp = helper.StringToInt(initBattle[i]["monster_hp"])
// 					} else {
// 						//血量变为恢复后血量（没超出）
// 						atkhp = int(currenthp)
// 					}
// 				}
// 				//检查有对我方触发的触发被动技能
// 				if atktrigger != 0 && trigger == 4 {
// 					//把触发被动技能加入攻击方正面buff
// 					atkpositive = atkpositive + strconv.Itoa(atktrigger) + ","
// 					//触发被动:20%单体回复UP 10%
// 					if atktrigger == 9051 {
// 						//攻击伤害为0
// 						atkdmg = 0
// 						//获取恢复后的血量（可能超出）
// 						currenthp = float64(atkhp) * healing
// 						//检查恢复血量有没有超出初始血量上限
// 						if currenthp > float64(helper.StringToInt(initBattle[i]["monster_hp"])) {
// 							//血量变为初始血量上限（超出）
// 							atkhp = helper.StringToInt(initBattle[i]["monster_hp"])
// 						} else {
// 							//血量变为恢复后血量（没超出）
// 							atkhp = int(currenthp)
// 						}
// 					}
// 					//触发被动:20%全体回复UP 5%
// 					if atktrigger == 9052 {
// 						//根据每个battleInfo
// 						for j := 0; j < len(battleInfo); j++ {
// 							//检查是我方的怪兽（user id) AND 该怪兽是否存活（HP > 0）
// 							if atkuuid == helper.StringToInt(battleInfo[j]["user_id"]) && helper.StringToInt(battleInfo[j]["monster_hp"]) > 0 {
// 								//获取恢复的血量
// 								heal = float64(helper.StringToInt(battleInfo[j]["monster_hp"])) * allhealing
// 								//获取恢复后的血量（可能超出）
// 								currenthp = float64(helper.StringToInt(battleInfo[j]["monster_hp"])) + heal
// 								//检查恢复血量有没有超出初始血量上限
// 								if currenthp > float64(helper.StringToInt(initBattle[j]["monster_hp"])) {
// 									//血量变为初始血量上限（超出）
// 									battleInfo[j]["monster_hp"] = initBattle[j]["monster_hp"]
// 								} else {
// 									//血量变为恢复后血量（没超出）
// 									battleInfo[j]["monster_hp"] = strconv.Itoa(int(currenthp))
// 								}
// 								//记录恢复后的血量，怪兽UID，位置，恢复的血量
// 								mulInfo = make(map[string]string)
// 								mulInfo["id"] = battleInfo[j]["monster_uid"]
// 								mulInfo["heal"] = strconv.Itoa(int(heal))
// 								mulInfo["hp"] = battleInfo[j]["monster_hp"]
// 								mulInfo["seq"] = battleInfo[j]["monster_sequence"]
// 								atkmultiple = append(atkmultiple, mulInfo)
// 							}
// 						}
// 					}
// 					//触发被动:20%单体攻击UP 5%
// 					if atktrigger == 9061 {
// 						//攻击伤害 * 1.05
// 						atkdmg = int(float64(atkdmg) * selfup)
// 						battleInfo[i]["monster_attack"] = strconv.Itoa(atkdmg)
// 					}
// 					//触发被动:20%单体防御UP 5%
// 					if atktrigger == 9062 {
// 						//防御数值 * 1.05
// 						atkarmor = int(float64(atkarmor) * selfup)
// 						battleInfo[i]["monster_armor"] = strconv.Itoa(atkarmor)
// 					}
// 					//触发被动:20%单体速度UP 5%
// 					if atktrigger == 9063 {
// 						//速度数值 * 1.05
// 						atkspeed = int(float64(atkspeed) * selfup)
// 						battleInfo[i]["monster_speed"] = strconv.Itoa(atkspeed)
// 					}
// 					//触发被动:20%单体生命UP 5%
// 					if atktrigger == 9064 {
// 						//生命数值 * 1.05
// 						atkhp = int(float64(atkhp) * selfup)
// 						battleInfo[i]["monster_hp"] = strconv.Itoa(atkhp)
// 					}
// 					//触发被动:20%单体闪避UP 5%
// 					if atktrigger == 9065 {
// 						//闪避数值 * 1.05
// 						atkhit = int(float64(atkhit) * selfup)
// 						battleInfo[i]["monster_hit"] = strconv.Itoa(atkhit)
// 					}
// 				}
// 				//根据每个battleInfo的怪兽
// 				for j := 0; j < len(battleInfo); j++ {
// 					//检查时敌方怪兽（不为自己的user_id) AND 该怪兽存活（HP > 0 ) AND 是否被攻击
// 					if uid != helper.StringToInt(battleInfo[j]["user_id"]) && helper.StringToInt(battleInfo[j]["monster_hp"]) > 0 && atk == 1 {
// 						//获取目标怪兽UID
// 						targetid = helper.StringToInt(battleInfo[j]["monster_uid"])
// 						//获取目标怪兽攻击数值
// 						targetdmg = helper.StringToInt(battleInfo[j]["monster_attack"])
// 						//获取目标怪兽防御数值
// 						targetarmor = helper.StringToInt(battleInfo[j]["monster_defend"])
// 						//获取目标怪兽生命数值
// 						targethp = helper.StringToInt(battleInfo[j]["monster_hp"])
// 						//获取目标怪兽速度数值
// 						targetspeed = helper.StringToInt(battleInfo[j]["monster_speed"])
// 						//获取目标怪兽命中/闪避数值
// 						targetmiss = helper.StringToInt(battleInfo[j]["monster_hit"])
// 						//获取目标怪兽的位置
// 						targetseq = helper.StringToInt(battleInfo[j]["monster_sequence"])
// 						//获取目标怪兽的元素属性
// 						targetelement = helper.StringToInt(battleInfo[j]["monster_element"])
// 						//获取目标怪兽的正面buff
// 						targetpositive = battleInfo[j]["monster_positive"]
// 						//获取目标怪兽的负面buff
// 						targetnegative = battleInfo[j]["monster_negative"]

// 						//检查双方属性是否为克制属性  1 > 2 > 3 > 4 > 5 > 1
// 						if atkelement == 5 { //属性克制
// 							if targetelement == 1 {
// 								atkdmg = atkdmg * 110 / 100
// 							}
// 						} else if atkelement != 5 {
// 							if atkelement == targetelement-1 {
// 								atkdmg = atkdmg * 110 / 100
// 							}
// 						}
// 						//检查触发被动技能不为空 AND 技能是对敌方触发
// 						if atktrigger != 0 && trigger == 3 {
// 							//触发被动:20%单体攻击DOWN 5%
// 							if atktrigger == 9031 {
// 								// targetdmg = targetdmg * int(enemydown)
// 								targetdmg = int(float64(targetdmg) * enemydown)
// 							}
// 							//触发被动:20%单体防御DOWN 5%
// 							if atktrigger == 9032 {
// 								// targetarmor = targetarmor * int(enemydown)
// 								targetarmor = int(float64(targetarmor) * enemydown)
// 							}
// 							//触发被动:20%单体速度DOWN 5%
// 							if atktrigger == 9033 {
// 								// targetspeed = targetspeed * int(enemydown)
// 								targetspeed = int(float64(targetspeed) * enemydown)
// 							}
// 							//触发被动:20%单体生命DOWN 5%
// 							if atktrigger == 9034 {
// 								// targethp = targethp * int(enemydown)
// 								targethp = int(float64(targethp) * enemydown)
// 							}
// 							//触发被动:20%单体闪避DOWN 5%
// 							if atktrigger == 9035 {
// 								// targetmiss = targetmiss * int(enemydown)
// 								targetmiss = int(float64(targetmiss) * enemydown)
// 							}
// 						}
// 						//检查晕眩攻击
// 						if stun == 1 {
// 							//加入敌方负面buff
// 							targetnegative = targetnegative + strconv.Itoa(atktrigger) + ","
// 						}
// 						//检查无视防御攻击
// 						if ignore == 1 {
// 							//敌方防御为0
// 							targetarmor = 0
// 							//重置无视防御攻击状态
// 							ignore = 0
// 						}
// 						//总伤害计算为 我方攻击数值 减 目标防御值
// 						totaldmg = atkdmg - targetarmor
// 						//总伤害 小于或等于 0时
// 						if totaldmg <= 0 { // armor higher than damage
// 							//总伤害为1
// 							totaldmg = 1
// 						}
// 						//检查攻击伤害为 0
// 						if atkdmg == 0 { // miss
// 							//总伤害为0
// 							totaldmg = 0
// 						}
// 						//检查有攻击方的正面buff
// 						if atkpositive != "" {
// 							//触发被动:20%全体回复UP 5%
// 							if atkpositive[:4] == "9052" { // heal
// 								totaldmg = 0 //总伤害为 0
// 							}
// 						}
// 						//检查目标方的正面buff
// 						if targetpositive != "" {
// 							//触发被动:20%单体护盾(抵挡一次伤害)
// 							if targetpositive[:4] == "9053" { // block
// 								totaldmg = 0 //总伤害为0
// 							}
// 						}
// 						//目标血量减少总伤害
// 						targethp = targethp - totaldmg
// 						//检查目标血量小于/等于 0 时
// 						if targethp <= 0 {
// 							//目标血量变为0，视为死亡
// 							targethp = 0
// 						}
// 						//存入怪兽的血量
// 						battleInfo[j]["monster_hp"] = strconv.Itoa(targethp)
// 						//检查啥？
// 						if target == 1 {
// 							//添加怪兽负面buff
// 							battleInfo[j]["monster_negative"] = strconv.Itoa(atkskill) + battleInfo[j]["monster_negative"]
// 							target = 0
// 						}
// 						//检查是范围攻击
// 						if aoeatk == 1 {
// 							//记录怪兽uid, 总伤害，生命数值，位置
// 							mulInfo = make(map[string]string)
// 							mulInfo["id"] = battleInfo[j]["monster_uid"]
// 							mulInfo["dmg"] = strconv.Itoa(totaldmg)
// 							mulInfo["hpdown"] = battleInfo[j]["monster_hp"]
// 							mulInfo["seq"] = battleInfo[j]["monster_sequence"]
// 							atkmultiple = append(atkmultiple, mulInfo)
// 						}
// 						//检查不是范围攻击，j不是最后一个
// 						if aoeatk == 0 || j == (len(battleInfo)-1) {
// 							atk = 0 //重置攻击状态
// 						}
// 					}
// 				}
// 			}
// 		}
// 	}
// 	return TurnInfo{
// 		AtkMonsterUid:         atkuid,
// 		AtkMonsterId:          atkmonster,
// 		AtkMonsterSide:        atkside,
// 		AtkMonsterSeq:         atkseq,
// 		AtkMonsterHp:          atkhp,
// 		AtkMonsterDamage:      totaldmg,
// 		AtkMonsterElement:     atkelement,
// 		AtkMonsterSkill:       atkskill,
// 		AtkMonsterHeal:        atkheal,
// 		AtkMonsterNegative:    atknegative,
// 		AtkMonsterPositive:    atkpositive,
// 		AtkMonsterTrigger:     atktrigger,
// 		AtkMultiple:           atkmultiple,
// 		AtkMonsterTriggered:   effect,
// 		AtkMonsterMiss:        attack,
// 		TargetMonsterId:       targetid,
// 		TargetMonsterSeq:      targetseq,
// 		TargetMonsterHp:       targethp,
// 		TargetMonsterNegative: targetnegative,
// 	}
// }

func getRankingHandler(r *http.Request) (map[string]interface{}, error) {

	uid, isValid := helper.VerifyJWT(r)
	if !isValid {
		return nil, errors.New("Invalid token")
	}
	reqJSON := helper.ReadParameters(r)
	rankNum := reqJSON["rankNum"]
	selfRankScore, err := db.GetRanking(uid, rankNum)

	if err != nil {
		return nil, errors.New("Invalid token")
	}
	payload := map[string]interface{}{
		"rankScore": selfRankScore,
		"rankNum":   rankNum,
	}
	return payload, nil
}

func setRankingHandler(r *http.Request) (int, error) {
	uid, isValid := helper.VerifyJWT(r)
	if !isValid {
		return 0, errors.New("Invalid token")
	}
	reqJSON := helper.ReadParameters(r)
	newRank := reqJSON["newRank"]
	rankNum := reqJSON["rankNum"]
	if rankNum == "3" || rankNum == "5" {
		db.UpdateRank(newRank, uid, rankNum)
		return 1, nil
	}
	return 0, errors.New("Invalid Rank Number")
}

func triggerEffect() int {
	if randomRange(1, 100) <= 20 {
		return 1
	} else {
		return 0
	}
}

func triggerEffect1(x int) int {
	if randomRange(1, 100) >= x {
		return 1
	} else {
		return 0
	}
}

func calculateAttackSuccess(hit int, miss int) bool {
	var calcHit int = hit - miss
	if rand.Intn(100) < calcHit {
		return true
	} else {
		return false
	}
}

func round(x float64) int { //四舍五入
	return int(math.Floor(x + 1.0))
}

func Sort(speed, monster []int) []int { //排序
	for i := 0; i < len(speed); i++ {
		min := i
		for j := i + 1; j < len(speed); j++ {
			//检查最小，找出最小speed的index
			if speed[j] < speed[min] {
				min = j
			}
		}
		tmp := speed[i]
		tmp1 := monster[i]
		speed[i] = speed[min]
		monster[i] = monster[min]
		speed[min] = tmp
		monster[min] = tmp1
	}
	return monster
}

func randomRange(min int, max int) int {
	return min + rand.Intn(max-min)
}

func giveStartBuff(uid string, matchInfo []map[string]string, self []map[string]string, enemy []map[string]string) ([]map[string]string, []map[string]string) {
	for a := 0; a < len(matchInfo); a++ {
		if matchInfo[a]["user_id"] == uid {
			for b := 0; b < len(self); b++ {
				if self[b]["monster_uid"] == matchInfo[a]["monster_uid"] {
					self[b] = matchInfo[a]
				}
			}
		} else {
			for c := 0; c < len(enemy); c++ {
				if enemy[c]["monster_uid"] == matchInfo[a]["monster_uid"] {
					enemy[c] = matchInfo[a]
				}
			}
		}
	}
	return self, enemy
}

func convertBuffArray(buffString string) []int {
	buffArray := []int{}
	if buffString == "" {
		buffArray = []int{}
	} else {
		var buffStringArray = strings.Split(buffString, ",")
		for a := 0; a > len(buffStringArray); a++ {
			buffArray = append(buffArray, helper.StringToInt(buffStringArray[a]))
		}
	}
	return buffArray
}

func convertBuffString(buffArray []int) string {
	var buffString string = ""
	if len(buffArray) > 0 {
		buffString = strconv.Itoa(buffArray[0])
		for a := 1; a > len(buffArray); a++ {
			buffString = buffString + "," + strconv.Itoa(buffArray[a])
		}
	}
	return buffString
}

func selectRandomTarget(targets []map[string]string, countMemberOfTeam int) map[string]string {

	switch countMemberOfTeam {
	case 3:
		//1
		if helper.StringToInt(targets[0]["monster_sequence"]) == 1 {
			return targets[0]
		} else {
			var backsideTarget []map[string]string
			for a := 0; a < len(targets); a++ {
				backsideTarget = append(backsideTarget, targets[a])
			}
			if len(backsideTarget) > 0 {
				rand.Seed(time.Now().UnixNano())
				return backsideTarget[rand.Intn(len(backsideTarget))]
			}
		}
		//23
	case 5:
		//12
		var backsideTarget []map[string]string
		for a := 0; a < len(targets); a++ {
			if helper.StringToInt(targets[a]["monster_sequence"]) == 1 || helper.StringToInt(targets[a]["monster_sequence"]) == 2 {
				backsideTarget = append(backsideTarget, targets[a])
			}
		}
		if len(backsideTarget) > 0 {
			rand.Seed(time.Now().UnixNano())
			return backsideTarget[rand.Intn(len(backsideTarget))]
		} else { //345
			for a := 0; a < len(targets); a++ {
				backsideTarget = append(backsideTarget, targets[a])
			}
			if len(backsideTarget) > 0 {
				rand.Seed(time.Now().UnixNano())
				return backsideTarget[rand.Intn(len(backsideTarget))]
			}
		}
	}
	return nil
}

func sortTargetAliveMonsterBasedSequence(targets []map[string]string) []map[string]string {
	var n = len(targets)
	for i := 1; i < n; i++ {
		j := i
		for j > 0 {
			if targets[j-1]["monster_sequence"] > targets[j]["monster_sequence"] {
				targets[j-1], targets[j] = targets[j], targets[j-1]
			}
			j = j - 1
		}
	}
	return targets
}

func calculateTriggerSuccess(triggerRateX ...int) bool {
	triggerRate := 20
	if len(triggerRateX) > 0 {
		triggerRate = triggerRateX[0]
	}
	rand.Seed(time.Now().UnixNano())
	randomNum := math.Floor(rand.Float64()) * 100
	if randomNum <= float64(triggerRate) {
		return true
	} else {
		return false
	}
}

func blockDefenseOrIgnoreDefense(target map[string]string, damage int, ignoreDefense bool) (map[string]string, int, int) {
	var blockDefense int
	var target_monster_positive = convertBuffArray(target["monster_positive"])
	if helper.ArrayContainsInt(target_monster_positive, 9053) {
		target_monster_positive = helper.ArrayRemoveInt(target_monster_positive, helper.ArrayIntIndex(target_monster_positive, 9053))
		damage = 0
	} else if ignoreDefense {
		blockDefense = 0
	} else {
		blockDefense = helper.StringToInt(target["monster_defend"])
	}
	return target, damage, blockDefense
}

func calcElementDmg(damage int, atkElement int, targetElement int) int {
	if atkElement == 1 {
		if targetElement == 5 {
			damage = int(math.Ceil(float64(damage) * 1.1))
		}
	} else {
		if atkElement == targetElement-1 {
			damage = int(math.Ceil(float64(damage) * 1.1))
		}
	}
	return damage
}
