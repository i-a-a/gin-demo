package logger

import (
	"os"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

// 对 logrus 加了按照文件大小进行分割的功能

var (
	LoggerHooks []func()
)

type Config struct {
	Filepath string
	Level    string

	// logrus.Formatter
	PrettyPrint     bool
	TimestampFormat string

	// lumberjack.Logger
	MaxSize int // 单个日志体积 MB
	MaxAge  int // 多少天
}

func NewLogrus(conf Config, filename string) *logrus.Logger {
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

	logger.AddHook(hook{})

	return logger
}
