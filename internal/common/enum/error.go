package enum

import (
	"app/pkg/carrot"
	"errors"
)

// 错误
var (
	ErrUserNotFound = errors.New("用户不存在") // 用户数据出错，不合理时返回
)

// 失败
var (
	MsgUserNotFound = carrot.Msg("用户不存在") // 一般情况，正常逻辑时返回
)
