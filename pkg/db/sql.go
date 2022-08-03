package db

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

type SQLConfig struct {
	Driver    string
	Host      string
	Port      int
	User      string
	Password  string
	DBName    string
	ParseTime bool

	TablePrefix     string
	LogPath         string
	SlowThresholdMs int
}

func ConnectSQL(confs map[string]*SQLConfig) map[string]*gorm.DB {
	dbs := make(map[string]*gorm.DB, len(confs))
	for k := range confs {
		switch confs[k].Driver {
		case "mysql":
			dbs[k] = connectMySQL(*confs[k])
		case "postgres":
			// dbs[k] = connectPostgres(*confs[k])
		default:
			panic("unknown driver: " + confs[k].Driver)
		}
	}

	return dbs
}

func connectMySQL(conf SQLConfig) *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=%t&loc=Local",
		conf.User,
		conf.Password,
		conf.Host,
		conf.Port,
		conf.DBName,
		conf.ParseTime,
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
		log.Println("DB dsn : " + dsn)
		panic("failed to connect database: " + err.Error())
	}

	SetDBConnParams(db, 10, 100, time.Hour)

	SetLogger(db, conf.LogPath, conf.SlowThresholdMs)

	return db
}

/*
func connectPostgres(conf SQLConfig) *gorm.DB {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s dbname=%s sslmode=disable password=%s TimeZone=%s",
		conf.Host,
		conf.Port,
		conf.User,
		conf.DBName,
		conf.Password,
		conf.TimeZone,
	)

	config := &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   conf.TablePrefix,
			SingularTable: true,
		},
		DisableForeignKeyConstraintWhenMigrating: true,
	}

	db, err := gorm.Open(postgres.Open(dsn), config)
	if err != nil {
		panic("failed to connect database: " + err.Error())
	}

	SetDBConnParams(db, 10, 100, time.Hour)

	SetLogger(db, conf.LogPath, conf.SlowThresholdMs)

	return db
}
*/

// 根据自己情况更改连接数
func SetDBConnParams(db *gorm.DB, maxIdleConns int, maxOpenConns int, maxLifeTime time.Duration) {
	sqlDB, _ := db.DB()

	sqlDB.SetMaxIdleConns(maxIdleConns)
	sqlDB.SetMaxOpenConns(maxOpenConns)

	sqlDB.SetConnMaxLifetime(maxLifeTime)
}

func SetLogger(db *gorm.DB, LogPath string, SlowThresholdMs int) {
	var writer io.Writer = os.Stderr
	if LogPath != "" {
		logFile := strings.TrimSuffix(LogPath, "/") + "/sql.log"
		file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			panic("failed to set logger :" + err.Error())
		}
		writer = file
	}

	db.Logger = logger.New(
		log.New(writer, "\r\n", log.Ltime),
		logger.Config{
			SlowThreshold:             time.Millisecond * time.Duration(SlowThresholdMs), // Slow SQL threshold
			LogLevel:                  logger.Info,                                       // Log level
			IgnoreRecordNotFoundError: true,                                              // Ignore ErrRecordNotFound error for logger
			Colorful:                  true,                                              // Disable color
		},
	)
}
