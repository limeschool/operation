package main

import (
	"operation/config"
	"operation/routers"
)

func main() {
	// 初始化路由
	engine := routers.InitRouter()

	// 初始化配置监听
	config.InitConfig()

	// 启动
	engine.Run(":8081")
}
