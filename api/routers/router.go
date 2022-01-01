package routers

import (
	"api/controllers"
	"api/controllers/conbase"
	"github.com/beego/beego/v2/server/web"
)

func init() {
	web.ErrorController(&conbase.ErrorController{})
	web.Include(&controllers.UserController{})
}
