package conbase

import (
	"api/services"
	"github.com/beego/beego/v2/server/web"
	"time"
	"waho/static"
)

const (
	VersionWeb = "web"
	VersionWechat = "wechat"
	VersionApp = "app"
)


type List struct {
	Count int64 `json:"count"`
	Data interface{} `json:"data"`
}

type Return struct {
	Code int `json:"code"`
	Message string `json:"message"`
	Data interface{} `json:"data"`
}

type BaseController struct {
	web.Controller
	Time    int64
	Version string
	User    static.UserCache
	Return  Return
	List    List
	handle  Handle
}


func (this *BaseController) Prepare() {
	this.Time = time.Now().Unix()
	this.setVersion()
	this.handle = NewHandle(this)
	if services.GetUserService().IsLoginRequired(this.GetControllerAndAction()) {
		var isLogin bool
		isLogin, this.User = services.GetUserService().CheckTokenWechat(this.Handle().GetStringTrim("token", ""))
		if !isLogin {
			this.Handle().Err(ReLogin)
		}
	}
}


// 设置属于哪版本
func (this *BaseController) setVersion() {
	this.Version = VersionWechat

	if false {
		this.Handle().Err(Fail)
	}
}