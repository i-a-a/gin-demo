package response

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Response struct {
	// 0:成功，-1:出错/拒绝，直接弹出消息。   其他：有含义的动作
	Code int `json:"code"`
	//
	Msg string `json:"msg"`
	//
	Data interface{} `json:"data"`
	// 用来排查问题
	RequestId string `json:"request_id,omitempty"`
}

func Echo(ctx *gin.Context, data interface{}, err error) {
	switch x := err.(type) {
	// 成功
	case nil:
		Success(ctx, data)
	// 失败：可预知，前端执行动作
	case Action:
		Fail(ctx, x.Code, x.Msg)
	// 拒绝：可预知，前端显示信息
	case String: // 这里的code码也可以使用别的。使用-1的话，就和error不区分。
		Fail(ctx, -1, x.Error())
	// 错误：不可预知。
	default:
		Error(ctx, err)
	}
}

func Success(ctx *gin.Context, data interface{}) {
	resp := Response{
		Code: 0,
		Msg:  "ok",
		Data: data,
	}
	done(ctx, resp)
}

func Fail(ctx *gin.Context, code int, msg string) {
	resp := Response{
		Code:      code,
		Msg:       msg,
		Data:      struct{}{},
		RequestId: uuid.NewString(),
	}
	done(ctx, resp)
}

// 增加一个错误钩子，用来排查问题
var (
	errorHook func(errorMsg string)
)

func AddErrorHook(fn func(string)) {
	errorHook = fn
}

func Error(ctx *gin.Context, err error) {
	if errorHook != nil {
		errorHook(err.Error())
	}
	resp := Response{
		Code:      -1,
		Msg:       err.Error(),
		Data:      struct{}{},
		RequestId: uuid.NewString(),
	}
	done(ctx, resp)
}

func done(ctx *gin.Context, resp Response) {
	byteData, err := json.Marshal(resp)
	if err != nil {
		panic(err)
	}

	ctx.Abort()
	ctx.Data(http.StatusOK, "application/json", byteData)

	// 这里和中间件log打配合，感觉分散在两处不太优雅，但是暂时找不到好方法。
	ctx.Set("response", byteData)
}
