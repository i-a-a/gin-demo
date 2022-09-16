package service

import (
	"app/internal/common/dto"
	"app/internal/common/enum"
	"app/internal/model"
	"app/pkg/carrot"
	"app/pkg/token"
	"app/pkg/util"

	"github.com/jinzhu/copier"

	log "github.com/sirupsen/logrus"

	"gorm.io/gorm"
)

type User struct{}

// 发送验证码
func (User) SendCode(req dto.UserSendCodeReq, resp *dto.UserSendCodeResp) error {
	// 验证账号状态
	var count int64
	model.DB().Model(model.UserPtr).Where("account = ?", req.Account).Count(&count)
	switch req.Type {
	case enum.CaptchaTypeSignUp:
		if count > 0 {
			return carrot.Msg("账号已被注册")
		}
	case enum.CaptchaTypeSignIn:
		fallthrough // go没有break，自动break，需要主动fallthrough穿透
	case enum.CaptchaTypeChangePwd:
		if count == 0 {
			return carrot.Msg("账号不存在")
		}
	default:
		return carrot.Msg("非法的验证码类型")
	}

	// 检查发送频率
	if err := model.CaptchaPtr.CanSend(req.Account, req.IP); err != nil {
		return err
	}

	captcha := util.RandString(6)
	data := model.Captcha{
		Account:  req.Account,
		Captcha:  captcha,
		Category: req.Type,
		IP:       req.IP,
	}
	err := model.DebugDB().Create(&data).Error
	if err != nil {
		return err
	}

	// TODO 发送验证码，这里不需要等待
	// 非引用的数据类型，在使用goroutine时，最好作为参数传进去。 （内存逃逸考虑、循环时的变量共享问题）
	go func(account, captcha string) {

	}(req.Account, captcha)

	return nil
}

// 注册
func (User) SignUp(req dto.UserSignUpReq, resp *dto.UserSignUpResp) error {
	SignUpMu.Lock()
	defer SignUpMu.Unlock()

	// // TO
	// if !helper.Regexp.IsEmail(req.Account) {

	// }

	// 检测账号重复
	var count int64
	model.DebugDB().Model(model.UserPtr).Where("account = ?", req.Account).Count(&count)
	if count > 0 {
		return carrot.Msg("账号已被注册")
	}

	// 检查验证码
	err := model.CaptchaPtr.Check(req.Account, req.Captcha)
	if err != nil {
		return err
	}

	salt := util.RandString(10)
	password := util.HmacSha256(req.Password, salt)
	err = model.DebugDB().Transaction(func(tx *gorm.DB) error {
		user := &model.User{
			Account:  req.Account,
			Password: password,
			Salt:     salt,
			State:    enum.UserStateOk,
		}
		err := tx.Create(user).Error
		if err != nil {
			return err
		}

		log.Info("新注册账号: " + req.Account)

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

// 登录
func (User) SignIn(req dto.UserSignInReq, resp *dto.UserSignInResp) error {
	user := model.UserPtr.GetByAccount(req.Account)
	if !user.IsValid() {
		return enum.MsgUserNotFound
	}
	if user.Password != util.HmacSha256(req.Password, user.Salt) {
		return carrot.Msg("密码错误")
	}
	if user.State == enum.UserStateFrozen {
		return carrot.Msg("账号已被禁用")
	}

	*resp = token.GenerateTokens(user.Id)

	log.Info("用户登录：", user.Id)

	return nil
}

func (User) GetInfo(uid int, req dto.UserGetInfoReq, resp *dto.UserGetInfoResp) error {
	model.DB().Model(model.UserPtr).Take(&resp, uid)
	if !resp.IsValid() {
		return enum.ErrUserNotFound
	}

	return nil
}

func (User) PostInfo(uid int, req dto.UserPostInfoReq, resp *dto.UserGetInfoResp) error {
	model.DB().Model(model.UserPtr).Take(&resp, uid)
	if !resp.IsValid() {
		return enum.ErrUserNotFound
	}

	// copy ，两个参数都是指针，相同字段后者覆盖前者
	copier.Copy(resp, &req)

	var update model.User
	copier.Copy(&update, &req)
	update.Id = uid
	err := model.DebugDB().Model(&update).Updates(&update).Error

	return err
}

func (User) List(req dto.UserListReq, resp *dto.UserListResp) error {
	q := model.DB().Model(model.UserPtr)
	if req.Nickname != "" {
		q = q.Where("nickname like ?", "%"+req.Nickname+"%")
	}
	err := req.AutoFind(q, &resp.Count, &resp.List)
	return err
}
