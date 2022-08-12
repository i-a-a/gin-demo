package main

import (
	"app/config"
	"app/internal/common"
	"app/internal/controller"
	"app/internal/middleware"
	"app/pkg"
	"app/pkg/db"
	"app/pkg/logger"
	"fmt"
	"strconv"

	"github.com/sirupsen/logrus"
)

func main() {
	// 设置环境
	common.SetEnv(config.App.Env)

	// 初始化数据库
	db.SQL.Connect()
	// 初始化Redis
	db.Redis.Connect(false)

	// 日志钩子：机器人
	// logger.AddErrorHooks(func(msg string) {
	// 	pkg.Robot.SendText(config.App.Env, config.App.System, msg)
	// })
	// 默认logrus输出到文件
	logrus.SetOutput(logger.NewLogrus("app.log").Out)
	// 接口日志
	middleware.ResponseLogger = logger.NewLogrus("/response.log").Out
	// 第三方请求日志
	pkg.HttpLogger = logger.NewLogrus("/http.log").Out

	// 启动HTTP服务
	fmt.Println("启动HTTP服务: " + config.App.Host + ":" + strconv.Itoa(config.App.Port))
	if err := controller.Engine.Run(":" + strconv.Itoa(config.App.Port)); err != nil {
		panic(err)
	}
}
