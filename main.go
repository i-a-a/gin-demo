package main

import (
	"app/config"
	"app/internal"
	"app/internal/controller"
	"fmt"
	"strconv"
)

func main() {

	internal.InitService()

	fmt.Println("启动HTTP服务: " + config.App.Host + ":" + strconv.Itoa(config.App.Port))
	if err := controller.Engine.Run(":" + strconv.Itoa(config.App.Port)); err != nil {
		panic(err)
	}
}
