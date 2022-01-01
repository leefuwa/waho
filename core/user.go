package core

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"time"
	"waho/comm"
	"waho/register/base"
	"waho/static"
)


//func (c *c) SetSendCodeCache(cacheList []static.SendPhoneCodeCurrentDayLogCache)  {
//	if cacheList == nil || len(cacheList) == 0 {
//		return
//	}
//
//	base.HSet(static.SendPhoneCodeCurrentDayLogKey, comm.ToJsonNotErr(cacheList))
//}

func (c *c) GetSendCodeCache(userId int) []static.SendPhoneCodeCurrentDayLogCache {
	str, _ := base.RdbHGet(static.SendPhoneCodeCurrentDayLogKey, comm.ToSting(userId))
	var sendPhoneCache []static.SendPhoneCodeCurrentDayLogCache
	_ = json.Unmarshal([]byte(str), &sendPhoneCache)

	return sendPhoneCache
}

func (c *c) GetSendCodeNewCache(userId int) static.SendPhoneCodeCurrentDayLogCache {
	str, _ := base.RdbHGet(static.SendPhoneCodeCurrentDayLogKey, comm.ToSting(userId))
	var sendPhoneCache []static.SendPhoneCodeCurrentDayLogCache
	_ = json.Unmarshal([]byte(str), &sendPhoneCache)
	sendPhoneNewCache := sendPhoneCache[len(sendPhoneCache)-1]

	return sendPhoneNewCache
}

func (c *c) GetSetUserCommKey(key string) string {
	md5str := key + "-wechat-comm-key-" + comm.ToSting(time.Now().Second())
	m := md5.Sum([]byte(md5str))
	commKey := hex.EncodeToString(m[:])
	return commKey
}