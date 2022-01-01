package base

import (
	"errors"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	vzh "github.com/go-playground/validator/v10/translations/zh"
	log "github.com/sirupsen/logrus"
	"reflect"
	"regexp"
	"waho/comm"
	"waho/register"
)

var validate *validator.Validate
var translator ut.Translator
type ValidatorRegister struct {
	register.BaseRegister
}


func (v *ValidatorRegister) Init() {
	validate = validator.New()
	validate.RegisterValidation("phone", phone)
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		return fld.Tag.Get("comment")
	})
	//创建消息国际化通用翻译器
	cn := zh.New()
	uni := ut.New(cn, cn)
	var found bool
	translator, found = uni.GetTranslator("zh")
	if found {
		err := vzh.RegisterDefaultTranslations(validate, translator)
		if err != nil {
			log.Error(err)
		}
	} else {
		log.Error("translator[中文] 不存在")
	}
}


func ValidateStruct(s interface{}) (err error) {
	err = validate.Struct(s)
	if err != nil {
		_, ok := err.(*validator.InvalidValidationError)
		if ok {
			log.Error("validator 尚未初始化 ， 需要注册validator", err)
		}
		errs, ok := err.(validator.ValidationErrors)
		if ok {
			for _, e := range errs {
				if comm.MapKeyExists(e.Tag(),registerValidationErr) {
					log.Error(registerValidationErr[e.Tag()])
					return errors.New(registerValidationErr[e.Tag()])
				} else {
					log.Error(e.Translate(translator))
					return errors.New(e.Translate(translator))
				}
			}
		}
		return err
	}
	return nil
}

var registerValidationErr = map[string]string{
	"phone": "手机号码填写有误！",
}
func phone(f validator.FieldLevel) bool {
	result, _ := regexp.MatchString(`^(1[3|4|5|8][0-9]\d{8})$`, comm.ToSting(f.Field().Interface()))
	if result {
		return true
	} else {
		return false
	}

}
