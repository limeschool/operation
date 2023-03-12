package system

import (
	"encoding/json"
	"github.com/limeschool/gin"
	"gorm.io/gorm"
	"operation/consts"
	"operation/middlewares/meta"
	"operation/models"
	"operation/tools"
	"operation/tools/tree"
	"time"
)

type User struct {
	gin.BaseModel
	TeamID      int64   `json:"team_id"`
	RoleID      int64   `json:"role_id"`
	Name        string  `json:"name"`
	Nickname    string  `json:"nickname"`
	Sex         *bool   `json:"sex,omitempty"`
	Phone       string  `json:"phone"`
	Password    string  `json:"password,omitempty"  gorm:"->:false;<-:create,update"`
	Avatar      string  `json:"avatar"`
	Email       string  `json:"email,omitempty"`
	Status      *bool   `json:"status,omitempty"`
	DisableDesc *string `json:"disable_desc"`
	LastLogin   int64   `json:"last_login"`
	Operator    string  `json:"operator"`
	OperatorID  int64   `json:"operator_id"`
	Role        Role    `json:"role" gorm:"->;foreignKey:role_id;reference:id"`
	Team        Team    `json:"team" gorm:"->;foreignKey:team_id;reference:id"`
}

func (u User) Table() string {
	return "tb_system_user"
}

// OneByID 通过id查询用户信息
func (u *User) OneByID(ctx *gin.Context, id int64) error {

	db := models.Database(ctx).Table(u.Table())

	db = db.Preload("Role", func(db *gorm.DB) *gorm.DB {
		return db.Table(u.Role.Table())
	})

	db = db.Preload("Team", func(db *gorm.DB) *gorm.DB {
		return db.Table(u.Role.Table())
	})

	return models.TransferErr(db.First(u, id).Error)
}

// OneByPhone 通过phone查询用户信息
func (u *User) OneByPhone(ctx *gin.Context, phone string) error {

	db := models.Database(ctx).Table(u.Table())

	db = db.Preload("Role", func(db *gorm.DB) *gorm.DB {
		return db.Table(u.Role.Table())
	})

	db = db.Preload("Team", func(db *gorm.DB) *gorm.DB {
		return db.Table(u.Role.Table())
	})

	return models.TransferErr(db.First(u, "phone=?", phone).Error)
}

// PasswordByPhone 查询全部字段信息包括密码
func (u *User) PasswordByPhone(ctx *gin.Context, phone string) (string, error) {
	m := map[string]any{}

	db := models.Database(ctx).Table(u.Table())

	if err := db.First(u, "phone = ?", phone).Scan(&m).Error; err != nil {
		return "", models.TransferErr(err)
	}

	return m["password"].(string), nil
}

// Page 查询分页数据
func (u *User) Page(ctx *gin.Context, options models.PageOptions) ([]*User, int64, error) {
	list, total := make([]*User, 0), int64(0)

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

	db = db.Preload("Role", func(db *gorm.DB) *gorm.DB {
		return db.Table(u.Role.Table())
	})

	db = db.Preload("Team", func(db *gorm.DB) *gorm.DB {
		return db.Table(u.Role.Table())
	})

	db = db.Offset((options.Page - 1) * options.Count).Limit(options.Count)

	return list, total, db.Find(&list).Error
}

// Create 创建用户信息
func (u *User) Create(ctx *gin.Context) error {
	// 操作者信息
	md, err := meta.Get(ctx)
	if err != nil {
		return err
	}

	u.OperatorID = md.UserID
	u.Operator = md.Username
	u.UpdatedAt = time.Now().Unix()
	u.Password, _ = tools.ParsePwd(u.Password)
	return models.TransferErr(models.Database(ctx).Table(u.Table()).Create(u).Error)
}

func (u *User) UpdateLastLogin(ctx *gin.Context, t int64) error {
	db := models.Database(ctx).Table(u.Table())
	return models.TransferErr(db.Where("id", u.ID).Update("last_login", t).Error)
}

