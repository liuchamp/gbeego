package services

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	"math/rand"
	"mul/models"
	"mul/utils"
	"time"
)

/*获取新的短信验证码
 */
func GetNewSmsVerCode() *models.SMSVerCode {

	redisC := utils.GetCacheClient()
breakHere:
	smscode := generate6SmsVerCode(0)

	for {
		v, e := redisC.Do("get", smscode.Code).String()
		if e != nil {
			logs.Warning(e)
		}
		if v == "" {
			break
		}
		smscode = generate6SmsVerCode(0)
	}
	b, e := json.Marshal(smscode)
	if e != nil {
		panic(e)
	}
	v, es := redisC.Do("setex", smscode.Code, smscode.TimeOut, b).String()
	if es != nil {
		panic(es)
	}
	logs.Debug(v)
	if v != "OK" {
		goto breakHere
	}

	return smscode

}

// 当redis存在这个码， 返回true 并且删除。
func CheckSmsCodeExsit(code string) bool {
	redisC := utils.GetCacheClient()
	v, es := redisC.Do("get", code).String()
	if es != nil {
		logs.Debug("Check SMS code meeting error -> ", es)
	}
	vd := models.SMSVerCode{}
	err := json.Unmarshal([]byte(v), &vd)
	if err != nil {
		logs.Warning("parse data meeting erro ->", err)
		return false
	}
	if vd.Code == code {
		defer redisC.Del(code)
		return true
	}
	return false

}

//** 生成6位的sms验证码
func generate6SmsVerCode(timeOut int64) *models.SMSVerCode {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	vcode := fmt.Sprintf("%06v", rnd.Int31n(1000000))
	smsVerCode := new(models.SMSVerCode)
	smsVerCode.Code = vcode
	smsVerCode.Time = time.Now().Unix()
	// 当设置的超时时间小于10时，那么需要将其改为超时时间
	if timeOut < 10 {
		timeOut = 60
	}
	smsVerCode.TimeOut = timeOut
	return smsVerCode
}

func init() {

}
