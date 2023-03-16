package config

import (
	"github.com/limeschool/gin"
	"github.com/spf13/viper"
	"operation/middlewares"
	"operation/tools/captcha"
	"operation/tools/upload"
)

// InitConfig 初始化并监听配置变更
func InitConfig() {
	gin.WatchConfigFunc(func(v *viper.Viper) {

		// 初始化jwt信息
		middlewares.NewJwt(v)

		// 初始化upload信息
		upload.InitUploadConfig(v)

		// 初始化验证器
		captcha.InitStore(v)
	})
}
