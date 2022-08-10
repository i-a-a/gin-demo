package dto

import (
	"app/pkg/token"
)

type UserInfo struct {
	Id        uint32 `json:"id"`
	Gender    uint8  `json:"gender"`
	Age       uint8  `json:"age"`
	Account   string `json:"account"`
	Nickname  string `json:"nickname"`
	Avatar    string `json:"avatar"`
	CreatedAt string `json:"created_at"`
}

type UserSendCodeReq struct {
	Type    uint8  `json:"type" form:"type" validate:"required"` // 1注册 2忘记密码 ...
	Account string `json:"account" form:"account" validate:"required"`
}

type UserSendCodeResp = Null

type UserSignUpReq struct {
	UserSignInReq
	VerifyCode string `json:"verify_code" validate:"required"`
	Nickname   string `json:"nickname" form:"nickname"`
}

type UserSignUpResp = Null

type UserSignInReq struct {
	Account  string `json:"account" form:"account" validate:"required,gte=5,lte=64"`
	Password string `json:"password" form:"password" validate:"required"`
	IP       string `json:"-" form:"-"`
}

type UserSignInResp = token.Token

type UserPostInfoReq struct {
	Nickname string `json:"nickname" form:"nickname" validate:"omitempty,lte=32"`
	Avatar   string `json:"avatar" form:"avatar" validate:"omitempty,lte=255"`
	Gender   uint8  `json:"gender" form:"gender" validate:"omitempty,min=1,max=2"`
	Age      uint8  `json:"age" form:"age" validate:"omitempty,min=1,max=130"`
}

type UserPostInfoResp = UserInfo

type UserGetInfoReq = Null

type UserGetInfoResp = UserInfo
