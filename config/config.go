package config

import (
	"errors"
	"flag"
	"strings"

	"app/pkg"
	"app/pkg/db"
	"app/pkg/logger"
	"app/pkg/token"

	"github.com/spf13/viper"
)

var (
	Run struct {
		ConfigFile string
		IsMigrate  bool
	}
	App struct {
		System string
		Env    string
		Host   string
		Port   int
	}
	DB struct {
		SQL   map[string]*db.SQLConfig
		Redis db.RedisConfig
	}
	Logger logger.Config
	JWT    token.JWTConfig
)

func init() {
	flag.StringVar(&Run.ConfigFile, "c", "./config/config.yaml", "指定配置文件")
	flag.BoolVar(&Run.IsMigrate, "m", false, "是否数据库迁移")
	flag.Parse()
}

func init() {
	// 读取配置
	if err := ReadConfigFile(Run.ConfigFile); err != nil {
		panic(err)
	}

	// 配置目录处理
	if Logger.Filepath != "" {
		Logger.Filepath = strings.TrimSuffix(Logger.Filepath, "/")
		if err := pkg.Filer.CheckDir(Logger.Filepath); err != nil {
			panic(err)
		}
		for k := range DB.SQL {
			DB.SQL[k].LogPath = Logger.Filepath
		}
	}
}

func ReadConfigFile(configFile string) error {
	if !pkg.Filer.IsExist(configFile) {
		return errors.New("config file not found:" + configFile)
	}

	viper.SetConfigFile(configFile)
	ext := configFile[strings.LastIndexByte(configFile, '.')+1:]
	viper.SetConfigType(ext)

	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	viper.UnmarshalKey("app", &App)
	viper.UnmarshalKey("db", &DB)
	viper.UnmarshalKey("logger", &Logger)
	viper.UnmarshalKey("jwt", &JWT)

	return nil
}
