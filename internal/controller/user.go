package controller

import (
	"app/internal/common/dto"
	"app/internal/middleware"
	"app/internal/service"
	"app/pkg/response"

	"github.com/gin-gonic/gin"
)

func init() {
	var r = Engine.Group("/user")
	var u User

	v1 := r.Group("/v1")

	// 无需登录
	v1.POST("register", u.Register)  // 注册
	v1.POST("verify", u.Verify)      // 验证账号
	v1.POST("login", u.Login)        // 登录
	v1.POST("send_code", u.SendCode) // 发送验证码

	// 需要登录
	v1.Use(middleware.Auth())
	v1.POST("info", u.PostInfo) // 修改个人信息
	v1.GET("info", u.GetInfo)   // 获取个人信息
}

// 这个结构体0字节， 不懂的话去查一下空struct
type User struct {
	Service service.User
}

func (u User) Register(c *gin.Context) {
	var req dto.UserRegisterReq
	if !BindAndTrim(c, &req) {
		return
	}
	var resp dto.UserRegisterResp
	err := u.Service.Register(req, &resp)
	response.Echo(c, &resp, err)
}

func (u User) Verify(c *gin.Context) {
	var req dto.UserVerifyReq
	if !Bind(c, &req) {
		return
	}
	var resp dto.UserVerifyResp
	err := u.Service.Verify(req, &resp)
	response.Echo(c, &resp, err)
}

func (u User) Login(c *gin.Context) {
	var req dto.UserLoginReq
	if !BindAndTrim(c, &req) {
		return
	}
	// 这个获取IP函数不准
	req.IP = c.ClientIP()
	var resp dto.UserLoginResp
	err := u.Service.Login(req, &resp)
	response.Echo(c, &resp, err)
}

func (u User) SendCode(c *gin.Context) {
	var req dto.UserSendCodeReq
	if !Bind(c, &req) {
		return
	}
	var resp dto.UserSendCodeResp
	err := u.Service.SendCode(req, &resp)
	response.Echo(c, &resp, err)
}

func (u User) PostInfo(c *gin.Context) {
	var req dto.UserPostInfoReq
	if !BindAndTrim(c, &req) {
		return
	}
	uid := GetUid(c)
	var resp dto.UserPostInfoResp
	err := u.Service.PostInfo(uid, req, &resp)
	response.Echo(c, &resp, err)
}

func (u User) GetInfo(c *gin.Context) {
	var req dto.UserGetInfoReq
	c.ShouldBind(&req)
	uid := GetUid(c)
	var resp dto.UserGetInfoResp
	err := u.Service.GetInfo(uid, req, &resp)
	response.Echo(c, &resp, err)
}
