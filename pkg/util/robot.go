package util

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	Robot     robot
	robotConf struct {
		Api string
		Key string
	}
)

func init() {
	viper.UnmarshalKey("robot", &robotConf)
}

// 飞书 https://open.feishu.cn/document/ukTMukTMukTM/ucTM5YjL3ETO24yNxkjN
// 钉钉 https://open.dingtalk.com/document/group/custom-robot-access

type robot struct{}

// 发送到机器人
func (robot) SendText(msg string) {
	if robotConf.Api == "" {
		return
	}

	// 这个是钉钉格式，别的自己去找一下
	body := map[string]interface{}{
		"msgtype": "text",
		"text": map[string]interface{}{
			"content": msg,
		},
	}
	data, _ := json.Marshal(body)

	api := robotConf.Api
	// 签名 https://open.dingtalk.com/document/robots/customize-robot-security-settings
	if robotConf.Key != "" {
		unix := strconv.Itoa(int(time.Now().Unix() * 1000))
		h := hmac.New(sha256.New, []byte(robotConf.Key))
		h.Write([]byte(unix + "\n" + robotConf.Key))
		sign := base64.StdEncoding.EncodeToString(h.Sum(nil))
		api += "&timestamp=" + unix + "&sign=" + sign
	}

	resp, err := http.Post(api, "application/json", bytes.NewReader(data))
	if err != nil {
		logrus.WithField("api", api).Warn("机器人挂机：", err)
		return
	}
	// b, _ := ioutil.ReadAll(resp.Body)
	// fmt.Println(string(b))
	// fmt.Println(api)
	resp.Body.Close()
}
