package controller

import (
	"app/pkg/response"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type test struct{}

func init() {
	var r = Engine.Group("/test")
	r.GET("/", test{}.Index)
	r.GET("/panic", test{}.Panic)
}

func (t test) Index(c *gin.Context) {
	sleep := c.Query("sleep")
	if sleep != "" {
		i, _ := strconv.Atoi(sleep)
		time.Sleep(time.Duration(i) * time.Second)
	}

	response.Echo(c, nil, nil)
}

func (t test) Panic(c *gin.Context) {
	// error + runtime.Error
	var x map[string]interface{}
	x["a"] = 1
	// error
	panic("测试panic")
}
