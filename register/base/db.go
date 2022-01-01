package base

import (
	"github.com/beego/beego/v2/client/orm"
	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
	"waho/register"
)

type DbRegister struct {
	register.BaseRegister
}

func (db *DbRegister) Init() {
	log.Info("db init")
	// 配置数据库 驱动用beego
	host := register.GetConf.Section("mysql").GetDefault("host", "127.0.0.1")
	port := register.GetConf.Section("mysql").GetDefault("port", "3306")
	name := register.GetConf.Section("mysql").GetDefault("name", "")
	user := register.GetConf.Section("mysql").GetDefault("user", "root")
	pwd := register.GetConf.Section("mysql").GetDefault("pwd", "")
	charset := register.GetConf.Section("mysql").GetDefault("charset", "utf8")
	if name == "" {
		panic("数据库名称[conf/app.conf => name=table]尚未设置")
	}
	dataSource := user + ":" + pwd + "@tcp(" + host + ":" + port + ")/" + name + "?charset=" + charset
	dataSourceLog := "tcp(" + host + ":" + port + ")/" + name + "?charset=" + charset
	log.Info(dataSourceLog)
	_ = orm.RegisterDriver("mysql", orm.DRMySQL)
	err := orm.RegisterDataBase("default", "mysql", dataSource)

	if err != nil {
		log.Error(err)
		panic("数据库连接出现问题")
	}
}