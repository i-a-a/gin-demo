package dto

type Null struct{}

type Pagination struct {
	Page int `form:"page" binding:"required,gte=1" label:"页码"`
	Size int `form:"size" binding:"required,gte=1,lte=100" label:"页大小"`
}
