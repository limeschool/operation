package system

import (
	"encoding/json"
	"github.com/limeschool/gin"
	"github.com/zeromicro/go-zero/core/syncx"
	"operation/consts"
	"operation/middlewares/meta"
	"operation/models"
	"operation/tools/tree"
	"time"
)

type Menu struct {
	ParentID   int64   `json:"parent_id"`
	Title      string  `json:"title"`
	Icon       string  `json:"icon"`
	Path       string  `json:"path"`
	Name       string  `json:"name"`
	Type       string  `json:"type"`
	Permission string  `json:"permission"`
	Method     string  `json:"method"`
	Component  string  `json:"component"`
	Redirect   *string `json:"redirect"`
	Weight     *int    `json:"weight"`
	IsHidden   *bool   `json:"is_hidden"`
	IsCache    *bool   `json:"is_cache"`
	Operator   string  `json:"operator"`
	OperatorID int64   `json:"operator_id"`
	Children   []*Menu `json:"children,omitempty" gorm:"-"`
	gin.BaseModel
}

func (u *Menu) ID() int64 {
	return u.BaseModel.ID
}

func (u *Menu) Parent() int64 {
	return u.ParentID
}

func (u *Menu) AppendChildren(child any) {
	menu := child.(*Menu)
	u.Children = append(u.Children, menu)
}

func (u *Menu) ChildrenNode() []tree.Tree {
	var list []tree.Tree
	for _, item := range u.Children {
		list = append(list, item)
	}
	return list
}

func (u *Menu) Table() string {
	return "tb_system_menu"
}

// Create 创建菜单
func (u *Menu) Create(ctx *gin.Context) error {
	md, err := meta.Get(ctx)
	if err != nil {
		return err
	}

	u.OperatorID = md.UserID
	u.Operator = md.Username
	u.UpdatedAt = time.Now().Unix()
	// 创建菜单为基础api时删除缓存
	if u.Permission == consts.BaseApi {
		models.DelayDelCache(ctx, consts.RedisBaseApi)
	}
	// 创建菜单
	return models.TransferErr(models.Database(ctx).Table(u.Table()).Create(&u).Error)
}

// OneByID 通过id查询指定菜单
func (u *Menu) OneByID(ctx *gin.Context, id int64) error {
	return models.TransferErr(models.Database(ctx).Table(u.Table()).First(u, id).Error)
}

// OneByCond 通过条件查询指定菜单
func (u *Menu) OneByCond(ctx *gin.Context, cond ...interface{}) error {
	return models.TransferErr(models.Database(ctx).Table(u.Table()).First(u, cond...).Error)
}

// GetBaseApiPath 获取基础菜单api列表
func (u *Menu) GetBaseApiPath(ctx *gin.Context) map[string]bool {
	lockKey := consts.RedisBaseApi + "_lock"
	flight := syncx.NewSingleFlight()
	resp, err := flight.Do(lockKey, func() (interface{}, error) {
		resp := map[string]bool{}
		// 从缓存中获取数据
		str, err := models.Cache(ctx).Get(ctx, consts.RedisBaseApi).Result()
		if err == nil && str != "" && json.Unmarshal([]byte(str), &resp) == nil {
			return resp, nil
		}

		// 缓存中读取失败，则查询数据库
		list, err := u.All(ctx, "permission = ? and type = 'A'", consts.BaseApi)
		if err != nil {
			return nil, err
		}

		for _, val := range list {
			resp[val.Method+":"+val.Path] = true
		}

		// 将数据库中的数据存入缓存
		byteData, _ := json.Marshal(resp)
		ctx.Redis(consts.REDIS).Set(ctx, consts.RedisBaseApi, string(byteData), 1*time.Hour)
		return resp, nil
	})

	if err != nil {
		return nil
	}
	return resp.(map[string]bool)
}

// All 获取全部的菜单列表
func (u *Menu) All(ctx *gin.Context, cond ...interface{}) ([]*Menu, error) {
	var list []*Menu
	db := models.Database(ctx).Table(u.Table())
	return list, models.TransferErr(db.Order("weight desc").Find(&list, cond...).Error)
}

// Tree 获取菜单树
func (u *Menu) Tree(ctx *gin.Context, cond ...interface{}) (tree.Tree, error) {
	list, err := u.All(ctx, cond...)
	if err != nil {
		return nil, err
	}
	var treeList []tree.Tree
	for _, item := range list {
		treeList = append(treeList, item)
	}
	return tree.BuildTree(treeList), nil
}

// Update 更新菜单
func (u *Menu) Update(ctx *gin.Context) error {
	md, err := meta.Get(ctx)
	if err != nil {
		return err
	}

	u.OperatorID = md.UserID
	u.Operator = md.Username

	// 删除基础api缓存
	if u.Permission == consts.BaseApi {
		models.DelayDelCache(ctx, consts.RedisBaseApi)
	}

	return models.TransferErr(models.Database(ctx).Table(u.Table()).Updates(u).Error)
}

// DeleteByCond 通过条件删除菜单
func (u *Menu) DeleteByCond(ctx *gin.Context, cond ...interface{}) error {
	if err := models.Database(ctx).Table(u.Table()).First(u, cond...).Error; err != nil {
		return models.TransferErr(err)
	}

	// 删除基础api缓存
	if u.Permission == consts.BaseApi {
		models.DelayDelCache(ctx, consts.RedisBaseApi)
	}
	return models.TransferErr(models.Database(ctx).Table(u.Table()).Delete(u).Error)
}
