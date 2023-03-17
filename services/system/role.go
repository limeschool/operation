package system

import (
	"github.com/jinzhu/copier"
	"github.com/limeschool/gin"
	"operation/errors"
	model "operation/models/system"
	"operation/tools/tree"
	types "operation/types/system"
)

// AllRole 返回所有的角色
func AllRole(ctx *gin.Context) (tree.Tree, error) {
	// 获取当前用户的角色
	user, err := CurUser(ctx)
	if err != nil {
		return nil, err
	}
	if err = user.OneByID(ctx, user.ID); err != nil {
		return nil, err
	}

	role := model.Role{}
	return role.Tree(ctx, user.RoleID)
}

// AddRole 新增角色
func AddRole(ctx *gin.Context, in *types.AddRoleRequest) error {
	role := model.Role{}
	if copier.Copy(&role, in) != nil {
		return errors.AssignError
	}
	return role.Create(ctx)
}

// UpdateRole 更新角色信息
func UpdateRole(ctx *gin.Context, in *types.UpdateRoleRequest) error {
	if in.ID == 1 {
		return errors.SuperAdminEditError
	}

	if in.Status != nil && !*in.Status {
		user, err := CurUser(ctx)
		if err != nil {
			return err
		}
		if in.ID == user.RoleID {
			return errors.DisableCurRoleError
		}
	}

	role := model.Role{}
	if copier.Copy(&role, in) != nil {
		return errors.AssignError
	}
	return role.Update(ctx)
}

// DeleteRole 删除角色信息
func DeleteRole(ctx *gin.Context, in *types.DeleteRoleRequest) error {
	// 超级管理员不允许删除
	if in.ID == 1 {
		return errors.SuperAdminDelError
	}

	// 删除角色时需要删除rbac权限表
	role := model.Role{}
	if err := role.OneByID(ctx, in.ID); err != nil {
		return err
	}
	_, _ = ctx.Rbac().Object().RemoveFilteredPolicy(0, role.Keyword)

	return role.DeleteByID(ctx, in.ID)
}
