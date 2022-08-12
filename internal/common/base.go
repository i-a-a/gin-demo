package common

import "fmt"

var (
	IsLocal bool
	IsDev   bool
	IsProd  bool
)

func SetEnv(env string) {
	switch env {
	case "local":
		IsLocal = true
	case "prod":
		IsProd = true
	default:
		IsDev = true
	}
	fmt.Println("当前环境: " + env)
}
