package dto

import (
	"app/internal/common/helper"
	"app/internal/model"
	"app/pkg/carrot"
	"app/pkg/token"

	"github.com/go-playground/validator/v10"
)

func init() {
	carrot.AddValidate("account", func(fl validator.FieldLevel) bool {
		return helper.Regexp.IsEmail(fl.Field().String())
	})
}

// 尽量最少的字段返回
type UserInfo struct {
	model.BaseID
	Account   string `json:"account"`
	Nickname  string `json:"nickname"`
	Gender    uint8  `json:"gender"`
	Age       uint8  `json:"age"`
	Avatar    string `json:"avatar"`
	CreatedAt string `json:"created_at"`
}

type UserSendCodeReq struct {
	Type    uint8  `json:"type" form:"type" validate:"required"` // 1注册 2登录 3忘记密码 ...
	Account string `json:"account" form:"account" validate:"required"`
	IP      int
}

type UserSendCodeResp = Null

type UserSignUpReq struct {
	Account  string `json:"account" form:"account" validate:"required,account,gte=5,lte=32"`
	Password string `json:"password" form:"password" validate:"omitempty"`
	Captcha  string `json:"captcha" form:"captcha" validate:"omitempty,number,eq=6" label:"验证码"`
	IP       string `json:"-" form:"-"`
}

type UserSignUpResp = Null

type UserSignInReq struct {
	Account  string `json:"account" form:"account" validate:"required,gte=5,lte=32"`
	Password string `json:"password" form:"password" validate:"required"`
	IP       string `json:"-" form:"-"`
}

type UserSignInResp = token.Token

// 更新任意属性
type UserPostInfoReq struct {
	Nickname *string `json:"nickname" form:"nickname" validate:"omitempty,gte=2,lte=20"`
	Avatar   *string `json:"avatar" form:"avatar" validate:"omitempty,gte=10,lte=100"`
	Gender   *uint8  `json:"gender" form:"gender" validate:"omitempty,max=2"`
	Age      *uint8  `json:"age" form:"age" validate:"omitempty,max=100"`
}

type UserPostInfoResp = UserInfo

type UserGetInfoReq = Null

type UserGetInfoResp = UserInfo

type UserListReq struct {
	Pagination
	Nickname string `form:"nickname"`
}

type UserListResp struct {
	Count int64 `json:"count"`
	List  []struct {
		Id       uint32 `json:"id"`
		Nickname string `json:"nickname"`
		Avatar   string `json:"avatar"`
	} `json:"list"`
}
