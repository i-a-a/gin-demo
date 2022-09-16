package controller

import (
	"app/internal/middleware"
	"app/pkg/carrot"

	"github.com/gin-gonic/gin"
)

var (
	Engine *gin.Engine
)

func init() {
	gin.SetMode(gin.ReleaseMode)

	Engine = gin.New()
	Engine.Use(middleware.CustomRecovery())

	// 404
	Engine.NoRoute(func(c *gin.Context) {
		carrot.New(c).Echo(nil, carrot.Msg("路由不存在"))
	})
	// ping
	Engine.GET("ping", func(c *gin.Context) {
		carrot.New(c).Echo("pong", nil)
	})
}
