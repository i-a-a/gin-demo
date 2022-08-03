package dto

import (
	"app/pkg/token"
)

type UserSimple struct {
	Id       uint64 `json:"id"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
}

type UserInfo struct {
	Id        uint32 `json:"id"`
	Gender    uint8  `json:"gender"`
	Age       uint8  `json:"age"`
	Account   string `json:"account"`
	Nickname  string `json:"nickname"`
	Avatar    string `json:"avatar"`
	CreatedAt string `json:"created_at"`
}

type UserRegisterReq struct {
	Nickname string `json:"nickname" form:"nickname"`
	UserLoginReq
}

type UserRegisterResp = Null

type UserVerifyReq struct {
	Account    string `json:"account" form:"account" binding:"required"`
	VerifyCode string `json:"verify_code" form:"verify_code" binding:"required"`
}

type UserVerifyResp = token.Token

type UserLoginReq struct {
	Account  string `json:"account" form:"account" binding:"required,lte=64"`
	Password string `json:"password" form:"password" binding:"required"`
	IP       string `json:"-" form:"-"`
}

type UserLoginResp struct {
	Token token.Token `json:"token"`
	User  UserInfo    `json:"user"`
}

type UserRefresgTokenReq struct {
	AccessToken  string `json:"access_token" form:"access_token" binding:"required"`
	RefreshToken string `json:"refresh_token" form:"refresh_token" binding:"required"`
}

type UserSendCodeReq struct {
	Type    uint8  `json:"type" form:"type" binding:"required"`
	Account string `json:"account" form:"account" binding:"required"`
}

type UserSendCodeResp = Null

type UserPostInfoReq struct {
	Nickname string `json:"nickname" form:"nickname" binding:"omitempty,lte=32"`
	Avatar   string `json:"avatar" form:"avatar" binding:"omitempty,lte=255"`
	Gender   uint8  `json:"gender" form:"gender" binding:"omitempty,min=1,max=2"`
	Age      uint8  `json:"age" form:"age" binding:"omitempty,min=1,max=130"`
}

type UserPostInfoResp = UserInfo

type UserGetInfoReq = Null

type UserGetInfoResp = UserInfo
