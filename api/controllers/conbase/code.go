package conbase

const (
	Success = 200
	Fail    = 404
	Error   = 500
	ClientNotToServiceTime = 900
	ReLogin = 901
	ParamsErr = 902
	Err404 = 903
	// 用户模块
	SendCodeExceedNumber = 1001
	UserIsBindPhone = 1002
	OtherBindCurrentUserPhone = 1003
	PhoneCodeExpired = 1004
	PhoneCodeErr = 1005
)

var MessageMap = map[int]string{
	Success:                   "操作成功！",
	Fail:                      "网络繁忙，请稍后重试！",
	Error:                     "服务繁忙，请稍后重试！",
	ReLogin:                   "网络繁忙，请稍后重试！",
	Err404:                    "网络繁忙，请稍后重试！",
	ParamsErr:                 "请按正常补充内容，再次重试！",
	SendCodeExceedNumber:      "当天超过发送次数，请明天再重新发送，如有问题请联系客服人员！",
	UserIsBindPhone:           "用户已绑定手机号码",
	OtherBindCurrentUserPhone: "此手机号码已被绑定",
	PhoneCodeExpired:          "用户验证码已过期，请重新发送",
	PhoneCodeErr:              "验证码输入错误，请重新输入",
}

func GetCodeMessage(code int, msg ...string) (int, string) {
	var message string
	if len(msg) > 0 {
		message = msg[0]
	} else {
		message = MessageMap[code]
	}

	return code, message
}

// 使用统一错误时 为了快速定位到错误
func GetCodeMessageAddMsg(code int, addMsg string) (int, string) {
	message := MessageMap[code] + "  " + addMsg

	return code, message
}