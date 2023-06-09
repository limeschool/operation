package system

import (
	"github.com/limeschool/gin"
	"operation/errors"
	model "operation/models/system"
	types "operation/types/system"
)

// UpdateRoleMenu 修改角色所属菜单
func UpdateRoleMenu(ctx *gin.Context, in *types.AddRoleMenuRequest) error {
	// 超级管理员不存在菜单权限，自动获取全部菜单
	if in.RoleID == 1 {
		return errors.SuperAdminEditError
	}

	// 获取当前role的数据
	role := model.Role{}
	if err := role.OneByID(ctx, in.RoleID); err != nil {
		return err
	}

	// 进行菜单修改
	rm := model.RoleMenu{}
	if err := rm.Update(ctx, in.RoleID, in.MenuIds); err != nil {
		return err
	}

	// 删除当前用户的全部rbac权限
	_, _ = ctx.Rbac().RemoveFilteredPolicy(0, role.Keyword)

	// 获取当前修改菜单的信息
	menu := model.Menu{}
	var policies [][]string
	apiList, _ := menu.All(ctx, "id in ? and type = 'A'", in.MenuIds)
	for _, item := range apiList {
		policies = append(policies, []string{role.Keyword, item.Path, item.Method})
	}

	// 将新的策略的策略写入rbac
	_, _ = ctx.Rbac().AddPolicies(policies)

	return nil
}

// RoleMenuIds 获取角色菜单的所有id
func RoleMenuIds(ctx *gin.Context, in *types.RoleMenuIdsRequest) ([]int64, error) {

	// 获取当前角色的所有菜单
	rm := model.RoleMenu{}
	rmList, err := rm.RoleMenus(ctx, in.RoleID)
	if err != nil {
		return nil, err
	}

	// 组装所有的菜单id
	var ids []int64
	for _, item := range rmList {
		ids = append(ids, item.MenuID)
	}

	return ids, nil
}
