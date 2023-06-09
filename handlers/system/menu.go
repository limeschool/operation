package system

import (
	"github.com/limeschool/gin"
	"operation/errors"
	service "operation/services/system"
	types "operation/types/system"
)

func AllMenu(ctx *gin.Context) {
	if resp, err := service.AllMenu(ctx); err != nil {
		ctx.RespError(err)
	} else {
		ctx.RespData(resp)
	}
}

func AddMenu(ctx *gin.Context) {
	in := types.AddMenuRequest{}
	if err := ctx.ShouldBindJSON(&in); err != nil {
		ctx.RespError(errors.ParamsError)
		return
	}
	if err := service.AddMenu(ctx, &in); err != nil {
		ctx.RespError(err)
	} else {
		ctx.RespSuccess()
	}
}

func UpdateMenu(ctx *gin.Context) {
	in := types.UpdateMenuRequest{}
	if ctx.ShouldBindJSON(&in) != nil {
		ctx.RespError(errors.ParamsError)
		return
	}
	if err := service.UpdateMenu(ctx, &in); err != nil {
		ctx.RespError(err)
	} else {
		ctx.RespSuccess()
	}
}

func DeleteMenu(ctx *gin.Context) {
	in := types.DeleteMenuRequest{}
	if ctx.ShouldBind(&in) != nil {
		ctx.RespError(errors.ParamsError)
		return
	}
	if err := service.DeleteMenu(ctx, &in); err != nil {
		ctx.RespError(err)
	} else {
		ctx.RespSuccess()
	}
}
