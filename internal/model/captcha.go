package model

import (
	"app/pkg/carrot"
	"errors"
	"strings"
	"time"
)

var (
	CaptchaPtr *Captcha
)

// 验证码
type Captcha struct {
	BaseID
	Account string `json:"account" gorm:"index:idx_account_createtime; size:32; not null; default:''; comment:账号; "`
	Captcha string `json:"captcha" gorm:"size:6; not null; default:''; comment:验证码;"`
	// mysql没有布尔类型，但是true和false会被转换成1和0
	IsUsed   bool  `json:"is_used" gorm:"type:tinyint unsigned; not null; default:0; comment:是否已使用;"`
	Category uint8 `json:"category" gorm:"type:tinyint unsigned; not null; default:0; comment:类型 1注册2更改密码;"`
	IP       int   `json:"ip" gorm:"index:idx_ip_createtime; type:int unsigned; not null; default:0; comment:IP地址;"`
	// 创建时间需要索引
	CreateTime Datetime `json:"create_time" gorm:"index:idx_ip_createtime; index:idx_account_createtime; <-:false; type:datetime; not null; default:now(); comment:创建时间;"`
	ChangeTime Datetime `json:"change_time" gorm:"<-:false; type:datetime; not null; default:now() ON UPDATE now(); comment:更新时间;"`
}

// 校验是否可以发送
func (*Captcha) CanSend(account string, ip int) error {
	var count int64

	// 每个账号每分钟只能发一次
	DB().Model(CaptchaPtr).Where("account = ? AND create_time > ?", account, time.Now().Add(-time.Minute)).Count(&count)
	if count > 0 {
		return carrot.Msg("请稍后再尝试发送验证码")
	}

	// TODO每天每个账号只能发5次

	// IP防刷，返回error触发预警
	DB().Model(CaptchaPtr).Where("ip = ? AND create_time > ?", ip, time.Now().Add(-time.Minute)).Count(&count)
	if count > 10 {
		return errors.New("验证码发送过多")
	}

	return nil
}

// 查询
func (*Captcha) Check(account, captcha string) error {
	var uc Captcha
	captcha = strings.ToUpper(captcha)

	// 事实上，建立一个 account,captcha的索引直接查询，会更快，也更好区分过期和错误，但是我觉得这非必要的。这里利用 account,create_time索引
	DB().Where("account = ? AND create_time > ? AND captcha = ?", account, time.Now().Add(-time.Hour*2), captcha).Take(&uc)
	// 错误或失效
	if !uc.IsValid() {
		return carrot.Msg("验证码错误或已失效")
	}

	// 被使用
	if uc.IsUsed {
		return carrot.Msg("验证码已被使用过，请重新获取")
	}

	// 标记为已使用。（这个更新不太影响业务不放在事务里）
	DebugDB().Model(&uc).Update("is_used", 1)

	return nil
}
