package controller

import (
	"app/internal/middleware"
	"app/pkg/response"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
)

var (
	Engine = NewEngine()
)

func NewEngine() *gin.Engine {

	gin.SetMode(gin.ReleaseMode)

	engine := gin.New()
	engine.Use(middleware.CustomRecovery(), middleware.ApiLog())

	// 404
	engine.NoRoute(func(c *gin.Context) {
		response.Fail(c, 1003, "路由不存在")
	})
	// ping
	engine.GET("ping", func(c *gin.Context) {
		response.Success(c, "pong")
	})

	return engine
}

func GetUid(c *gin.Context) int {
	return c.GetInt("uid")
}

// 参数验证，务必在false时return。 若非强校验，使用c.ShouldBind(&req)
func Bind(c *gin.Context, to interface{}) bool {
	if err := c.ShouldBind(to); err != nil {
		response.Fail(c, 1002, err.Error())
		return false
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
