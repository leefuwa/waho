package base

import (
	"github.com/beego/beego/v2/server/web"
	log "github.com/sirupsen/logrus"
	"waho/register"
)

type BeegoRegister struct {
	register.BaseRegister
}

func (beego *BeegoRegister) Init() {
	log.Info("beego init")
}

func (beego *BeegoRegister) Start() {
	log.Info("beego start")
	web.Run()
}

func (beego *BeegoRegister) IsRoutine() bool {
	return false
}
