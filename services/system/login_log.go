package system

import (
	"github.com/limeschool/gin"
	"operation/models"
	model "operation/models/system"
	"operation/tools/address"
	"operation/tools/ua"
	types "operation/types/system"
)

func AddLoginLog(ctx *gin.Context, phone string, err error) error {
	ip := ctx.RemoteIP()
	userAgent := ctx.Request.Header.Get("User-Agent")
	info := ua.Parse(userAgent)
	desc := ""
	code := 0

	if err != nil {
		customErr, _ := err.(*gin.CustomError)
		code = customErr.Code
		desc = customErr.Msg
	}

	log := model.LoginLog{
		Phone:       phone,
		IP:          ip,
		Address:     address.GetAddress(ip),
		Browser:     info.Name,
		Status:      err == nil,
		Description: desc,
		Code:        code,
		Device:      info.OS + " " + info.OSVersion,
	}
	return log.Create(ctx)
}

func PageLoginLog(ctx *gin.Context, in *types.LoginLogRequest) ([]model.LoginLog, int64, error) {
	log := model.LoginLog{}
	return log.Page(ctx, models.PageOptions{
		Page:  in.Page,
		Count: in.Count,
		Model: in,
	})
}
