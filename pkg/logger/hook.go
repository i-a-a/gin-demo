package logger

import "github.com/sirupsen/logrus"

var (
	hookMap = make(map[logrus.Level][]func(), 2)

	_ logrus.Hook = hook{}
)

// Error触发钩子
func AddErrorHooks(fns ...func()) {
	hookMap[logrus.ErrorLevel] = append(hookMap[logrus.ErrorLevel], fns...)
}

// Warn触发钩子
func AddWarnHooks(fns ...func()) {
	hookMap[logrus.WarnLevel] = append(hookMap[logrus.WarnLevel], fns...)
}

type hook struct{}

func (h hook) Fire(entry *logrus.Entry) error {
	for _, f := range hookMap[entry.Level] {
		f()
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

// msg := fmt.Sprintf("环境: %s\n系统: %s\n错误: %v", config.App.Env, config.App.System, strings.Split(entry.Message, " "))
