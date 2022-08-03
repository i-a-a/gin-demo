package common

import (
	"app/config"
	"app/pkg/db"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var (
	IsLocal bool
	IsDev   bool
	IsProd  bool
)

var (
	DB    *gorm.DB
	Redis *db.Redis
)

type LogrusHook struct{}

func (l LogrusHook) Fire(entry *logrus.Entry) error {
	switch entry.Level {
	case logrus.ErrorLevel:
		stack, ok := entry.Data["stack"].([]string)
		if ok && !IsLocal && config.Alarm.RobotAPI != "" {
			msg := fmt.Sprintf(`环境: %s
系统: %s
错误: %s
堆栈: %v`, config.App.Env, config.App.System, entry.Message, strings.Join(stack, "\n\t\t\t\t\t"))
			Robot(msg)
		}
	}
	return nil
}

func (l LogrusHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

type send struct {
	MsgType string `json:"msgtype"`
	Text    struct {
		Content string `json:"content"`
	} `json:"text"`
}

func Robot(content string) {
	body := send{MsgType: "text"}
	body.Text.Content = content
	data, _ := json.Marshal(body)
	resp, err := http.Post(config.Alarm.RobotAPI, "application/json", bytes.NewReader(data))
	if err != nil {
		logrus.Warn("机器人挂机", err)
		return
	}
	resp.Body.Close()
}
