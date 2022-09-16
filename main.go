package main

import (
	"app/config"
	"app/internal/controller"
	"app/internal/model"
	"app/pkg/db"
	"net/http"
	"time"
)

func main() {

	// 连接数据库
	db.Connect()
	// 数据迁移
	if config.Flag.IsMigrate {
		model.AutoMigrate()
		return
	}

	// 启动HTTP服务
	s := &http.Server{
		Addr:           "0.0.0.0:" + config.App.Port,
		Handler:        controller.Engine,
		WriteTimeout:   6 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1M
	}
	if err := s.ListenAndServe(); err != nil {
		panic(err)
	}
}
