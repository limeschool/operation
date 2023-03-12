package system

import (
	"github.com/limeschool/gin"
	service "operation/services/system"
)

func Captcha(ctx *gin.Context) {
	if resp, err := service.Captcha(ctx); err != nil {
		ctx.RespError(err)
	} else {
		ctx.RespData(resp)
	}
}
