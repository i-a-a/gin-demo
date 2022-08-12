package response

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// code为0，忽略msg； 1表示拒绝，弹出msg； -1表示未知错误，弹出msg；其它表示约定动作。
type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func Echo(ctx *gin.Context, data interface{}, err error) {
	var r Response
	if err == nil {
		r.Msg = "ok"
		r.Data = data
	} else {
		r.Msg = err.Error()
		r.Data = struct{}{}
		// 区分 拒绝、动作、未知错误
		switch x := err.(type) {
		case Msg:
			r.Code = 1
			logrus.Info(r.Msg)
		case Code:
			r.Code = int(x)
			logrus.Warn(r.Msg)
		default:
			r.Code = -1
			s := getErrorStack(err.Error(), "internal")
			logrus.WithField("stack", s).Error(r.Msg)
		}
	}

	byteData, err := json.Marshal(&r)
	if err != nil {
		panic(err)
	}

	ctx.Abort()
	ctx.Data(http.StatusOK, "application/json", byteData)

	// 这里和中间件log配合
	ctx.Set("response", byteData)
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
