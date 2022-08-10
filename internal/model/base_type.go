package model

import (
	"time"

	"gorm.io/gorm"
)

// 统一主键！ 判断数据是否有效，使用IsValid()
// mysql的int unsigned是4个字节，对应go的uint32。  big int unsigned 对应uint64。
type BaseID struct {
	Id uint32 `json:"id" gorm:"primaryKey; <-:false; type:int unsigned auto_increment; comment:主键;"`
}

// 替代 xx.Id > 0
func (b BaseID) IsValid() bool {
	return b.Id > 0
}

// 统一时间字段
// 禁用gorm更改这两个字段，依赖数据库维护
type BaseTimes struct {
	CreatedAt Datetime `json:"created_at" gorm:"<-:false; type:datetime; not null; default:now(); comment:创建时间;"`
	UpdatedAt Datetime `json:"updated_at" gorm:"<-:false; type:datetime; not null; default:now() ON UPDATE now(); comment:更新时间;"`
}

// 带上软删字段，则Delete为标记删除。 查询时默认带上排除删除条件
// 使用软删字段须注意唯一索引问题
type BaseDeleteTime struct {
	DeletedAt gorm.DeletedAt `json:"-" gorm:"type:datetime; comment:删除时间;"`
}

// 我认为时间类型大多时候只是简单的存取，不需要进行判断，所以用string取代time.Time。
// 使用Datetime，须关闭gorm连接时的Parsetime，且数据库该字段自动为当前时间戳。
// 如果你不认同，自己去封装一个 time.Time。 参考https://gorm.io/zh_CN/docs/data_types.html
type Datetime string

func (d Datetime) ToTime() time.Time {
	t, _ := time.ParseInLocation("2006-01-02 15:04:05", string(d), time.Local)
	return t
}
