package model

var (
	UserPtr *User
)

// user虽然是关键字，但是并不影响使用。 （不像order、group那样不能使用）
type User struct {
	BaseID

	// 主要字段放在前边。  账号唯一，昵称是否唯一按需求。
	Nickname string `json:"nickname" gorm:"type:varchar(20); not null; default:''; comment:昵称;"`
	Account  string `json:"account" gorm:"uniqueIndex:uk_account; type:varchar(32); not null; default:''; comment:账号;"`

	// 数字尽量小，尽量注意内存对齐
	State  uint8 `json:"state" gorm:"type:tinyint unsigned; not null; default:3; comment:状态 1正常2冻结3未验证;"`
	Gender uint8 `json:"gender" gorm:"type:tinyint; not null; default:3; comment:性别 1男2女3未知;"`
	Age    uint8 `json:"age" gorm:"type:tinyint; not null; default:0; comment:年龄;"`

	// 很多链接，可以不存储域名，只存储路径。倒不是节省空间，而是迁移方便。  例如头像，这里使用自定义类型，通过gorm钩子自动处理。
	Avatar   `json:"avatar" gorm:"type:varchar(100); not null; default:''; comment:头像;"`
	Password string `json:"-" gorm:"column:passwd; type:varchar(100); not null; default:''; comment:密码;"`
	Salt     string `json:"-" gorm:"type:varchar(10); not null; default:''; comment:密码盐;"`

	BaseTimes
}

// func (u *User) DB() *gorm.DB {
// 	if u != nil && u.IsValid() {
// 		return DB().Debug().Model(u)
// 	}
// 	return DB().Model(u)
// }

// 一些常用的方法，可以直接写在model里
func (*User) GetByAccount(account string) (row User) {
	DB().Where("account = ?", account).Take(&row)
	return
}
