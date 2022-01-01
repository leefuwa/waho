package core

import (
	"sync"
	"waho/static"
)

var once sync.Once
var core Core

func init()  {
	once.Do(func() {
		core = new(c)
	})
}

func GetCore() Core {
	return core
}

type Core interface{
	// 获取用户发送短信验证码的列表
	GetSendCodeCache(userId int) []static.SendPhoneCodeCurrentDayLogCache
	// 获取用户发送最新的短信验证码
	GetSendCodeNewCache(userId int) static.SendPhoneCodeCurrentDayLogCache
	// 获取设置用户加密密钥key
	GetSetUserCommKey(key string) string
	// 加密价格
}

type c struct {}

