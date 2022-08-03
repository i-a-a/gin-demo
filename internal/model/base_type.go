package model

import "time"

// 统一主键！ 判断数据是否有效，使用IsValid()
// mysql的int unsigned是4个字节，对应go的uint32。  big int unsigned 对应uint64。
type ID struct {
	Id uint32 `json:"id" gorm:"primaryKey; <-:false; type:int unsigned auto_increment; comment:主键;"`
}

func (b ID) IsValid() bool {
	return b.Id > 0
}

// 我认为时间类型大多时候只是简单的存取，不需要进行判断，所以用string取代time.Time。
// 使用Datetime，关闭gorm连接时的Parsetime。
// 如果你不认同，自己去封装一个 time.Time。 参考https://gorm.io/zh_CN/docs/data_types.html
type Datetime string

func (d Datetime) ToTime() time.Time {
	t, _ := time.Parse("2006-01-02 15:04:05", string(d))
	return t
}
