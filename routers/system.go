package routers

import (
	"github.com/limeschool/gin"
	"operation/handlers/system"
)

func UseSystemRouter(root *gin.RouterGroup) {
	api := root.Group("/system")
	{
		// 验证码
		api.POST("/captcha", system.Captcha)

		//文件上传相关
		api.POST("/upload", system.UploadFile)

		// 菜单相关
		api.GET("/menus", system.AllMenu)
		api.POST("/menu", system.AddMenu)
		api.PUT("/menu", system.UpdateMenu)
		api.DELETE("/menu", system.DeleteMenu)

		// 角色相关
		api.GET("/roles", system.AllRole)
		api.POST("/role", system.AddRole)
		api.PUT("/role", system.UpdateRole)
		api.DELETE("/role", system.DeleteRole)

		// 角色菜单相关
		api.GET("/role/menu/ids", system.RoleMenuIds)
		api.PUT("/role/menu", system.UpdateRoleMenu)

		// 部门相关
		api.GET("/teams", system.AllTeam)
		api.POST("/team", system.AddTeam)
		api.PUT("/team", system.UpdateTeam)
		api.DELETE("/team", system.DeleteTeam)

		// 用户管理相关
		api.GET("/users", system.PageUser)
		api.GET("/user", system.CurUser)
		api.POST("/user", system.AddUser)
		api.PUT("/user", system.UpdateUser)
		api.DELETE("/user", system.DeleteUser)
		api.GET("/user/menus", system.UserMenus)

		// 用户其他操作
		api.POST("/user/login", system.UserLogin)
		api.POST("/token/refresh", system.RefreshToken)
		api.GET("/login/log", system.LoginLog)
	}
}
