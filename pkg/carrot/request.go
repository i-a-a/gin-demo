package carrot

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// 参数验证，务必在false时return。
func (c *Carrot) Bind(to interface{}) bool {
	// 读取body
	if c.Switch.RequestLog {
		c.requestData = copyRequestBody(c.ctx)
	}

	c.ctx.ShouldBind(to) // 这里也返回error，但是这个error没法翻译。 例如传的类型不对，忽略问题不大，如果是必要的，后续会认为其是零值产生错误

	err := validate.Struct(to)
	if err != nil {
		errs := err.(validator.ValidationErrors)
		for _, e := range errs {
			s := e.Translate(trans)
			c.Echo(nil, errors.New(s))
			return false
		}
	}

	return true
}

func (c *Carrot) BindAndTrim(to interface{}) bool {
	if !c.Bind(to) {
		return false
	}
	return true
}

func trimString(obj interface{}) {
	elem := reflect.Indirect(reflect.ValueOf(obj))
	for i := 0; i < elem.NumField(); i++ {
		field := elem.Field(i)
		switch field.Kind() {
		case reflect.String:
			if field.String() != "" {
				field.SetString(strings.TrimSpace(field.String()))
			}
		case reflect.Struct:
			trimString(field.Addr().Interface())
		}
	}
}

// 拷贝请求body数据
func copyRequestBody(ctx *gin.Context) []byte {
	if ctx.Request.Method == http.MethodGet {
		return nil
	}

	bodyBytes, _ := ctx.GetRawData()
	ctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	return bodyBytes
}
