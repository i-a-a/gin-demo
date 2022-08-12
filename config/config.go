package config

import (
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
)

func init() {
	flag.StringVar(&Run.ConfigFile, "c", "./config/config.yaml", "指定配置文件")
	flag.BoolVar(&Run.IsMigrate, "m", false, "是否数据库迁移")
	flag.Parse()
}

func init() {
	if !pkg.Filer.IsExist(Run.ConfigFile) {
		panic("config file not found:" + Run.ConfigFile)
	}

	// 解析配置
	viper.SetConfigFile(Run.ConfigFile)
	ext := Run.ConfigFile[strings.LastIndexByte(Run.ConfigFile, '.')+1:]
	viper.SetConfigType(ext)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	// 配置赋值
	viper.UnmarshalKey("app", &App)
	viper.UnmarshalKey("sql", &db.SQL.Configs)
	viper.UnmarshalKey("redis", &db.Redis.Config)
	viper.UnmarshalKey("logger", &logger.Config)
	viper.UnmarshalKey("jwt", &token.JwtConfig)
	viper.UnmarshalKey("robotApi", &pkg.RobotApi)

	// 日志目录检查及设置
	if logger.Config.Filepath != "" {
		logDir := strings.TrimSuffix(logger.Config.Filepath, "/")
		if err := pkg.Filer.CheckDir(logDir); err != nil {
			panic(err)
		}
		for i := range db.SQL.Configs {
			db.SQL.Configs[i].LogPath = logDir
		}
	}

}
