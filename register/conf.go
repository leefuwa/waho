package register

import (
	"gopkg.in/ini.v1"
)

var GetConf AppConf
// 注册服务的配置信息，贯穿整个生命周期
type AppConf map[string]interface{}

const (
	appConfType = "ini"
)

func (conf AppConf) getConfKey() string {
	return "confKey"
}
func (conf AppConf) getConfKeyVal() string {
	return conf[conf.getConfKey()].(string)
}
func (conf AppConf) Section(key string) AppConf {
	conf[conf.getConfKey()] = key
	return conf
}
func (conf AppConf) HasKey(key string) bool {
	val := conf[appConfType].(*ini.File).Section(conf.getConfKeyVal()).HasKey(key)
	conf[conf.getConfKey()] = ""
	return val
}
func (conf AppConf) Get(key string) string {
	val := conf[appConfType].(*ini.File).Section(conf.getConfKeyVal()).Key(key).String()
	conf[conf.getConfKey()] = ""
	return val
}
func (conf AppConf) GetInt(key string) (int, error) {
	val, err := conf[appConfType].(*ini.File).Section(conf.getConfKeyVal()).Key(key).Int()
	conf[conf.getConfKey()] = ""
	return val,err
}
func (conf AppConf) GetDefault(key string, def string) string {
	if !conf[appConfType].(*ini.File).Section(conf.getConfKeyVal()).HasKey(key) {
		return def
	}

	return conf.Get(key)
}
func (conf AppConf) GetIntDefault(key string, def int) int {
	if !conf[appConfType].(*ini.File).Section(conf.getConfKeyVal()).HasKey(key) {
		return def
	}
	val, err := conf.GetInt(key)
	if err != nil {
		return def
	}

	return val
}

func (conf AppConf) Set(cfg *ini.File) {
	conf[appConfType] = cfg
	conf[conf.getConfKey()] = ""
}
