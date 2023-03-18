package captcha

import (
	"context"
	"crypto/md5"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/limeschool/gin"
	"github.com/spf13/viper"
	"time"
)

var (
	Store *captchaStore
)

type captchaStore struct {
	duration time.Duration
	redis    *redis.Client
}

func InitStore(v *viper.Viper) {
	ctx := gin.NewContext()

	duration := v.GetDuration("captcha.expire")
	cache := ctx.Redis(v.GetString("captcha.store"))

	if duration == 0 || cache == nil {
		panic("验证器初始化失败")
	}
	Store = &captchaStore{
		duration: duration,
		redis:    cache,
	}
}

// ReStoreClientCaptcha 只缓存一次用户的id
func (s *captchaStore) ReStoreClientCaptcha(ip, id string) {
	uuid := fmt.Sprintf("captcha_%x", md5.Sum([]byte(ip)))
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

func (s *captchaStore) DurationSecond() int {
	return int(s.duration / time.Second)
}

func (s *captchaStore) DurationMinute() int {
	return int(s.duration / time.Minute)
}
