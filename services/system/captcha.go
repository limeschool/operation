package system

import (
	"context"
	"crypto/md5"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/limeschool/gin"
	"github.com/mojocn/base64Captcha"
	"operation/consts"
	"operation/tools"
	"time"
)

func Captcha(ctx *gin.Context) (any, error) {
	dt := base64Captcha.NewDriverDigit(80, 240, 4, 0.7, 80)
	duration := 30
	store := NewCaptchaStore(ctx, time.Duration(duration)*time.Second)
	captcha := base64Captcha.NewCaptcha(dt, store)
	id, base64, err := captcha.Generate()

	// ip+id做关联。
	store.ReStoreClientCaptcha(ctx, id)

	return map[string]any{
		"id":     id,
		"base64": base64,
		"expire": duration,
	}, err
}

type captchaStore struct {
	duration time.Duration
	redis    *redis.Client
}

func NewCaptchaStore(ctx *gin.Context, duration time.Duration) *captchaStore {
	return &captchaStore{
		duration: duration,
		redis:    ctx.Redis(consts.REDIS),
	}
}

// ReStoreClientCaptcha 只缓存一次用户的id
func (s *captchaStore) ReStoreClientCaptcha(ctx *gin.Context, id string) {
	uuid := fmt.Sprintf("captcha_%x", md5.Sum([]byte(tools.ClientIP(ctx))))
	oldId := s.Get(uuid, false)
	if oldId != "" {
		s.Clear(oldId)
	}
	_ = s.Set(uuid, id)
}

func (s *captchaStore) Set(id string, value string) error {
	return s.redis.Set(context.Background(), id, value, s.duration).Err()
}

func (s *captchaStore) Get(id string, clear bool) string {
	res := s.redis.Get(context.Background(), id).String()
	if clear {
		s.Clear(id)
	}
	return res
}

func (s *captchaStore) Clear(id string) {
	s.redis.Del(context.Background(), id)
}

func (s *captchaStore) Verify(id string, answer string, clear bool) bool {
	res, err := s.redis.Get(context.Background(), id).Result()
	if err != nil {
		return false
	}
	if clear {
		s.Clear(id)
	}
	return res == answer
}
