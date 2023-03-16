package system

import (
	"github.com/limeschool/gin"
	"github.com/mojocn/base64Captcha"
	"operation/tools"
	ct "operation/tools/captcha"
)

func Captcha(ctx *gin.Context) (any, error) {
	dt := base64Captcha.NewDriverDigit(80, 240, 4, 0.7, 80)

	captcha := base64Captcha.NewCaptcha(dt, ct.Store)
	id, base64, err := captcha.Generate()

	// ip+id做关联。
	ct.Store.ReStoreClientCaptcha(tools.ClientIP(ctx), id)

	return map[string]any{
		"id":     id,
		"base64": base64,
		"expire": ct.Store.Duration(),
	}, err
}
