package db

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var (
	// 慢查询阈值
	slowThresholdMs time.Duration = 300 * time.Millisecond
	// 连接参数
	maxIdleConns = 10
	maxOpenConns = 0
	maxLifeTime  = time.Hour
)

func connectMySQL(conf config) *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=False&loc=Local",
		conf.User,
		conf.Password,
		conf.Host,
		conf.Port,
		conf.DBName,
	)

	config := &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   conf.TablePrefix,
			SingularTable: true,
		},
		DisableForeignKeyConstraintWhenMigrating: true,
	}

	db, err := gorm.Open(mysql.Open(dsn), config)
	if err != nil {
		fmt.Println("DB dsn : " + dsn)
		panic("failed to connect database: " + err.Error())
	}

	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(maxIdleConns)
	sqlDB.SetMaxOpenConns(maxOpenConns)
	sqlDB.SetConnMaxLifetime(maxLifeTime)

	db.Logger = logger.New(
		log.New(os.Stdout, "\r\n", log.Ltime),
		logger.Config{
			SlowThreshold:             slowThresholdMs, // Slow SQL threshold
			LogLevel:                  logger.Info,     // Log level
			IgnoreRecordNotFoundError: true,            // Ignore ErrRecordNotFound error for logger
			Colorful:                  true,            // Disable color
		},
	)

	return db
}
