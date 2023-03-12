package system

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/limeschool/gin"
	"github.com/mojocn/base64Captcha"
	"operation/consts"
	"time"
)

func Captcha(ctx *gin.Context) (any, error) {
	dt := base64Captcha.NewDriverDigit(80, 240, 4, 0.7, 80)
	captcha := base64Captcha.NewCaptcha(dt, NewCaptchaStore(ctx))
	id, base64, err := captcha.Generate()
	return map[string]string{
		"id":     id,
		"base64": base64,
	}, err
}

type captchaStore struct {
	Redis *redis.Client
}

func NewCaptchaStore(ctx *gin.Context) *captchaStore {
	return &captchaStore{
		Redis: ctx.Redis(consts.REDIS),
	}
}

func (s *captchaStore) Set(id string, value string) error {
	return s.Redis.Set(context.Background(), id, value, 10*time.Minute).Err()
}

func (s *captchaStore) Get(id string, clear bool) string {
	return s.Redis.Get(context.Background(), id).String()
}

func (s *captchaStore) Verify(id string, answer string, clear bool) bool {
	res, err := s.Redis.Get(context.Background(), id).Result()
	if err != nil {
		return false
	}
	if clear {
		s.Redis.Del(context.Background(), id)
	}
	return res == answer
}
