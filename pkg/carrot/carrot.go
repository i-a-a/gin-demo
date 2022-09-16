package carrot

import (
	"encoding/json"
	"time"

	"github.com/gin-gonic/gin"
)

type Carrot struct {
	ctx   *gin.Context
	Uid   int
	timer time.Time

	requestData json.RawMessage // 请求数据
	Switch      struct {
		ResponseLog bool // 是否记录响应日志
		RequestLog  bool // 是否记录请求日志
	}
}

func New(ctx *gin.Context) *Carrot {
	carrot := &Carrot{
		ctx:   ctx,
		Uid:   ctx.GetInt("uid"),
		timer: time.Now(),
	}
	carrot.Switch.RequestLog = true
	carrot.Switch.ResponseLog = true

	return carrot
}

// // 关闭日志。有些接口因为敏感性例如登录，有些接口调用太频繁或者返回数据   一个是为了示范选项模式
// func CloseLog() CarrotOption {
// 	return func(c *Carrot) {
// 		c.canLog = false
// 	}
// }
