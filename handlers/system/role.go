package system

import (
	"github.com/limeschool/gin"
	"operation/consts"
	"operation/errors"
	service "operation/services/system"
	"operation/tools"
	types "operation/types/system"
)

func AllRole(ctx *gin.Context) {
	if resp, err := service.AllRole(ctx); err != nil {
		ctx.RespError(err)
	} else {
		ctx.RespData(resp)
	}
}

func AddRole(ctx *gin.Context) {
	in := types.AddRoleRequest{}
	if ctx.ShouldBindJSON(&in) != nil {
		ctx.RespError(errors.ParamsError)
		return
	}

	if !tools.InList([]string{consts.ALLTEAM, consts.DOWNTEAM, consts.CURTEAM, consts.CUSTOM}, in.DataScope) {
		ctx.RespError(errors.ParamsError)
		return
	}

	if err := service.AddRole(ctx, &in); err != nil {
		ctx.RespError(err)
	} else {
		ctx.RespSuccess()
	}
}

func UpdateRole(ctx *gin.Context) {
	in := types.UpdateRoleRequest{}
	if ctx.ShouldBindJSON(&in) != nil {
		ctx.RespError(errors.ParamsError)
		return
	}

	if !tools.InList([]string{consts.ALLTEAM, consts.DOWNTEAM, consts.CURTEAM, consts.CUSTOM}, in.DataScope) {
		ctx.RespError(errors.ParamsError)
		return
	}

	if err := service.UpdateRole(ctx, &in); err != nil {
		ctx.RespError(err)
	} else {
		ctx.RespSuccess()
	}
}

func DeleteRole(ctx *gin.Context) {
	in := types.DeleteRoleRequest{}
	if ctx.ShouldBind(&in) != nil {
		ctx.RespError(errors.ParamsError)
		return
	}
	if err := service.DeleteRole(ctx, &in); err != nil {
		ctx.RespError(err)
	} else {
		ctx.RespSuccess()
	}
}
