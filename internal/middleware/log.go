package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// #### 打印请求和返回信息 ####

var (
	ResponseLogger io.Writer = os.Stdout
)

func ApiLog() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer printLog(ctx, copyRequestBody(ctx), time.Now())

		ctx.Next()
	}
}

type apiLogFormat struct {
	Time       string          `json:"time"`
	Match      string          `json:"match"` // ${用户uid}-${路由}
	Method     string          `json:"method"`
	Cost       int64           `json:"cost,omitempty"` // 毫秒
	BodyString string          `json:"reqString,omitempty"`
	BodyJson   json.RawMessage `json:"reqJson,omitempty"`
	Response   json.RawMessage `json:"response"`
}

func copyRequestBody(ctx *gin.Context) []byte {
	if ctx.Request.Method == http.MethodGet {
		return nil
	}

	bodyBytes, _ := ctx.GetRawData()
	ctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	return bodyBytes
}

func printLog(ctx *gin.Context, bodyData []byte, start time.Time) {
	apiLoger := apiLogFormat{
		Time:   time.Now().Format("2006-01-02 15:04:05"),
		Cost:   time.Since(start).Milliseconds(),
		Method: ctx.Request.Method,
	}
	// match 在搜索时，使用 ${用户uid}-${路由} 格式
	uid := ctx.GetInt("uid")
	if uid > 0 {
		apiLoger.Match += strconv.Itoa(uid) + "-"
	}
	apiLoger.Match += ctx.Request.URL.Path
	if ctx.Request.URL.RawQuery != "" {
		apiLoger.Match += "?" + ctx.Request.URL.RawQuery
	}

	// 请求body，json时使用json.RawMessage便于观看
	if len(bodyData) > 0 && bodyData[0] == '{' {
		apiLoger.BodyJson = bodyData
	} else {
		apiLoger.BodyString = string(bodyData)
	}

	// 返回数据也是json.RawMessage
	resp, has := ctx.Get("response")
	if has && resp != nil {
		apiLoger.Response = resp.([]byte)
	}

	// 自定义的json编码，避免转义
	b, err := JSONMarshal(apiLoger)
	if err != nil {
		logrus.Error(err)
		return
	}
	ResponseLogger.Write(b)
	ResponseLogger.Write([]byte("\n"))
}

func JSONMarshal(t interface{}) ([]byte, error) {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(t)
	return buffer.Bytes(), err
}
