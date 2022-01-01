package core

import (
	"api/controllers/conbase"
	"api/services"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/rand"
	"sync"
	"time"
	"waho/comm"
	core2 "waho/core"
	"waho/models"
	"waho/register/base"
	"waho/static"
)

var once sync.Once

func init() {
	once.Do(func() {
		services.UserService = new(user)
	})
}

type user struct {
	core
}

func (s *user) IsLoginRequired(controller string, action string) bool {
	flag, _ := comm.InArray(controller+":"+action, conbase.NotLogin)
	if flag {
		return false
	} else {
		return true
	}
}
func (s *user) LoginWechat(params services.LoginWechatParams) (token string, err error) {
	err = base.ValidateStruct(&params)

	if err != nil {
		return "", err
	}
	openid, err := s.getOpenId(params.Code)
	if err != nil {
		return "", err
	}
	where := models.MakeWhereMap()
	where["openid"] = openid
	exist, info := models.NewUser().Where(where).Info()

	// 如果不存在则先创建一条
	if !exist {
		var add models.User
		add.Openid = openid
		add.CreTime = comm.ToInt(time.Now().Unix())
		id, err := models.NewUser().AddObj(&add)
		if err != nil {
			return "", err
		}

		info = add
		info.Id = comm.ToInt(id)
	} else {
		// 存在则删除原有token
		userStr, _ := base.RdbHGet(static.UserKey, comm.ToSting(info.Id))
		var userCache static.UserCache
		err := json.Unmarshal([]byte(userStr), &userCache)
		if err == nil {
			base.RdbHDel(static.UserTokenKey, userCache.Token)
		}
	}

	token = s.createTokenWechat(info)
	var userCache static.UserCache
	userCache.User = info
	userCache.Token = token
	userCache.CommKey = core2.GetCore().GetSetUserCommKey(token)

	jsonUserStatic, _ := json.Marshal(userCache)
	// 把user id放到缓存里面去
	base.RdbHSet(static.UserKey, info.Id, string(jsonUserStatic))
	// 把token保存在缓存里
	base.RdbHSet(static.UserTokenKey, token, info.Id)

	return token, nil
}

func (s *user) getOpenId(code string) (openid string, err error) {
	// 根据code获取openID
	openid = "openid_test_001"
	return
}

func (s *user) createTokenWechat(userInfo models.User) (token string) {
	md5str := comm.ToSting(userInfo.Id) + "-wechat-" + comm.ToSting(time.Now().Second())
	m := md5.Sum([]byte(md5str))
	token = hex.EncodeToString(m[:])
	return
}

func (s *user) CheckTokenWechat(token string) (bool, static.UserCache) {
	if token == "" {
		return false, static.UserCache{}
	}
	// Token不存在
	str, err := base.RdbHGet(static.UserTokenKey, token)
	if str == "" || err != nil {
		return false, static.UserCache{}
	}
	// 用户不存在
	userStr, err := base.RdbHGet(static.UserKey, str)
	if userStr == "" || err != nil {
		return false, static.UserCache{}
	}

	var userCache static.UserCache
	_ = json.Unmarshal([]byte(userStr), &userCache)
	if token != userCache.Token || userCache.Id <= 0 {
		return false, static.UserCache{}
	}

	return true, userCache
}

func (s *user) SendPhoneCode(params services.Phone) error {
	err := base.ValidateStruct(&params)

	if err != nil {
		return err
	}

	// 生成code
	code := fmt.Sprintf("%06v", rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(1000000))

	cache := static.SendPhoneCodeCurrentDayLogCache{
		Phone: params.Phone,
		Code: code,
		UserId: params.UserId,
	}
	// 加入任务 发送验证码
	cacheByte, _ := json.Marshal(cache)
	err = base.Publish(static.SendPhoneCodeQueueName, string(cacheByte))
	if err != nil {
		return err
	}

	return nil
}

func (s *user) UserBindPhone(phone string, code string, user static.UserCache) (int, string) {
	if phone == "" || code == "" {
		return conbase.GetCodeMessageAddMsg(conbase.Fail, "[1002-000]")
	}

	// 用户已绑定手机号
	if user.Phone != "" {
		return conbase.GetCodeMessage(conbase.UserIsBindPhone)
	}

	where := models.MakeWhereMap()
	where["phone"] = phone
	// 其他用户绑定了此手机号码
	count, err := models.NewUser().Where(where).Count()
	if err != nil {
		return conbase.GetCodeMessageAddMsg(conbase.Fail, "[1002-001]")
	}
	if count > 0 {
		return conbase.GetCodeMessage(conbase.OtherBindCurrentUserPhone)
	}

	// 验证手机号码和验证码是否一致
	codeCache := core2.GetCore().GetSendCodeNewCache(user.Id)
	// 当前时间
	currentTime := time.Now().Unix()
	if currentTime - 600 > codeCache.SendTime {
		return conbase.GetCodeMessage(conbase.PhoneCodeExpired)
	}
	// 验证码错误
	if code != codeCache.Code || phone != codeCache.Phone {
		return conbase.GetCodeMessage(conbase.PhoneCodeErr)
	}

	// 更新数据库
	updSave := map[string]interface{}{
		"phone": phone,
		"upd_time": currentTime,
	}
	updWhere := map[string]interface{}{
		"id": user.Id,
	}
	count, err = models.NewUser().Update(updSave, updWhere)
	if err != nil || count == 0 {
		return conbase.GetCodeMessageAddMsg(conbase.Fail, "[1002-002]")
	}

	// 更新缓存
	user.Phone = phone
	user.UpdTime = comm.ToInt(currentTime)
	jsonUserStatic, _ := json.Marshal(user)
	base.RdbHSet(static.UserKey, user.Id, string(jsonUserStatic))

	// TODO 放进队列 再update一次 防止失败

	return conbase.GetCodeMessage(conbase.Success)
}