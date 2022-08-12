package model

import (
	"app/pkg/db"

	"gorm.io/gorm"
)

func AutoMigrate() {

}

func DB() *gorm.DB {
	return db.GetDB("default")
}
