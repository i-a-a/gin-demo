package pkg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
)

var (
	Robot    robot
	RobotApi string
)

// 飞书 https://open.feishu.cn/document/ukTMukTMukTM/ucTM5YjL3ETO24yNxkjN
// 钉钉 https://open.dingtalk.com/document/group/custom-robot-access

type robot struct{}

// 发送到机器人
func (robot) SendText(env string, system string, msg string) {
	if RobotApi == "" {
		logrus.Warn("未配置机器人api")
		return
	}
	content := fmt.Sprintf("环境: %s\n系统: %s\n错误: %v", env, system, msg)
	// 这个是钉钉格式，别的自己去找一下
	body := map[string]interface{}{
		"msgtype": "text",
		"text": map[string]interface{}{
			"content": content,
		},
	}
	data, _ := json.Marshal(body)
	resp, err := http.Post(RobotApi, "application/json", bytes.NewReader(data))
	if err != nil {
		fmt.Println("机器人挂机", err)
		return
	}
	// b, _ := ioutil.ReadAll(resp.Body)
	// fmt.Println(string(b))
	resp.Body.Close()
}
