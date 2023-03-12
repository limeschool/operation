package system

import (
	"github.com/limeschool/gin"
	"operation/errors"
	service "operation/services/system"
	types "operation/types/system"
)

func UpdateRoleMenu(ctx *gin.Context) {
	in := types.AddRoleMenuRequest{}
	if err := ctx.ShouldBindJSON(&in); err != nil {
		ctx.RespError(errors.ParamsError)
		return
	}
	if err := service.UpdateRoleMenu(ctx, &in); err != nil {
		ctx.RespError(err)
	} else {
		ctx.RespSuccess()
	}
}

func RoleMenuIds(ctx *gin.Context) {
	in := types.RoleMenuIdsRequest{}
	if err := ctx.ShouldBind(&in); err != nil {
		ctx.RespError(errors.ParamsError)
		return
	}
	if ids, err := service.RoleMenuIds(ctx, &in); err != nil {
		ctx.RespError(err)
	} else {
		ctx.RespData(ids)
	}
}
