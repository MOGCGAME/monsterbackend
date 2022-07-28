package helper

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var hmacSampleSecret = []byte("M45IFQ7QhG")

const (
	PASSWORDSEPERATOR = ":::"
	TOKENUIDKEY       = "uid"
	TOKENADMINUIDKEY  = "adminuid"
	USERINFOKEY       = "userinfo"
)

func VerifyJWTString(tokenString string) (jwt.MapClaims, bool) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return hmacSampleSecret, nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// fmt.Println(claims[TOKENUIDKEY])
		return claims, true //fmt.Sprintf("%v", claims[TOKENUIDKEY]), true
	} else {
		fmt.Println(err)
		return nil, false
	}
}

// return uid only
func VerifyJWT(httpRequest *http.Request) (string, bool) {
	//从Request的header获取Authorization
	tokenString := httpRequest.Header.Get("Authorization")
	//检查Authorization没有东西
	if tokenString == "" {
		return "", false
	}
	//验证Authorization的东西
	claims, isValid := VerifyJWTString(tokenString)
	if isValid {
		return fmt.Sprintf("%v", claims[TOKENUIDKEY]), true
	}
	return "", false
}

func NewJWT(uid string, userInfo map[string]interface{}) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		TOKENUIDKEY: uid,
		USERINFOKEY: userInfo,
		// "nbf": time.Date(2015, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
	})
	tokenString, err := token.SignedString(hmacSampleSecret)

	if err != nil {
		fmt.Println(tokenString, err)
	}

	return tokenString
}

func ReadParameters(httpRequest *http.Request) map[string]string {
	reqBody, err := ioutil.ReadAll(httpRequest.Body)
	if err != nil {
		log.Println(err)
	}
	var reqJson map[string]interface{}
	err = json.Unmarshal(reqBody, &reqJson)
	if err != nil {
		log.Println(err)
	}
	reqFormat := map[string]string{}

	for k, v := range reqJson {
		// fmt.Printf("%d of type %T", v, v)
		if float64_v, ok := v.(float64); ok && (k == "uid" || k == "Uid") {

			s := strconv.FormatFloat(float64_v, 'f', -1, 64)
			s = strings.Trim(s, ".")
			reqFormat[k] = s

		} else if int64_v, ok := v.(int64); ok {
			reqFormat[k] = Int64ToString(int64_v)
		} else {
			reqFormat[k] = fmt.Sprintf("%v", v)
		}
	}

	return reqFormat
}

func StringToInt(str string) int {
	i, err := strconv.Atoi(str)
	if err != nil {
		fmt.Printf("%d of type %T", i, i)
		fmt.Printf("StringToInt err %v", err)
	}
	return i
}

func StringToInt64(str string) int64 {
	n, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		fmt.Printf("%d of type %T", n, n)
		fmt.Printf("StringToInt64 err %v", err)
	}
	return n
}

func Int64ToString(n int64) string {
	return strconv.FormatInt(n, 10)
}

func hexaNumberToInteger(hexaString string) int64 {
	// replace 0x or 0X with empty String
	numberStr := strings.Replace(hexaString, "0x", "", -1)
	numberStr = strings.Replace(numberStr, "0X", "", -1)

	if s, err := strconv.ParseInt(numberStr, 16, 32); err == nil {
		fmt.Printf("%T, %v\n", s, s)

		return s
	}

	return 0
}

func ArrayContainsString(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func ArrayContainsInt(i []int, num int) bool {
	for _, v := range i {
		if v == num {
			return true
		}
	}

	return false
}

func ArrayIntIndex(i []int, num int) int {
	for a, v := range i {
		if v == num {
			return a
		}
	}
	return -1
}

func ArrayRemoveInt(i []int, index int) []int {
	return append(i[:index], i[index+1:]...)
}

func GetExpiryDateHanoiTime() string {
	loc, _ := time.LoadLocation("Asia/Ho_Chi_Minh")
	timein := time.Now().In(loc).Add(
		time.Minute * time.Duration(30))
	createdFormat := "20060102150405" // yyyyMMddHHmmss
	return timein.Format(createdFormat)
}

func GetCurrentShanghaiTime() time.Time {
	loc, _ := time.LoadLocation("Asia/Shanghai")
	// fmt.Println(time.Now().Add(time.Hour * time.Duration(8)).Unix())

	timeStr := time.Now().In(loc).Format("2006-01-02 15:04:05")
	timeDB, err := time.ParseInLocation("2006-01-02 15:04:05", timeStr, time.Local)

	if err != nil {
		log.Println(err)
	}

	return timeDB //.Add(time.Hour * time.Duration(8)).Unix()
}

func GetCurrentShanghaiTimeString() string {
	createdFormat := "2006-01-02 15:04:05"
	return GetCurrentShanghaiTime().Format(createdFormat)
	// return strconv.Itoa(int(GetCurrentShanghaiTimeUnix()))
}

func GetCurrentShanghaiDateOnlyString() string {
	createdFormat := "2006-01-02"
	return GetCurrentShanghaiTime().Format(createdFormat)
	// return strconv.Itoa(int(GetCurrentShanghaiTimeUnix()))
}
