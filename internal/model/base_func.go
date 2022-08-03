package model

import (
	"app/internal/common"
)

func AutoMigrate() {
	common.DB.AutoMigrate(
		&User{},
	)
}
