package system

import (
	"github.com/limeschool/gin"
	"operation/middlewares/meta"
	"operation/models"
	"operation/tools/tree"
	"time"
)

type Role struct {
	gin.BaseModel
	ParentID    int64   `json:"parent_id"`
	Name        string  `json:"name" `
	Keyword     string  `json:"keyword"`
	Status      *bool   `json:"status,omitempty" `
	Weight      *int    `json:"weight"`
	Description *string `json:"description,omitempty"`
	TeamIds     *string `json:"team_ids,omitempty"`
	DataScope   string  `json:"data_scope,omitempty"`
	Operator    string  `json:"operator"`
	OperatorID  int64   `json:"operator_id"`
	Children    []*Role `json:"children"  gorm:"-"`
}

func (r *Role) ID() int64 {
	return r.BaseModel.ID
}

func (r *Role) Parent() int64 {
	return r.ParentID
}

func (r *Role) AppendChildren(child any) {
	menu := child.(*Role)
	r.Children = append(r.Children, menu)
}

func (r *Role) ChildrenNode() []tree.Tree {
	var list []tree.Tree
	for _, item := range r.Children {
		list = append(list, item)
	}
	return list
}

func (r *Role) Table() string {
	return "tb_system_role"
}

// Create 创建角色信息
func (r *Role) Create(ctx *gin.Context) error {
	// 操作者信息
	md, err := meta.Get(ctx)
	if err != nil {
		return err
	}

	r.OperatorID = md.UserID
	r.Operator = md.Username
	r.UpdatedAt = time.Now().Unix()
	return models.TransferErr(models.Database(ctx).Table(r.Table()).Create(&r).Error)
}

// OneByID 通过ID查询角色信息
func (r *Role) OneByID(ctx *gin.Context, id int64) error {
	return models.TransferErr(models.Database(ctx).Table(r.Table()).First(r, "id = ?", id).Error)
}

// All 查询全部角色信息
func (r *Role) All(ctx *gin.Context, cond ...any) ([]*Role, error) {
	var list []*Role
	db := models.Database(ctx).Table(r.Table()).Order("weight desc")
	return list, models.TransferErr(db.Find(&list, cond...).Error)
}

func (r *Role) RoleStatus(ctx *gin.Context, roleId int64) bool {
	team, err := r.Tree(ctx, 1)
	if err != nil {
		return false
	}
	res := false
	dfsRoleStatus(team.(*Role), roleId, true, &res)
	return res
}

func dfsRoleStatus(role *Role, roleId int64, status bool, res *bool) bool {
	if roleId == role.BaseModel.ID {
		is := *role.Status && status
		*res = is
	}

	for _, item := range role.Children {
		dfsRoleStatus(item, roleId, status && *item.Status, res)
	}

	return status
}

// Tree 查询全部角色树
func (r *Role) Tree(ctx *gin.Context, roleId int64) (tree.Tree, error) {
	list, err := r.All(ctx)
	if err != nil {
		return nil, err
	}
	var treeList []tree.Tree
	for _, item := range list {
		treeList = append(treeList, item)
	}
	return tree.BuildTreeByID(treeList, roleId), nil
}

// Update 更新角色信息
func (r *Role) Update(ctx *gin.Context) error {
	md, err := meta.Get(ctx)
	if err != nil {
		return err
	}

	r.OperatorID = md.UserID
	r.Operator = md.Username
	return models.TransferErr(models.Database(ctx).Table(r.Table()).Updates(r).Error)
}

// DeleteByID 通过ID删除角色信息
func (r *Role) DeleteByID(ctx *gin.Context, id int64) error {
	return models.TransferErr(models.Database(ctx).Table(r.Table()).Where("id = ?", id).Delete(&r).Error)
}
