package system

import (
	"github.com/limeschool/gin"
	"gorm.io/gorm"
	"operation/middlewares/meta"
	"operation/models"
)

type RoleMenu struct {
	gin.BaseModel
	RoleID     int64  `json:"role_id"`
	MenuID     int64  `json:"menu_id"`
	Operator   string `json:"operator"`
	OperatorID int64  `json:"operator_id"`
}

func (RoleMenu) Table() string {
	return "tb_system_role_menu"
}

// Update 批量更新角色所属菜单
func (u *RoleMenu) Update(ctx *gin.Context, roleId int64, menuIds []int64) error {
	// 操作者信息
	md, err := meta.Get(ctx)
	if err != nil {
		return err
	}

	// 组装新的菜单数据
	list := make([]RoleMenu, 0)
	for _, menuId := range menuIds {
		list = append(list, RoleMenu{
			RoleID:     roleId,
			MenuID:     menuId,
			OperatorID: md.UserID,
			Operator:   md.Username,
		})
	}

	err = models.Database(ctx).Table(u.Table()).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("role_id=?", roleId).Delete(u).Error; err != nil {
			return err
		}
		return tx.Create(&list).Error
	})

	return models.TransferErr(err)
}

// RoleMenus 通过角色ID获取角色菜单
func (u *RoleMenu) RoleMenus(ctx *gin.Context, roleId int64) ([]RoleMenu, error) {
	var list []RoleMenu
	db := models.Database(ctx).Table(u.Table())
	return list, models.TransferErr(db.Find(&list, "role_id=?", roleId).Error)
}

// MenuRoles 通过菜单ID获取角色菜单列表
func (u *RoleMenu) MenuRoles(ctx *gin.Context, menuId int64) ([]RoleMenu, error) {
	var list []RoleMenu
	db := models.Database(ctx).Table(u.Table())
	return list, models.TransferErr(db.Find(&list, "menu_id=?", menuId).Error)
}

// DeleteByRoleID 通过角色id删除 角色所属菜单
func (u *RoleMenu) DeleteByRoleID(ctx *gin.Context, roleId int64) error {
	return models.TransferErr(models.Database(ctx).Table(u.Table()).Delete(u, "role_id=?", roleId).Error)
}
