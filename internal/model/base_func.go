package model

import (
	"app/config"
	"app/pkg/db"

	"gorm.io/gorm"
)

// 不要直接使用db.DB
func DB() *gorm.DB {
	if config.IsDebug {
		return db.DB.Debug()
	}
	return db.DB
}

// 很多查询语句是非常简单和确定的，就不必要打印。  而非查询是一定要使用 .Debug()记录的
func DebugDB() *gorm.DB {
	return db.DB.Debug()
}

func AutoMigrate() {
	DebugDB().AutoMigrate(
		&User{}, &Captcha{},
	)
}
