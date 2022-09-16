package db

import (
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

var (
	DB *gorm.DB
)

type config struct {
	Driver   string
	Host     string
	Port     int
	User     string
	Password string
	DBName   string

	TablePrefix string
}

func Connect() {
	var conf config
	viper.UnmarshalKey("db", &conf)

	switch conf.Driver {
	case "mysql":
		DB = connectMySQL(conf)
	default:
		panic("unsupported driver: " + conf.Driver)
	}

}
