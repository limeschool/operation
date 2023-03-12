package models

import (
	"github.com/go-redis/redis/v8"
	"github.com/limeschool/gin"
	"gorm.io/gorm"
	"operation/consts"
	"operation/errors"
	"strings"
	"time"
)

type PageOptions struct {
	Page   int
	Count  int
	Model  any
	Scopes func(db *gorm.DB) *gorm.DB
}

type AllOptions struct {
	Model  any
	Scopes func(db *gorm.DB) *gorm.DB
}

// dataMap 数据字典
var dataMap = map[string]string{
	"phone":   "手机号码",
	"email":   "电子邮箱",
	"keyword": "标志",
	"name":    "名称",
}

// TransferErr 将数据库的错误转换成中文
func TransferErr(err error) error {
	if err == nil {
		return nil
	}

	if customErr, ok := err.(*gin.CustomError); ok {
		return customErr
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.DBNotFoundError
	}

	if strings.Contains(err.Error(), "Duplicate") {
		str := err.Error()
		str = strings.ReplaceAll(str, "'", "")
		str = strings.TrimPrefix(str, "Error 1062: Duplicate entry ")
		arr := strings.Split(str, " for key ")
		return errors.NewF(`%v "%v" 已存在`, dataMap[arr[1]], arr[0])
	}

	if strings.Contains(err.Error(), "FOREIGN KEY") {
		return errors.NewF(`数据正在被使用中，无法删除`)
	}

	return errors.DBError
}

// Database 进行数据库选择的快捷方法
func Database(ctx *gin.Context) *gorm.DB {
	return ctx.Orm(consts.DATABASE)
}

// ExecCallback 执行数据库回调的快捷方法
func ExecCallback(db *gorm.DB, fs ...func(db *gorm.DB) *gorm.DB) *gorm.DB {
	if fs != nil {
		for _, f := range fs {
			db = f(db)
		}
	}
	return db
}

// DelayDelCache 数据延迟双删
func DelayDelCache(ctx *gin.Context, key string) {
	ctx.Redis(consts.REDIS).Del(ctx, key)
	go func() {
		time.Sleep(1 * time.Second)
		ctx.Redis(consts.REDIS).Del(ctx, key)
	}()
}

// Cache 进行数据缓存
func Cache(ctx *gin.Context) *redis.Client {
	return ctx.Redis(consts.REDIS)
}
