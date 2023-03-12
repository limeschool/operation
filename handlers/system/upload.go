package system

import (
	"github.com/limeschool/gin"
	"operation/errors"
	service "operation/services/system"
	types "operation/types/system"
)

func UploadFile(ctx *gin.Context) {
	// 检验参数
	in := types.UploadRequest{}
	if ctx.ShouldBind(&in) != nil {
		ctx.RespError(errors.ParamsError)
		return
	}

	if data, err := service.UploadFile(ctx, &in); err != nil {
		ctx.RespError(err)
	} else {
		ctx.RespData(data)
	}
}
