package register

import (
	"errors"
	"gopkg.in/ini.v1"
	"os"
	"path/filepath"
)

//应用程序
type BootApplication struct {}

func init() {
	WorkPath, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	var filePath[]string
	var path string
	path = filepath.Join(WorkPath, "..", "conf", "app.conf")
	_, err = os.Lstat(path)
	if err == nil {
		filePath = append(filePath, path)
	}
	path = filepath.Join(WorkPath, "conf", "app.conf")
	_, err = os.Lstat(path)
	if err == nil {
		filePath = append(filePath, path)
	}

	var cfg *ini.File
	filePathLen := len(filePath)
	if filePathLen == 0 {
		panic(errors.New("请在目录项目下新建 conf/app.conf 配置文件！"))
	} else if filePathLen == 1 {
		cfg, err = ini.Load(filePath[0])
	} else {
		cfg, err = ini.Load(filePath[0],filePath[1])
	}
	if err != nil {
		panic(err)
	}
	appConf := AppConf{}
	appConf.Set(cfg)
	GetConf = appConf
}

func New() *BootApplication {
	return &BootApplication{}
}

func (boot *BootApplication) Start() {
	boot.init()
	boot.start()
}

func (boot *BootApplication) init() {
	for _, v := range GetRegister() {
		v.Init()
	}
}

func (boot *BootApplication) start() {
	for _, v := range GetRegister() {
		if v.IsRoutine() {
			go v.Start()
		} else {
			v.Start()
		}
	}
}
