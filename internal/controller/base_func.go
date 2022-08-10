package controller

import (
	"app/pkg/response"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/locales/zh"
	translator "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
)

var (
	validate *validator.Validate
	trans    translator.Translator
)

// 自定义中文验证器
func init() {
	uni := translator.New(zh.New())
	trans, _ = uni.GetTranslator("zh")
	validate = validator.New()
	//注册一个函数，获取struct tag里自定义的label作为字段名
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		return "「" + fld.Name + "」"
	})
	//注册翻译器
	err := zh_translations.RegisterDefaultTranslations(validate, trans)
	if err != nil {
		panic("翻译器初始化失败 : " + err.Error())
	}
}

func GetUid(c *gin.Context) int {
	return c.GetInt("uid")
}

// 参数验证，务必在false时return。 若非强校验，使用c.ShouldBind(&req)
func Bind(c *gin.Context, to interface{}) bool {
	c.ShouldBind(to)
	err := validate.Struct(to)
	if err != nil {
		errs := err.(validator.ValidationErrors)
		for _, e := range errs {
			s := e.Translate(trans)
			response.Echo(c, nil, response.Msg(s))
			return false
		}
	}

	return true
}

func BindAndTrim(c *gin.Context, to interface{}) bool {
	if !Bind(c, to) {
		return false
	}
	TrimString(to)
	return true
}

func TrimString(obj interface{}) {
	elem := reflect.Indirect(reflect.ValueOf(obj))
	for i := 0; i < elem.NumField(); i++ {
		field := elem.Field(i)
		switch field.Kind() {
		case reflect.String:
			if field.String() != "" {
				field.SetString(strings.TrimSpace(field.String()))
			}
		case reflect.Struct:
			TrimString(field.Addr().Interface())
		}
	}
}
