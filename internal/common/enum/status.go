package enum

// 用户账号状态
const (
	UserStateOk     = iota + 1 // 正常
	UserStateFrozen            // 冻结
)

// 验证码类型
const (
	CaptchaTypeSignUp    = iota + 1 // 注册
	CaptchaTypeSignIn               // 登录
	CaptchaTypeChangePwd            // 修改密码
)
