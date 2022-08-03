package service

import (
	"app/internal/common"
	"app/internal/common/constant"
	"app/internal/common/dto"
	"app/internal/model"
	"app/pkg/response"
	"app/pkg/token"
	"app/pkg/util"
	"sync"

	"github.com/jinzhu/copier"

	log "github.com/sirupsen/logrus"

	"gorm.io/gorm"
)

var (
	mu sync.Mutex
)

type User struct{}

func (u User) Register(req dto.UserRegisterReq, resp *dto.UserRegisterResp) error {
	mu.Lock()
	defer mu.Unlock()

	// 检测账号重复
	var count int64
	common.DB.Model(model.UserPtr).Where("account = ?", req.Account).Count(&count)
	if count > 0 {
		return response.String("账号已被注册")
	}

	err := common.DB.Transaction(func(tx *gorm.DB) error {
		salt := util.RandomString(10)
		password := util.HmacSha256(req.Password, salt)
		user := model.User{
			Account:  req.Account,
			Password: password,
			Salt:     salt,
			State:    constant.UserStateOk,
		}
		if err := tx.Create(&user).Error; err != nil {
			return err
		}

		// TODO 其它逻辑 。

		log.Info("新注册账号: " + req.Account)

		return nil
	})

	return err
}

func (u User) Verify(req dto.UserVerifyReq, resp *dto.UserVerifyResp) error {
	var user model.User

	common.DB.Select("id", "state").Where("account = ?", req.Account).Take(&user)
	if user.Id == 0 {
		return response.String("账号不存在")
	}

	if user.State != constant.UserStateUncheck {
		return response.String("账号已验证")
	}

	// ------TODO 验证-------

	// ---------------------

	common.DB.Model(&user).UpdateColumn("state", constant.UserStateOk)

	*resp = token.GenerateTokens(int(user.Id))

	return nil
}

func (User) Login(req dto.UserLoginReq, resp *dto.UserLoginResp) error {
	user := model.UserPtr.GetByAccount(req.Account)
	if !user.IsValid() {
		return response.String("用户不存在")
	}

	// 检查密码
	if user.Password != util.HmacSha256(req.Password, user.Salt) {
		return response.String("密码错误")
	}

	switch user.State {
	case constant.UserStateUncheck:
		return response.String("请先验证账号")
	case constant.UserStateFrozen:
		return response.String("账号已被冻结")
	}

	copier.Copy(&resp.User, &user)
	resp.Token = token.GenerateTokens(int(resp.User.Id))

	log.Info("用户登录：", resp.User.Id)
	// === TODO 登录记录 ===

	return nil
}

func (User) SendCode(req dto.UserSendCodeReq, resp *dto.UserSendCodeResp) error {

	return nil
}

func (u User) PostInfo(uid int, req dto.UserPostInfoReq, resp *dto.UserGetInfoResp) error {
	err := common.DB.Model(model.UserPtr).Take(&resp, uid).Error
	if err != nil {
		return err
	}

	// 合并信息
	copier.Copy(&resp, &req)

	return common.DB.Model(model.UserPtr).Where("id = ?", resp.Id).Updates(&req).Error
}

func (u User) GetInfo(uid int, req dto.UserGetInfoReq, resp *dto.UserGetInfoResp) error {
	err := model.UserPtr.DB().Take(&resp, uid).Error
	return err
}
