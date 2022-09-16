package carrot

import (
	"reflect"
	"strings"

	"github.com/go-playground/locales/zh"
	translator "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
)

var (
	validate *validator.Validate
	trans    translator.Translator
)

func init() {
	uni := translator.New(zh.New())
	trans, _ = uni.GetTranslator("zh")
	validate = validator.New()
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		return strings.ToLower(fld.Name)
	})
	err := zh_translations.RegisterDefaultTranslations(validate, trans)
	if err != nil {
		panic("翻译器初始化失败 : " + err.Error())
	}

}

// 自定义校验规则
func AddValidate(name string, fn validator.Func) error {
	return validate.RegisterValidation(name, fn)
}
