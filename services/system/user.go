package system

import (
	"encoding/json"
	"github.com/jinzhu/copier"
	"github.com/limeschool/gin"
	"gorm.io/gorm"
	"operation/consts"
	"operation/errors"
	"operation/middlewares"
	"operation/middlewares/meta"
	"operation/models"
	model "operation/models/system"
	"operation/tools"
	"operation/tools/tree"
	types "operation/types/system"
	"time"
)

func CurTeamIds(ctx *gin.Context) ([]int64, error) {
	user := model.User{}
	teamIds, err := user.CurUserTeamIds(ctx)
	if err != nil {
		return nil, err
	}
	return teamIds, nil
}

func CurUser(ctx *gin.Context) (*model.User, error) {
	md, err := meta.Get(ctx)
	if err != nil {
		return nil, err
	}
	user := model.User{}
	return &user, user.OneByID(ctx, md.UserID)
}

func PageUser(ctx *gin.Context, in *types.PageUserRequest) ([]*model.User, int64, error) {
	user := model.User{}
	// 获取用户所管理的部门
	teamIds, err := user.CurUserTeamIds(ctx)
	if err != nil {
		return nil, 0, err
	}

	return user.Page(ctx, models.PageOptions{
		Page:  in.Page,
		Count: in.Count,
		Model: in,
		Scopes: func(db *gorm.DB) *gorm.DB {
			return db.Where("team_id in ?", teamIds)
		},
	})
}

func GetUser(ctx *gin.Context, in *types.GetUserRequest) (*model.User, error) {
	user := model.User{}
	return &user, user.OneByID(ctx, in.ID)
}

func AddUser(ctx *gin.Context, in *types.AddUserRequest) error {
	user := model.User{}
	if in.Nickname == "" {
		in.Nickname = in.Name
	}

	if copier.Copy(&user, in) != nil {
		return errors.AssignError
	}

	// 获取用户所管理的部门
	teamIds, err := user.CurUserTeamIds(ctx)
	if err != nil {
		return err
	}

	// 添加用户时，只允许添加当前所属部门的用户
	if !tools.InList(teamIds, in.TeamID) {
		return errors.NotAddTeamUserError
	}

	return user.Create(ctx)
}

func UpdateUser(ctx *gin.Context, in *types.UpdateUserRequest) error {

	user := model.User{}
	if user.OneByID(ctx, in.ID) != nil {
		return errors.DBNotFoundError
	}

	//超级管理员不允许修改所在部门和角色
	if in.ID == 1 {
		in.RoleID = 0
		in.TeamID = 0
		if *user.Status != *in.Status {
			return errors.SuperAdminEditError
		}
	}

	// 获取用户所管理的部门
	teamIds, err := user.CurUserTeamIds(ctx)
	if err != nil {
		return err
	}

	// 只允许更新当前部门的用户信息
	if !tools.InList(teamIds, user.TeamID) {
		return errors.NotEditTeamUserError
	}

	// 修改部门时，也只允许修改到自己所管辖的部门
	if in.TeamID != 0 && in.TeamID != user.TeamID && !tools.InList(teamIds, in.TeamID) {
		return errors.NotAddTeamUserError
	}

	if copier.Copy(&user, in) != nil {
		return errors.AssignError
	}

	return user.Update(ctx)
}

func DeleteUser(ctx *gin.Context, in *types.DeleteUserRequest) error {
	// 超级管理员不允许删除
	if in.ID == 1 {
		return errors.SuperAdminDelError
	}

	user := model.User{}
	if user.OneByID(ctx, in.ID) != nil {
		return errors.DBNotFoundError
	}

	teamIds, err := user.CurUserTeamIds(ctx)
	if err != nil {
		return err
	}

	// 只允许删除当前所管理部门的人员
	if !tools.InList(teamIds, user.TeamID) {
		return errors.NotDelTeamUserError
	}

	return user.DeleteByID(ctx, in.ID)
}

func UserLogout(ctx *gin.Context) error {
	metadata, err := meta.Get(ctx)
	if err != nil {
		return err
	}
	return middlewares.JwtAuth.ClearToken(ctx, metadata.UserID)
}

