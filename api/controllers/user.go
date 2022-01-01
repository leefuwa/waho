package controllers

import (
	"api/controllers/conbase"
	"api/services"
	"waho/comm"
	"waho/core"
)

type UserController struct {
	conbase.BaseController
}

// 用户登录 根据code换取token
// @router /login [*]
func (this *UserController) Login()  {
	params := services.LoginWechatParams{
		Code: this.Handle().GetStringTrim("code"),
	}
	token, err := services.GetUserService().LoginWechat(params)
	this.Handle().CheckErrDef(err)
	this.Return.Data = token
	this.Out()
}

// 发送验证码
// @router /send_phone_code [*]
func (this *UserController) SendPhoneCode()  {
	// 用户每天只能发送10条验证码
	sendPhoneCache := core.GetCore().GetSendCodeCache(this.User.Id)
	this.Handle().CheckErrBool(len(sendPhoneCache) >= 10, conbase.SendCodeExceedNumber)

	params := services.Phone{Phone: this.Handle().GetStringTrim("phone", ""),UserId: this.User.Id}
	err := services.GetUserService().SendPhoneCode(params)
	this.Handle().CheckErrDef(err)

	this.Handle().Msg()
}

// 绑定手机号
// @router /bind/phone [*]
func (this *UserController) UserBindPhone()  {
	code := this.Handle().GetStringTrim("code", "")
	phone := this.Handle().GetStringTrim("phone", "")
	this.Return.Code, this.Return.Message = services.GetUserService().UserBindPhone(phone, code, this.User)
	this.Msg()
}

// 获取用户信息
// @router /user/info [*]
func (this *UserController) GetUserInfo()  {
	type UserEcho struct {
		Username string `json:"username"`  //用户名(非帐号)
		Head string `json:"head"`  //头像
		Sex int `json:"sex"`  //0: 未知, 1: 男, 2: 女
		Birth int `json:"birth"`  //生日
		Phone string `json:"phone"`  //电话
		Email string `json:"email"`  //邮箱
		Money int `json:"money"`  //金额
		Score int `json:"score"`  //积分
	}
	var user UserEcho
	comm.StructFieldsIdenticalCopy(&user, this.User.User)

	this.Return.Data = user
	this.Handle().Out()
}

