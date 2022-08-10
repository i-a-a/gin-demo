package constant

// 用户账号状态
const (
	UserStateOk     uint8 = iota + 1 // 正常
	UserStateFrozen                  // 冻结
)
