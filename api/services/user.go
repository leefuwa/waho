package services

import (
	"waho/static"
)

var UserService User

func GetUserService() User {
	return UserService
}
type User interface {
	// 判断是否需要登陆
	IsLoginRequired(controller string, action string) bool
	// 小程序登录操作
	LoginWechat(params LoginWechatParams) (token string, err error)
	// 验证token是否有效 并判断帐户是否有异常
	CheckTokenWechat(token string) (isLogin bool, user static.UserCache)
	// 发送手机验证码
	SendPhoneCode(params Phone) error
	// 用户绑定手机号码
	UserBindPhone(phone string, code string, user static.UserCache) (errCode int, errMsg string)
	// 获取验证码记录
	//GetSendCodeCache(userId int) []static.SendPhoneCodeCurrentDayLogCache
}

type LoginWechatParams struct {
	Code string `json:"code" validate:"required"`
}

type Phone struct {
	Phone string `json:"phone" validate:"phone"`
	UserId int `json:"user_id"`
}