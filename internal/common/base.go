package common

import (
	"app/pkg/db"

	"gorm.io/gorm"
)

var (
	IsLocal bool
	IsDev   bool
	IsProd  bool
)

var (
	DB    *gorm.DB
	Redis *db.Redis
)
