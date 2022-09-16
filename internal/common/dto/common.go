package dto

import "gorm.io/gorm"

type Null struct{}

type Pagination struct {
	Page int `json:"page" form:"page" binding:"omitempty" label:"页码"`
	Size int `json:"size" form:"size" binding:"omitempty,lte=100" label:"页大小"`
}

// 深度分页的话，这么写肯定不好。
func (p Pagination) AutoFind(tx *gorm.DB, count *int64, to interface{}) error {
	if p.Page == 0 {
		p.Page = 1
	}
	if p.Size == 0 {
		p.Size = 10
	}
	return tx.Count(count).Offset((p.Page - 1) * p.Size).Limit(p.Size).Find(to).Error
}
