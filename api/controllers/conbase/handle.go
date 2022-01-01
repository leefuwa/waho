package conbase

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"strings"
	"waho/comm"
)

var NotLogin = []string{
	"UserController:Login",
	//"UserController:SendPhoneCode", // 发送手机验证码，而微信小程序必须登录过后才能
}

// 方便用于展示使用
type Handle interface {
	Err(code int, msg ...string) // 如有异常就中断 继续输出
	CheckErr(err error, code int, msg ...string) // 如有异常就中断 继续输出
	CheckErrDef(err error)
	CheckErrBool(err_is_true bool, code int, msg ...string)
	Lis() // 输出List
	Out()
	Msg() // 输出结果
	GetPageIndex() int // 获取当前页
	GetPageSize() int // 获取每页最大条数 默认最大30条
	GetIntNotErr(key string , def int) int
	GetStringTrim(key string, def ...string) string
}


func NewHandle(handle Handle) Handle {
	return handle
}

func (this *BaseController) Handle () Handle {
	return this.handle
}

type requestLog struct {
	Referer string `json:"referer"`
	Url     string `json:"url"`
	Method  string `json:"method"`
	Form    string `json:"form"`
	Header  string `json:"header"`
	Out     string `json:"out"`
}

func (this *BaseController) writeRequestLog()  {
	r := requestLog{}
	r.Referer = this.Ctx.Request.Referer()
	r.Url = this.Ctx.Request.URL.String()
	r.Method = this.Ctx.Request.Method
	form , _ := json.Marshal(this.Ctx.Request.Form)
	r.Form = string(form)
	header , _ := json.Marshal(this.Ctx.Request.Header)
	r.Header = string(header)
	out, _ := json.Marshal(this.Return)
	r.Out = string(out)
	rByte, _ := json.Marshal(r)
	rStr := string(rByte)

	log.Info("[RequestLog]",rStr)
}

func (this *BaseController) Err (code int, msg ...string)  {
	this.Return.Code, this.Return.Message = GetCodeMessage(code, msg...)
	this.Data["json"] = this.Return
	this.writeRequestLog()
	this.ServeJSON()
	this.StopRun()
}

func (this *BaseController) CheckErr (err error, code int, msg ...string)  {
	if err != nil {
		this.Return.Code, this.Return.Message = GetCodeMessage(code, msg...)
		this.Data["json"] = this.Return
		this.writeRequestLog()
		this.ServeJSON()
		this.StopRun()
	}
}

func (this *BaseController) CheckErrDef (err error)  {
	if err != nil {
		this.Handle().CheckErr(err, Fail, err.Error())
	}
}

func (this *BaseController) CheckErrBool (err_is_true bool, code int, msg ...string)  {
	if err_is_true {
		this.Return.Code, this.Return.Message = GetCodeMessage(code, msg...)
		this.Data["json"] = this.Return
		this.writeRequestLog()
		this.ServeJSON()
		this.StopRun()
	}
}

func (this *BaseController) Out () {
	if this.Return.Code == 0 {
		this.Return.Code = Success
	}
	if this.Return.Message == "" {
		this.Return.Message = MessageMap[Success]
	}
	this.Data["json"] = this.Return
	this.writeRequestLog()
	this.ServeJSON()
	this.StopRun()
}

func (this *BaseController) Lis () {
	this.Return.Data = this.List
	this.Data["json"] = this.Return
	this.writeRequestLog()
	this.ServeJSON()
	this.StopRun()
}

func (this *BaseController) Msg () {
	if this.Return.Code == 0 && this.Return.Message == "" {
		this.Return.Code, this.Return.Message = GetCodeMessage(Success)
	}
	this.Data["json"] = this.Return
	this.writeRequestLog()
	this.ServeJSON()
	this.StopRun()
}

// 获取当前页
func (this *BaseController) GetPageIndex() int {
	index, err := this.GetInt("p", 1)
	if err != nil {
		return 1
	}

	return index
}

// 获取每页最大条数
func (this *BaseController) GetPageSize() int {
	pageMax := 30
	size, err := this.GetInt("s")
	if err != nil {
		return pageMax
	}

	size = comm.Min(size, pageMax)

	return size
}

func (this *BaseController) GetIntNotErr(key string , def int) int {
	val, err := this.GetInt(key, def)
	if err != nil {
		return def
	}

	return val
}


func (this *BaseController) GetStringTrim(key string , def ...string) string {
	return strings.Trim(this.GetString(key, def...), " ")
}