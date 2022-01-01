package models

import (
	"github.com/beego/beego/v2/client/orm"
	log "github.com/sirupsen/logrus"
	"waho/comm"
)

type User struct {
	Id int `json:"id" orm:"pk"`  //
	Username string `json:"username"`  //用户名(非帐号)
	Head string `json:"head"`  //头像
	Sex int `json:"sex"`  //0: 未知, 1: 男, 2: 女
	Birth int `json:"birth"`  //生日
	Phone string `json:"phone"`  //电话
	Email string `json:"email"`  //邮箱
	Password string `json:"password"`  //
	Openid string `json:"openid"`  //
	Status int `json:"status"`  //状态，1：正常，2：禁用
	Money int `json:"money"`  //金额
	Score int `json:"score"`  //积分
	CreTime int `json:"cre_time"`  //创建时间
	UpdTime int `json:"upd_time"`  //更新时间
	DelTime int `json:"del_time"`  //删除时间
	IsDelete int `json:"is_delete"`  //是否已删除
	where interface{}
	pageIndex int
	pageSize int
	defaultSoftDelete bool
	txOrm orm.TxOrmer
}

func init()  {
	orm.RegisterModelWithPrefix(GetDbPrefix(), new(User))
}

func NewUser() *User {
	User := &User{defaultSoftDelete: true}
	return User
}
func (m *User) DefaultSoftDelete(soft bool) (*User) {
	m.defaultSoftDelete = soft
	return m
}
func (m *User) Where(where interface{}) (*User) {
	m.where = where
	return m
}
func (m *User) Limit(index int, size int) (*User) {
	m.pageIndex = index
	m.pageSize = size
	return m
}
func (m *User) Info(where ...interface{}) (bool, User) {
	if where != nil{
		m.where = where[0]
	}
	model := NewModels(m)
	exist, _ := model.GetInfo(m.where)
	log.Info(model.GetLastSql())

	return exist, *m
}
func (m *User) List(where ...interface{}) (int64, []User) {
	if where != nil{
		m.where = where[0]
	}
	if m.pageIndex == 0 {
		m.pageIndex = 1
		m.pageSize = comm.IntMaxUnsigned()
	}
	var arr []User
	model := NewModels(&arr)
	count, err := model.GetList(m.where, m.pageIndex, m.pageSize)

	if err != nil {
		log.Error("User List Error：", model.GetLastSql())
		return 0, nil
	}

	return count, arr
}
func (m *User) Count(where ...interface{}) (int64, error) {
	if where != nil{
		m.where = where[0]
	}
	var arr []User
	model := NewModels(&arr)
	count, err := model.GetCount(m.where)

	if err != nil {
		log.Error("User Count Error：", model.GetLastSql())
		return 0, err
	}

	return count, nil
}
func (m *User) AddObj(save ...interface{}) (int64, error) {
	model := NewModels(m)
	model.StartTrans(m.txOrm)
	return model.AddObj(save...)
}
func (m *User) Add(save map[string]interface{}) (int64, error) {
	model := NewModels(m)
	model.StartTrans(m.txOrm)
	return model.Add(save)
}
func (m *User) UpdateObj(save ...interface{}) (int64, error) {
	model := NewModels(m)
	model.StartTrans(m.txOrm)
	return model.UpdateObj(save...)
}
func (m *User) Update(save ...map[string]interface{}) (int64, error) {
	model := NewModels(m)
	model.StartTrans(m.txOrm)
	return model.Update(save...)
}
func (m *User) StartTrans(ormer orm.TxOrmer) *User {
	m.txOrm = ormer
	return m
}
