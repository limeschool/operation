package system

import (
	"github.com/limeschool/gin"
	"operation/errors"
	service "operation/services/system"
	types "operation/types/system"
)

func PageUser(ctx *gin.Context) {
	// 检验参数
	in := types.PageUserRequest{}
	if ctx.ShouldBind(&in) != nil {
		ctx.RespError(errors.ParamsError)
		return
	}

	if list, total, err := service.PageUser(ctx, &in); err != nil {
		ctx.RespError(err)
	} else {
		ctx.RespList(in.Page, len(list), int(total), list)
	}
}

func UserTeamIds(ctx *gin.Context) {
	if data, err := service.CurTeamIds(ctx); err != nil {
		ctx.RespError(err)
	} else {
		ctx.RespData(data)
	}
}

func GetUser(ctx *gin.Context) {
	in := types.GetUserRequest{}
	if ctx.ShouldBind(&in) != nil {
		ctx.RespError(errors.ParamsError)
		return
	}
	if data, err := service.GetUser(ctx, &in); err != nil {
		ctx.RespError(err)
	} else {
		ctx.RespData(data)
	}
}

func CurUser(ctx *gin.Context) {
	if user, err := service.CurUser(ctx); err != nil {
		ctx.RespError(err)
	} else {
		ctx.RespData(user)
	}
}

func AddUser(ctx *gin.Context) {
	// 检验参数
	in := types.AddUserRequest{}
	if err := ctx.ShouldBindJSON(&in); err != nil {
		ctx.RespError(errors.ParamsError)
		return
	}
	// 调用实现
	if err := service.AddUser(ctx, &in); err != nil {
		ctx.RespError(err)
	} else {
		ctx.RespSuccess()
	}
}

func UpdateUser(ctx *gin.Context) {
	// 检验参数
	in := types.UpdateUserRequest{}
	if ctx.ShouldBindJSON(&in) != nil {
		ctx.RespError(errors.ParamsError)
		return
	}
	// 调用实现
	if err := service.UpdateUser(ctx, &in); err != nil {
		ctx.RespError(err)
	} else {
		ctx.RespSuccess()
	}
}

func UpdateUserinfo(ctx *gin.Context) {
	// 检验参数
	in := types.UpdateUserinfoRequest{}
	if ctx.ShouldBindJSON(&in) != nil {
		ctx.RespError(errors.ParamsError)
		return
	}
	// 调用实现
	if err := service.UpdateUserinfo(ctx, &in); err != nil {
		ctx.RespError(err)
	} else {
		ctx.RespSuccess()
	}
}

func UpdatePassword(ctx *gin.Context) {
	// 检验参数
	in := types.UpdatePasswordRequest{}
	if ctx.ShouldBindJSON(&in) != nil {
		ctx.RespError(errors.ParamsError)
		return
	}
	// 调用实现
	if err := service.UpdatePassword(ctx, &in); err != nil {
		ctx.RespError(err)
	} else {
		ctx.RespSuccess()
	}
}

func DeleteUser(ctx *gin.Context) {
	// 检验参数
	in := types.DeleteUserRequest{}
	if ctx.ShouldBindJSON(&in) != nil {
		ctx.RespError(errors.ParamsError)
		return
	}
	// 调用实现
	if err := service.DeleteUser(ctx, &in); err != nil {
		ctx.RespError(err)
	} else {
		ctx.RespSuccess()
	}
}

func UserLogin(ctx *gin.Context) {
	// 检验参数
	in := types.UserLoginRequest{}
	if ctx.ShouldBindJSON(&in) != nil {
		ctx.RespError(errors.ParamsError)
		return
	}
	// 调用实现
	if resp, err := service.UserLogin(ctx, &in); err != nil {
		ctx.RespError(err)
	} else {
		ctx.RespData(resp)
	}
}

func RefreshToken(ctx *gin.Context) {
	if resp, err := service.RefreshToken(ctx); err != nil {
		ctx.RespError(err)
	} else {
		ctx.RespData(resp)
	}
}

func UserLogout(ctx *gin.Context) {
	if err := service.UserLogout(ctx); err != nil {
		ctx.RespError(err)
	} else {
		ctx.RespSuccess()
	}
}

// UserMenus 获取用户的菜单列表
func UserMenus(ctx *gin.Context) {
	if tree, err := service.CurUserMenuTree(ctx); err != nil {
		ctx.RespError(err)
	} else {
		ctx.RespData(tree)
	}
}
