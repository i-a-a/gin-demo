package internal

import (
	"app/config"
	"app/internal/common"
	"app/internal/middleware"
	"app/internal/model"
	"app/pkg"
	"app/pkg/db"
	"app/pkg/token"
	"fmt"

	"app/pkg/logger"

	"github.com/sirupsen/logrus"
)

func InitService() {

	// 初始化数据库 （暂未增加gorm的主从）
	dbs := db.ConnectSQL(config.DB.SQL)
	common.DB = dbs["default"]

	// 初始化redis
	if config.DB.Redis.Host != "" {
		common.Redis = db.ConnRedis(config.DB.Redis)
	}

	// 初始化组件日志
	if config.Logger.Filepath != "" {
		fmt.Println("日志将打印到：" + config.Logger.Filepath)
		// 接口日志
		middleware.ResponseLog = logger.NewLogrus(config.Logger, "/response.log").Out
		// 第三方HTTP请求日志
		pkg.HttpLogger = logger.NewLogrus(config.Logger, "/http.log").Out
	}

	// 默认logrus日志
	defaultLogger := logger.NewLogrus(config.Logger, "app.log")
	logrus.SetOutput(defaultLogger.Out)
	logrus.SetFormatter(defaultLogger.Formatter)
	logrus.SetLevel(defaultLogger.Level)

	// 初始化JWT设置
	token.SetConf(config.JWT)

	// 设置环境
	switch config.App.Env {
	case "local":
		common.IsLocal = true
	case "dev":
		common.IsDev = true
	case "prod":
		common.IsProd = true
	}

	// 是否执行数据库迁移
	if config.Run.IsMigrate {
		model.AutoMigrate()
		return
	}

}