// Update 更新用户信息
func (u *User) Update(ctx *gin.Context) error {
	// 操作者信息
	md, err := meta.Get(ctx)
	if err != nil {
		return err
	}

	u.OperatorID = md.UserID
	u.Operator = md.Username
	if u.Password != "" {
		u.Password, _ = tools.ParsePwd(u.Password)
	}

	// 执行更新
	return models.TransferErr(models.Database(ctx).Table(u.Table()).Updates(&u).Error)
}

// DeleteByID 通过id删除用户信息
func (u *User) DeleteByID(ctx *gin.Context, id int64) error {
	return models.TransferErr(models.Database(ctx).Table(u.Table()).Delete(u, id).Error)
}

// GetTeamIdsByID 通过用户id获取用户所管理的部门id
func (u *User) GetTeamIdsByID(ctx *gin.Context, userId int64) ([]int64, error) {
	// 操作者信息
	user := User{}
	if err := user.OneByID(ctx, userId); err != nil {
		return nil, err
	}

	// 查询角色信息
	role := Role{}
	if err := role.OneByID(ctx, user.RoleID); err != nil {
		return nil, err
	}

	// 当用户权限是当前部门时，直接返回当前部门的id
	if role.DataScope == consts.CURTEAM {
		return []int64{user.TeamID}, nil
	}

	ids := make([]int64, 0)
	// 当用户权限是自定义部门时，直接返回自定义部门id
	if role.DataScope == consts.CUSTOM {
		return ids, json.Unmarshal([]byte(*role.TeamIds), &ids)
	}

	// 以当前部门为根节点构造部门树
	team := Team{}
	teamList, _ := team.All(ctx)
	var treeList []tree.Tree
	for _, item := range teamList {
		treeList = append(treeList, item)
	}
	teamTree := tree.BuildTreeByID(treeList, user.TeamID)

	// 根据部门树取值
	switch role.DataScope {
	case consts.ALLTEAM:
		// 全部数据权限时返回所有部门id
		ids = tree.GetTreeID(teamTree)
	case consts.DOWNTEAM:
		// 下级部门权限时，排除当前部门id
		ids = tree.GetTreeID(teamTree)
		if len(ids) > 2 {
			ids = ids[1:]
		} else {
			ids = []int64{}
		}
	}
	return ids, nil
}

// CurUserTeamIds 获取当前用户的部门管理ID
func (u *User) CurUserTeamIds(ctx *gin.Context) ([]int64, error) {
	// 操作者信息
	user, err := meta.Get(ctx)
	if err != nil {
		return nil, err
	}

	// 查询角色信息
	role := Role{}
	if err = role.OneByID(ctx, user.RoleID); err != nil {
		return nil, err
	}

	// 当用户权限是当前部门时，直接返回当前部门的id
	if role.DataScope == consts.CURTEAM {
		return []int64{user.TeamID}, nil
	}

	ids := make([]int64, 0)
	// 当用户权限是自定义部门时，直接返回自定义部门id
	if role.DataScope == consts.CUSTOM {
		return ids, json.Unmarshal([]byte(*role.TeamIds), &ids)
	}

	// 以当前部门为根节点构造部门树
	team := Team{}
	teamList, _ := team.All(ctx)
	var treeList []tree.Tree
	for _, item := range teamList {
		treeList = append(treeList, item)
	}
	teamTree := tree.BuildTreeByID(treeList, user.TeamID)

	// 根据部门树取值
	switch role.DataScope {
	case consts.ALLTEAM:
		// 全部数据权限时返回所有部门id
		ids = tree.GetTreeID(teamTree)
	case consts.DOWNTEAM:
		// 下级部门权限时，排除当前部门id
		ids = tree.GetTreeID(teamTree)
		if len(ids) > 2 {
			ids = ids[1:]
		} else {
			ids = []int64{}
		}
	}
	return ids, nil
}
