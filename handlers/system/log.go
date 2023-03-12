package system

import (
	"github.com/limeschool/gin"
	"operation/errors"
	service "operation/services/system"
	types "operation/types/system"
)

func LoginLog(ctx *gin.Context) {
	// 检验参数
	in := types.LoginLogRequest{}
	if ctx.ShouldBind(&in) != nil {
		ctx.RespError(errors.ParamsError)
		return
	}

	if list, total, err := service.PageLoginLog(ctx, &in); err != nil {
		ctx.RespError(err)
	} else {
		ctx.RespList(in.Page, len(list), int(total), list)
	}
}
