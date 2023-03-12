package routers

import (
	"github.com/limeschool/gin"
	"operation/middlewares"
)

func InitRouter() *gin.Engine {
	engine := gin.Default()

	// 全局404处理
	engine.NoRoute(gin.Resp404())

	// 健康检查暴露接口
	engine.GET("/healthy", gin.Success())

	// 静态资源文件
	engine.Static("/static", "./static")

	// 注册鉴权中间件
	root := engine.Group("/api", middlewares.JwtAuth.Auth())

	UseSystemRouter(root)

	return engine
}