func UserLogin(ctx *gin.Context, in *types.UserLoginRequest) (resp *types.UserLoginResponse, err error) {
	// 记录登陆日志
	resp = new(types.UserLoginResponse)
	defer func() {
		if !(errors.Is(err, errors.UserDisableError) ||
			errors.Is(err, errors.CaptchaError)) {
			_ = AddLoginLog(ctx, in.Phone, err)
		}
	}()

	if middlewares.JwtAuth.Captcha {
		if !NewCaptchaStore(ctx).Verify(in.CaptchaID, in.Captcha, true) {
			err = errors.CaptchaError
			return
		}
	}

	// 用户密码解密
	if in.Password, err = ctx.Rsa(consts.RsaPrivate).Decode(in.Password); err != nil {
		err = errors.RsaPasswordError
		return
	}

	var pw types.Password
	if json.Unmarshal([]byte(in.Password), &pw) != nil {
		err = errors.RsaPasswordError
		return
	}

	// 判断当前时间戳是否过期,超过10s则拒绝
	if time.Now().UnixMilli()-pw.Time > 10*1000 {
		err = errors.PasswordExpireError
		return
	}

	in.Password = pw.Password

	// 通过手机号获取用户信息
	user := model.User{}
	if err = user.OneByPhone(ctx, in.Phone); err != nil {
		return
	}

	// 由于屏蔽了password，需要调用指定方法查询
	password, err := user.PasswordByPhone(ctx, in.Phone)
	if err != nil {
		err = errors.UserNotFoundError
		return
	}

	// 用户被禁用则拒绝登陆
	if !*user.Status {
		err = errors.UserDisableError
		return
	}

	// 所属角色被禁用则拒绝登陆
	role := model.Role{}
	if !role.RoleStatus(ctx, user.RoleID) {
		return nil, errors.RoleDisableError
	}

	// 对比用户密码，错误则拒绝登陆
	if !tools.CompareHashPwd(password, in.Password) {
		err = errors.PasswordError
		return
	}

	// 生成登陆token
	if resp.Token, err = middlewares.JwtAuth.CreateToken(ctx, &meta.Metadata{
		UserID:    user.ID,
		RoleID:    user.RoleID,
		RoleKey:   user.Role.Keyword,
		DataScope: user.Role.DataScope,
		Username:  user.Name,
		TeamID:    user.TeamID,
	}); err != nil {
		return nil, err
	}

	// 将用户的token信息写入redis

	// 修改登陆时间
	return resp, user.UpdateLastLogin(ctx, time.Now().Unix())
}

func RefreshToken(ctx *gin.Context) (*types.UserLoginResponse, error) {
	claims, expired, maxExpired := middlewares.JwtAuth.MapClaimsAndExpired(ctx)
	if claims == nil {
		return nil, errors.TokenDataError
	}

	if !expired {
		return nil, errors.RefreshActiveTokenError
	}

	if maxExpired {
		return nil, errors.TokenExpiredError
	}

	metadata, err := meta.Parse(claims)
	if err != nil {
		return nil, err
	}

	token, err := middlewares.JwtAuth.CreateToken(ctx, metadata)
	if err != nil {
		return nil, err
	}

	return &types.UserLoginResponse{
		Token: token,
	}, err
}

// CurUserMenuTree 获取当前用户的菜单树
func CurUserMenuTree(ctx *gin.Context) (tree.Tree, error) {
	md, err := meta.Get(ctx)
	if err != nil {
		return nil, err
	}

	// 如果是超级管理员就直接返回全部菜单
	if md.RoleKey == consts.JwtSuperAdmin {
		return AllMenu(ctx)
	}

	// 查询角色所属菜单
	rm := model.RoleMenu{}
	rmList, err := rm.RoleMenus(ctx, md.RoleID)
	if err != nil {
		return nil, err
	}

	// 获取菜单的所有id
	var ids []int64
	for _, item := range rmList {
		ids = append(ids, item.MenuID)
	}

	// 获取指定id的所有菜单
	var menu model.Menu
	menuList, _ := menu.All(ctx, "id in ?", ids)
	var listTree []tree.Tree
	for _, item := range menuList {
		listTree = append(listTree, item)
	}

	return tree.BuildTree(listTree), nil
}
