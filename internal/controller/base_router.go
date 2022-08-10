package controller

import (
	"app/internal/middleware"
	"app/pkg/response"

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
		response.Echo(c, nil, response.Msg("路由不存在"))
	})
	// ping
	engine.GET("ping", func(c *gin.Context) {
		response.Echo(c, "pong", nil)
	})

	return engine
}
