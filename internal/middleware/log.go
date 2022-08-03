package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// #### 打印请求和返回信息 ####

var (
	ResponseLog io.Writer = os.Stdout
)

func ApiLog() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer printLog(ctx, copyRequestBody(ctx), time.Now())

		ctx.Next()
	}
}

type apiLogFormat struct {
	Time       string          `json:"time"`
	Path       string          `json:"path"`
	Method     string          `json:"method"`
	Cost       int64           `json:"cost" `          // 毫秒
	Slow       bool            `json:"slow,omitempty"` // 慢请求
	Uid        int             `json:"uid,omitempty"`
	Query      string          `json:"query,omitempty"`
	BodyString string          `json:"body_string,omitempty"`
	BodyJson   json.RawMessage `json:"body_json,omitempty"`
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
		Path:   ctx.Request.URL.Path,
		Method: ctx.Request.Method,
		Uid:    ctx.GetInt("user_id"),
		Query:  ctx.Request.URL.RawQuery,
	}

	if len(bodyData) > 0 && bodyData[0] == '{' {
		apiLoger.BodyJson = bodyData
	} else {
		apiLoger.BodyString, _ = url.QueryUnescape(string(bodyData))
	}

	if apiLoger.Cost > 300 {
		apiLoger.Slow = true
	}

	resp, has := ctx.Get("response")
	if has && resp != nil {
		apiLoger.Response = resp.([]byte)
	}

	b, err := json.Marshal(apiLoger)
	if err != nil {
		logrus.Error(err)
		return
	}
	ResponseLog.Write(b)
	ResponseLog.Write([]byte("\n"))
}
