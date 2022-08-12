package service

import (
	"app/pkg/db"

	"gorm.io/gorm"
)

func DB() *gorm.DB {
	return db.GetDB("default")
}
