package main

import (
	_ "api/core"
	_ "api/routers"
	"waho/register"
	"waho/register/base"
)

func init() {
	register.Register(&base.RedisRegister{})
	register.Register(&base.ValidatorRegister{})
	register.Register(&base.LogRegister{})
	register.Register(&base.DbRegister{})
	register.Register(&base.BeegoRegister{})
}
