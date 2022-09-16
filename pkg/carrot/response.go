package carrot

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type Response struct {
	Code int         `json:"code"` // 0:成功，忽略msg，使用data。 1:失败，弹出msg。  -1:未知错误，需要排查。
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
	Uuid string      `json:"uuid,omitempty"` // 用来定位日志
}

func (c *Carrot) Echo(data interface{}, err error) {
	if c.ctx.IsAborted() {
		logrus.Error("api is aborted")
		return
	}

	var r Response
	if err == nil {
		r.Msg = "ok"
		r.Data = data
	} else {
		r.Msg = err.Error()
		r.Data = struct{}{}
		r.Uuid = uuid.NewString()
		l := logrus.WithField("uuid", r.Uuid)
		switch err.(type) {
		case Msg:
			r.Code = 1
			l.Info(r.Msg)
		default:
			r.Code = -1
			l.Error(r.Msg)
		}
	}

	byteData, err := json.Marshal(&r)
	if err != nil {
		panic(err)
	}

	c.ctx.Abort()
	c.ctx.Data(http.StatusOK, "application/json", byteData)

	printLog(c, byteData)
}
