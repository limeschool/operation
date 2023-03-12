package system

import (
	"github.com/limeschool/gin"
	"operation/models"
)

type LoginLog struct {
	gin.CreateModel
	Phone       string `json:"phone"`
	IP          string `json:"ip"`
	Address     string `json:"address"`
	Browser     string `json:"browser"`
	Device      string `json:"device"`
	Status      bool   `json:"status"`
	Code        int    `json:"code"`
	Description string `json:"description"`
}

func (u LoginLog) Table() string {
	return "tb_system_login_log"
}

func (u *LoginLog) Create(ctx *gin.Context) error {
	return models.TransferErr(models.Database(ctx).Table(u.Table()).Create(u).Error)
}

func (u *LoginLog) Page(ctx *gin.Context, options models.PageOptions) ([]LoginLog, int64, error) {
	list, total := make([]LoginLog, 0), int64(0)

	db := models.Database(ctx).Table(u.Table())

	if options.Model != nil {
		db = gin.GormWhere(db, u.Table(), options.Model)
	}

	if options.Scopes != nil {
		db = db.Scopes(options.Scopes)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	db = db.Order("created_at").Offset((options.Page - 1) * options.Count).Limit(options.Count)

	return list, total, db.Find(&list).Error
}
