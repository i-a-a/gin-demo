package controller

import (
	"app/internal/common/dto"
	"app/internal/common/util"
	"app/internal/middleware"
	"app/internal/service"
	"app/pkg/response"

	"github.com/gin-gonic/gin"
)

// v1
func init() {
	var v1 = Engine.Group("/user/v1")
	var u User
	// 无需登录
	v1.POST("send-code", u.SendCode) // 发送验证码
	v1.POST("sign-up", u.SignUp)     // 注册
	v1.POST("sign-in", u.SignIn)     // 登录

	// 需要登录
	v1.Use(middleware.Auth())
	v1.POST("info", u.PostInfo) // 修改个人信息
	v1.GET("info", u.GetInfo)   // 获取个人信息
}

// v2
func init() {
	// var v2 = Engine.Group("user/v2")
	// var u User
}

// 这个结构体0字节，查一下空struct
type User struct {
	Service service.User
}

func (u User) SendCode(c *gin.Context) {
	var req dto.UserSendCodeReq
	var resp dto.UserSendCodeResp
	if !util.Bind(c, &req) {
		return
	}
	err := u.Service.SendCode(req, &resp)
	response.Echo(c, &resp, err)
}

func (u User) SignUp(c *gin.Context) {
	var req dto.UserSignUpReq
	var resp dto.UserSignUpResp
	if !util.BindAndTrim(c, &req) {
		return
	}
	err := u.Service.SignUp(req, &resp)
	response.Echo(c, &resp, err)
}

func (u User) SignIn(c *gin.Context) {
	var req dto.UserSignInReq
	var resp dto.UserSignInResp
	if !util.BindAndTrim(c, &req) {
		return
	}
	// 这个获取IP函数不准
	req.IP = c.ClientIP()
	err := u.Service.SignIn(req, &resp)
	response.Echo(c, &resp, err)
}

func (u User) PostInfo(c *gin.Context) {
	var req dto.UserPostInfoReq
	var resp dto.UserPostInfoResp
	if !util.BindAndTrim(c, &req) {
		return
	}
	err := u.Service.PostInfo(util.GetUid(c), req, &resp)
	response.Echo(c, &resp, err)
}

func (u User) GetInfo(c *gin.Context) {
	var req dto.UserGetInfoReq
	var resp dto.UserGetInfoResp
	if !util.Bind(c, &req) {
		return
	}
	err := u.Service.GetInfo(util.GetUid(c), req, &resp)
	response.Echo(c, &resp, err)
}
