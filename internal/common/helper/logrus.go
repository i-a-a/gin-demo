package helper

import (
	"app/config"
	"app/pkg/util"
	"fmt"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
)

func init() {
	if config.IsDebug {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}
	if config.App.LogDir != "" {
		logrus.SetOutput(util.NewWriter(config.App.LogDir + "/app.log"))
		// db.SetLogWriter(util.NewWriter(config.App.LogDir + "/sql.log"))
	}
	logrus.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		PrettyPrint:     false,
	})
	logrus.AddHook(&Hook{})
}

// logrus日志钩子，在Error时自动发送机器人通知
// TODO: 优化机器人通知，相同错误短时间内不重复提醒
type Hook struct{}

func (h *Hook) Fire(entry *logrus.Entry) error {
	if entry.Message == "" {
		return nil
	}
	stack := getErrorStack(entry.Message, "internal")
	entry.Data["stack"] = stack

	const format = "环境: %s\n系统: %s\n错误: %v"
	content := fmt.Sprintf(format, config.App.Env, config.App.System, strings.Join(stack, "\n"))
	util.Robot.SendText(content)
	return nil
}

func (h *Hook) Levels() []logrus.Level {
	return []logrus.Level{logrus.ErrorLevel}
}

// 获取错误文件和行号。 去除go自带函数和外部包。
func getErrorStack(errString string, splitDirName string) []string {
	var pcs [32]uintptr
	n := runtime.Callers(3, pcs[:])
	idx := 0
	recorder := []string{errString}

	for i, pc := range pcs[:n] {
		fn := runtime.FuncForPC(pc)
		filepath, line := fn.FileLine(pc)

		if strings.Contains(filepath, splitDirName) {
			if idx == 0 {
				idx = strings.Index(filepath, splitDirName)
			}
			recorder = append(recorder, fmt.Sprintf("%s:%d", filepath[idx:], line))
		}

		if i >= 20 {
			break
		}
	}

	return recorder
}
