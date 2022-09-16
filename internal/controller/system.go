package controller

import (
	"app/config"
	"app/pkg/carrot"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func init() {
	var r = Engine.Group("/system")
	r.POST("/debug", system{}.ChangeDebugMode) // 切换debug模式
}

type system struct{}

func (system) ChangeDebugMode(c *gin.Context) {
	app := carrot.New(c)
	var req struct {
		Key string `json:"key" form:"key"`
	}
	c.ShouldBind(&req)
	if req.Key != config.App.Key {
		app.Echo(nil, carrot.Msg("key不正确"))
		return
	}

	config.IsDebug = !config.IsDebug
	if config.IsDebug {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}

	app.Echo("ok", nil)
}
