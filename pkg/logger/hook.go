package logger

import (
	"github.com/sirupsen/logrus"
)

var (
	hookMap = make(map[logrus.Level][]hookFunc, 2)

	_ logrus.Hook = hook{}
)

type hookFunc func(string)

// Error触发钩子
func AddErrorHooks(fns ...hookFunc) {
	hookMap[logrus.ErrorLevel] = append(hookMap[logrus.ErrorLevel], fns...)
}

// Warn触发钩子
func AddWarnHooks(fns ...hookFunc) {
	hookMap[logrus.WarnLevel] = append(hookMap[logrus.WarnLevel], fns...)
}

type hook struct{}

func (h hook) Fire(entry *logrus.Entry) error {
	for _, f := range hookMap[entry.Level] {
		f(entry.Message)
	}
	return nil
}

func (h hook) Levels() []logrus.Level {
	levels := make([]logrus.Level, 0, len(hookMap))
	for k := range hookMap {
		levels = append(levels, k)
	}
	return levels
}
