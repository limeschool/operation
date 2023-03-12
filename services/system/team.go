package system

import (
	"github.com/jinzhu/copier"
	"github.com/limeschool/gin"
	"operation/errors"
	model "operation/models/system"
	"operation/tools"
	"operation/tools/tree"
	types "operation/types/system"
)

func AllTeam(ctx *gin.Context) (tree.Tree, error) {
	team := model.Team{}
	return team.Tree(ctx)
}

func AddTeam(ctx *gin.Context, in *types.AddTeamRequest) error {
	team := model.Team{}

	// 获取用户所管理的部门
	user := model.User{}
	teamIds, err := user.CurUserTeamIds(ctx)
	if err != nil {
		return err
	}

	if !tools.InList(teamIds, in.ParentID) {
		return errors.NotAddTeamError
	}

	if copier.Copy(&team, in) != nil {
		return errors.AssignError
	}

	return team.Create(ctx)
}

func UpdateTeam(ctx *gin.Context, in *types.UpdateTeamRequest) error {
	team := model.Team{}
	if in.ParentID != 0 && in.ID == in.ParentID {
		return errors.TeamParentIdError
	}

	// 获取用户所管理的部门
	user := model.User{}
	teamIds, err := user.CurUserTeamIds(ctx)
	if err != nil {
		return err
	}

	if !tools.InList(teamIds, in.ID) {
		return errors.NotEditTeamError
	}
	if copier.Copy(&team, in) != nil {
		return errors.AssignError
	}
	return team.Update(ctx)
}

func DeleteTeam(ctx *gin.Context, in *types.DeleteTeamRequest) error {
	team := model.Team{}

	// 获取用户所管理的部门
	user := model.User{}
	teamIds, err := user.CurUserTeamIds(ctx)
	if err != nil {
		return err
	}

	if !tools.InList(teamIds, in.ID) {
		return errors.NotDelTeamError
	}
	return team.DeleteByID(ctx, in.ID)
}
