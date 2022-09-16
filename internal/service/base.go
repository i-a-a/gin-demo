package service

import "sync"

var (
	// 注册互斥锁
	SignUpMu sync.Mutex
)
