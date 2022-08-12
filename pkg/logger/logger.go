package logger

import (
	"os"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

// 对 logrus 加了按照文件大小进行分割的功能
var (
	Config struct {
		Filepath string
		Level    string

		// logrus.Formatter
		PrettyPrint     bool
		TimestampFormat string

		// lumberjack.Logger
		MaxSize int // 单个日志体积 MB
		MaxAge  int // 多少天
	}
	once sync.Once
)

func NewLogrus(filename string) *logrus.Logger {
	// 统一的设置
	once.Do(func() {
		level, err := logrus.ParseLevel(Config.Level)
		if err != nil {
			panic(err)
		}
		logrus.SetLevel(level)
		logrus.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: Config.TimestampFormat,
			PrettyPrint:     Config.PrettyPrint,
		})
		// 钩子
		logrus.AddHook(hook{})

		Config.Filepath = strings.TrimSuffix(Config.Filepath, "/")
	})

	logger := logrus.New()

	// 默认输出到终端，输出到文件则进行切割
	if Config.Filepath == "" || filename == "" {
		logger.SetOutput(os.Stdout)
	} else {
		logger.SetOutput(&lumberjack.Logger{
			Filename:   Config.Filepath + "/" + filename,
			MaxSize:    Config.MaxSize, // MB
			MaxAge:     Config.MaxAge,  // Day
			MaxBackups: 10,
			LocalTime:  true,
			Compress:   true,
		})
	}

	return logger
}
