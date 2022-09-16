package controller

import (
	"app/internal/common/dto"
	"app/internal/middleware"
	"app/internal/service"
	"app/pkg/carrot"
	"app/pkg/util"

	"github.com/gin-gonic/gin"
)

func init() {
	var v1 = Engine.Group("/user")
	var u User
	// 无需登录
	v1.POST("send-code", u.SendCode) // 发送验证码
	v1.POST("sign-up", u.SignUp)     // 注册
	v1.POST("sign-in", u.SignIn)     // 登录

	v1.GET("list", u.List) // 用户列表

	// 需要登录
	v1.Use(middleware.Auth())
	v1.POST("info", u.PostInfo) // 修改个人信息
	v1.GET("info", u.GetInfo)   // 获取个人信息

}

// 这个结构体是空结构体，0字节
// 此User的函数都是结构体的函数，都是拷贝一个空结构体
type User struct {
	Service service.User
}

func (u User) SendCode(c *gin.Context) {
	app := carrot.New(c)
	var (
		req  dto.UserSendCodeReq
		resp dto.UserSendCodeResp
	)
	if !app.BindAndTrim(&req) {
		return
	}
	req.IP = util.Ip2Int(c.ClientIP())
	err := u.Service.SendCode(req, &resp)
	app.Echo(&resp, err)
}

func (u User) SignUp(c *gin.Context) {
	app := carrot.New(c)
	app.Switch.RequestLog = false // 关闭请求数据
	var req dto.UserSignUpReq
	if !app.BindAndTrim(&req) {
		return
	}
	var resp dto.UserSignUpResp
	err := u.Service.SignUp(req, &resp)
	app.Echo(&resp, err)
}

func (u User) SignIn(c *gin.Context) {
	app := carrot.New(c)
	var req dto.UserSignInReq
	if !app.BindAndTrim(&req) {
		return
	}
	req.IP = c.ClientIP()
	var resp dto.UserSignInResp
	err := u.Service.SignIn(req, &resp)
	app.Echo(&resp, err)
}

func (u User) PostInfo(c *gin.Context) {
	app := carrot.New(c)
	var req dto.UserPostInfoReq
	if !app.BindAndTrim(&req) {
		return
	}
	var resp dto.UserPostInfoResp
	err := u.Service.PostInfo(app.Uid, req, &resp)
	app.Echo(resp, err)
}

func (u User) GetInfo(c *gin.Context) {
	app := carrot.New(c)
	var req dto.UserGetInfoReq
	if !app.Bind(&req) {
		return
	}
	var resp dto.UserGetInfoResp
	err := u.Service.GetInfo(app.Uid, req, &resp)
	app.Echo(&resp, err)
}

func (u User) List(c *gin.Context) {
	app := carrot.New(c)
	var req dto.UserListReq
	if !app.Bind(&req) {
		return
	}
	var resp dto.UserListResp
	err := u.Service.List(req, &resp)
	app.Echo(resp, err)
}
