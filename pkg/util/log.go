package util

import (
	"os"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

// 对 logrus 加了按照文件大小进行分割的功能
// 设置默认 logrus 输出位置

type LoggerConfig struct {
	Filepath string
	Level    string

	// logrus.Formatter
	PrettyPrint     bool
	TimestampFormat string

	// lumberjack.Logger
	MaxSize int // 单个日志体积 MB
	MaxAge  int // 多少天
}

// 设置默认 logrus 输出位置
func SetDefultLogrus(conf LoggerConfig, filename string) {
	level, err := logrus.ParseLevel(conf.Level)
	if err != nil {
		panic(err)
	}
	logrus.SetLevel(level)

	logrus.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: conf.TimestampFormat,
		PrettyPrint:     conf.PrettyPrint,
	})

	if conf.Filepath == "" || filename == "" {
		logrus.SetOutput(os.Stdout)
	} else {
		logrus.SetOutput(&lumberjack.Logger{
			Filename:   conf.Filepath + "/" + filename,
			MaxSize:    conf.MaxSize, // MB
			MaxAge:     conf.MaxAge,  // Day
			MaxBackups: 10,
			LocalTime:  true,
			Compress:   true,
		})
	}
}

func NewLogrus(conf LoggerConfig, filename string) *logrus.Logger {
	logger := logrus.New()

	level, err := logrus.ParseLevel(conf.Level)
	if err != nil {
		panic(err)
	}
	logger.SetLevel(level)

	logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: conf.TimestampFormat,
		PrettyPrint:     conf.PrettyPrint,
	})

	if conf.Filepath == "" || filename == "" {
		logger.SetOutput(os.Stdout)
	} else {
		logger.SetOutput(&lumberjack.Logger{
			Filename:   conf.Filepath + "/" + filename,
			MaxSize:    conf.MaxSize, // MB
			MaxAge:     conf.MaxAge,  // Day
			MaxBackups: 10,
			LocalTime:  true,
			Compress:   true,
		})
	}

	return logger
}
