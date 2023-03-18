package system

import (
	"crypto/md5"
	"fmt"
	"github.com/google/uuid"
	"github.com/limeschool/gin"
	"github.com/mojocn/base64Captcha"
	"operation/errors"
	"operation/middlewares/meta"
	model "operation/models/system"
	"operation/tools"
	ct "operation/tools/captcha"
	"operation/tools/email"
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
		"expire": ct.Store.DurationSecond(),
	}, err
}

func EmailCaptcha(ctx *gin.Context) (any, error) {
	md, err := meta.Get(ctx)
	if err != nil {
		return nil, err
	}

	id := fmt.Sprintf("%x", md5.Sum([]byte(uuid.New().String())))
	code := tools.CreateCode(email.Sender.Length)
	if ct.Store.Set(id, code) != nil {
		return nil, errors.New("存储验证码失败")
	}

	// 获取用户邮箱信息
	user := model.User{}
	if err = user.OneByID(ctx, md.UserID); err != nil {
		return nil, err
	}

	// 组装邮件信息
	sender := email.Sender.New(email.DefaultTemplate)
	str := "您的邮箱验证码为：%s，该验证码%v分钟内有效，为了保证您的账户安全，请勿向他人泄露验证码吗信息"
	content := fmt.Sprintf(str, code, ct.Store.DurationMinute())
	title := "验证码通知"

	// 发送邮件
	if sender.Send(user.Email, title, content) != nil {
		return nil, errors.CaptchaSendError
	}

	ct.Store.ReStoreClientCaptcha(tools.ClientIP(ctx), id)

	return map[string]any{
		"id":     id,
		"expire": ct.Store.DurationSecond(),
	}, nil

}
