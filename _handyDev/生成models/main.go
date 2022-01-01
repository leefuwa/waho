package main

import (
	"fmt"
	"github.com/beego/beego/v2/client/orm"
	"io"
	"os"
	"path/filepath"
	"strings"
	"waho/comm"
	"waho/models"
	"waho/register"
	"waho/register/base"
)

func main()  {
	register.Register(&base.DbRegister{})
	register.New().Start()

	GetTableStructText("user")
}

func GetTableStructText(where string)  {
	type fields struct {
		Field string `orm:"column(Field)" json:"Field"`
		Type string `orm:"column(Type)" json:"Type"`
		Key string `orm:"column(Key)" json:"Key"`
		Comment string `orm:"column(Comment)" json:"Comment"`
	}
	var data []fields

	modelOrm := orm.NewOrm()
	_, _ = modelOrm.Raw("SHOW FULL COLUMNS FROM " + where).QueryRows(&data)
	tableObject := comm.CamelString(strings.TrimLeft(where, models.GetDbPrefix()))


	text := ""
	for _, v := range data{
		fieldType := ""
		if strings.Contains(v.Type, "int") {
			fieldType = "int"
		} else if strings.Contains(v.Type, "char") {
			fieldType = "string"
		} else if strings.Contains(v.Type, "timestamp") {
			fieldType = "string"
		} else if strings.Contains(v.Type, "text") {
			fieldType = "string"
		} else if strings.Contains(v.Type, "float") {
			fieldType = "float64"
		} else {
			fmt.Println("异常错误", v.Field, "---", v.Type)
			os.Exit(0)
		}
		pri := ""
		if strings.ToUpper(v.Key) == "PRI" {
			pri = " orm:\"pk\""
		}
		text += fmt.Sprintf("\t%s %s `json:\"%s\"%s`  //%s\n", comm.CamelString(v.Field), fieldType, v.Field, pri, v.Comment)
		_ = fieldType
	}


	text = `package models

import (
	"github.com/beego/beego/v2/client/orm"
	log "github.com/sirupsen/logrus"
	"waho/comm"
)

type `+tableObject+` struct {
`+text+`	where interface{}
	pageIndex int
	pageSize int
	defaultSoftDelete bool
	txOrm orm.TxOrmer
}

func init()  {
	orm.RegisterModelWithPrefix(GetDbPrefix(), new(`+tableObject+`))
}

func New`+tableObject+`() *`+tableObject+` {
	`+tableObject+` := &`+tableObject+`{defaultSoftDelete: true}
	return `+tableObject+`
}
func (m *`+tableObject+`) DefaultSoftDelete(soft bool) (*`+tableObject+`) {
	m.defaultSoftDelete = soft
	return m
}
func (m *`+tableObject+`) Where(where interface{}) (*`+tableObject+`) {
	m.where = where
	return m
}
func (m *`+tableObject+`) Limit(index int, size int) (*`+tableObject+`) {
	m.pageIndex = index
	m.pageSize = size
	return m
}
func (m *`+tableObject+`) Info(where ...interface{}) (bool, `+tableObject+`) {
	if where != nil{
		m.where = where[0]
	}
	model := NewModels(m)
	exist, _ := model.GetInfo(m.where)
	log.Info(model.GetLastSql())

	return exist, *m
}
func (m *`+tableObject+`) List(where ...interface{}) (int64, []`+tableObject+`) {
	if where != nil{
		m.where = where[0]
	}
	if m.pageIndex == 0 {
		m.pageIndex = 1
		m.pageSize = comm.IntMaxUnsigned()
	}
	var arr []`+tableObject+`
	model := NewModels(&arr)
	count, err := model.GetList(m.where, m.pageIndex, m.pageSize)

	if err != nil {
		log.Error("`+tableObject+` List Error：", model.GetLastSql())
		return 0, nil
	}

	return count, arr
}
func (m *`+tableObject+`) Count(where ...interface{}) (int64, error) {
	if where != nil{
		m.where = where[0]
	}
	var arr []`+tableObject+`
	model := NewModels(&arr)
	count, err := model.GetCount(m.where)

	if err != nil {
		log.Error("`+tableObject+` Count Error：", model.GetLastSql())
		return 0, err
	}

	return count, nil
}
func (m *`+tableObject+`) AddObj(save ...interface{}) (int64, error) {
	model := NewModels(m)
	model.StartTrans(m.txOrm)
	return model.AddObj(save...)
}
func (m *`+tableObject+`) Add(save map[string]interface{}) (int64, error) {
	model := NewModels(m)
	model.StartTrans(m.txOrm)
	return model.Add(save)
}
func (m *`+tableObject+`) UpdateObj(save ...interface{}) (int64, error) {
	model := NewModels(m)
	model.StartTrans(m.txOrm)
	return model.UpdateObj(save...)
}
func (m *`+tableObject+`) Update(save ...map[string]interface{}) (int64, error) {
	model := NewModels(m)
	model.StartTrans(m.txOrm)
	return model.Update(save...)
}
func (m *`+tableObject+`) StartTrans(ormer orm.TxOrmer) *`+tableObject+` {
	m.txOrm = ormer
	return m
}
`



	var filePath string
	WorkPath, _ := os.Getwd()
	var f *os.File


	filePath = filepath.Join(WorkPath, "models",  comm.Lcfirst(comm.CamelString(tableObject))+".go")

	//if comm.CheckFileIsExist(filePath) {
	//	_ = os.Remove(filePath)
	//}

	if comm.CheckFileIsExist(filePath) {
		panic(filePath+"  文件已存在，若需要重新生成 先删除")
	} else {
		f, _ = os.Create(filePath)
		_, _ = io.WriteString(f, text)
	}
}
