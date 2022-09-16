package carrot

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"strconv"
	"time"
)

// ====== 记录Api请求 ======

var (
	logWriter io.Writer = os.Stdout
)

func SetLogger(w io.Writer) {
	logWriter = w
}

type logFormat struct {
	Time          string          `json:"time"`
	Cost          int64           `json:"cost,omitempty"` // 毫秒
	Method        string          `json:"method"`
	Match         string          `json:"match"` // ${用户uid}-${路由} 例：1001-/user/info
	RequestString string          `json:"request_string,omitempty"`
	RequestJSON   json.RawMessage `json:"request_json,omitempty"`
	Response      json.RawMessage `json:"response"`
}

func printLog(c *Carrot, response []byte) {
	f := logFormat{
		Time:   c.timer.Format("2006-01-02 15:04:05"),
		Cost:   time.Since(c.timer).Milliseconds(),
		Method: c.ctx.Request.Method,
	}
	// match 在搜索时，使用 ${用户uid}-${路由} 格式
	if c.Uid > 0 {
		f.Match += strconv.Itoa(c.Uid) + "-"
	}
	f.Match += c.ctx.Request.URL.Path
	if c.ctx.Request.URL.RawQuery != "" {
		f.Match += "?" + c.ctx.Request.URL.RawQuery
	}
	// 请求body，json时使用json.RawMessage便于观看
	if len(c.requestData) > 0 {
		if c.requestData[0] == '{' {
			f.RequestJSON = c.requestData
		} else {
			f.RequestString = string(c.requestData)
		}
	}
	// 返回数据也是json.RawMessage
	if c.Switch.ResponseLog {
		f.Response = json.RawMessage(response)
	}

	bs, _ := jsonEncode(&f)
	logWriter.Write(bs)
	logWriter.Write([]byte("\n"))
}

// 解决json.Marshal时转义问题
func jsonEncode(v interface{}) ([]byte, error) {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(v)
	return buffer.Bytes(), err
}
