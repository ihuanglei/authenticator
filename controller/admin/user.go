package admin

import (
	"github.com/casbin/casbin/v2"
	"github.com/ihuanglei/authenticator/models"
	"github.com/ihuanglei/authenticator/pkg/consts"
	"github.com/ihuanglei/authenticator/pkg/context"
	"github.com/ihuanglei/authenticator/pkg/errors"
	"github.com/ihuanglei/authenticator/pkg/mapper/st"
	"github.com/simplexwork/common"
)

// GetUserLogins 管理员获取用户登录历史
// @tags 管理 - 用户管理
// @Summary 管理员获取用户登录历史
// @Accept x-www-form-urlencoded
// @Success 200 {object} context.JSONResult
// @Param id path string true "用户编号"
// @Router /admin/user/{id}/login [get]
// @Security AdminKeyAuth
func GetUserLogins(query st.EmptyQuery, ctx *context.Context) {
	userID := ctx.ParamsID("userID")
	count, userLogins, err := models.GetUserLoginByID(userID, query)
	if err != nil {
		ctx.BadRequestByError(err)
		return
	}
	ctx.JSONList(count, "logins", userLogins)
}

// GetUsers 管理员获取用户列表
// @tags 管理 - 用户管理
// @Summary 管理员获取用户列表
// @Accept x-www-form-urlencoded
// @Success 200 {object} context.JSONResult
// @Param name query string false "用户名"
// @Router /admin/user [get]
// @Security AdminKeyAuth
func GetUsers(query st.UserQuery, ctx *context.Context) {
	count, users, err := models.GetUsers(query)
	if err != nil {
		ctx.BadRequestByError(err)
		return
	}
	ctx.JSONList(count, "users", users)
}

// GetUser 管理员获取用户信息
// @tags 管理 - 用户管理
// @Summary 管理员获取用户信息
// @Accept x-www-form-urlencoded
// @Success 200 {object} context.JSONResult
// @Param id path string true "用户编号"
// @Router /admin/user/{id} [get]
// @Security AdminKeyAuth
func GetUser(ctx *context.Context) {
	userID := ctx.ParamsID("userID")
	if _, err := models.GetUserByID(userID); err != nil {
		ctx.BadRequestByError(err)
		return
	}
	userInfo, err := models.GetUserInfoByID(userID)
	if err != nil {
		ctx.BadRequestByError(err)
		return
	}
	ctx.JSON(userInfo)
}

// ResetLoginError 重置错误登录
// @tags 管理 - 用户管理
// @Summary 重置错误登录
// @Accept x-www-form-urlencoded
// @Success 200 {object} context.JSONResult
// @Param id path string true "用户编号"
// @Router /admin/user/{id}/reset [post]
// @Security AdminKeyAuth
func ResetLoginError(ctx *context.Context) {
	userID := ctx.ParamsID("userID")
	err := models.ResetLoginErrorForUser(userID)
	if err != nil {
		ctx.BadRequestByError(err)
		return
	}
	ctx.JSONEmpty()
}

// Forbidden 禁止/恢复用户
// @tags 管理 - 用户管理
// @Summary 禁止/恢复用户
// @Accept x-www-form-urlencoded
// @Success 200 {object} context.JSONResult
// @Param id path string true "用户编号"
// @Param forbidden path string true "状态" Enums(available,unAvailable)
// @Router /admin/user/{id}/forbidden/{forbidden} [post]
// @Security AdminKeyAuth
func Forbidden(ctx *context.Context) {
	userID := ctx.ParamsID("userID")
	forbidden := consts.NewForbidden(ctx.Params("forbidden"))
	err := models.UpdateForbiddenForUser(userID, forbidden)
	if err != nil {
		ctx.BadRequestByError(err)
		return
	}
	ctx.JSONEmpty()
}

// ChangePassword 修改密码
// @tags 管理 - 用户管理
// @Summary 修改密码
// @Accept x-www-form-urlencoded
// @Success 200 {object} context.JSONResult
// @Param id path string true "用户编号"
// @Param password formData string true "密码"
// @Router /admin/user/{id}/password [post]
// @Security AdminKeyAuth
func ChangePassword(ctx *context.Context) {
	userID := ctx.ParamsID("userID")
	password := ctx.QueryTrim("password")
	if !common.IsSimplePassword(password) {
		ctx.BadRequestByError(errors.ErrPassword)
		return
	}
	err := models.UpdatePassword1ForUser(userID, password)
	if err != nil {
		ctx.BadRequestByError(err)
		return
	}
	ctx.JSONEmpty()
}

// ActivateUser 激活用户
// @tags 管理 - 用户管理
// @Summary 激活用户
// @Accept x-www-form-urlencoded
// @Success 200 {object} context.JSONResult
// @Param id path string true "用户编号"
// @Router /admin/user/{id}/activate [post]
// @Security AdminKeyAuth
func ActivateUser(ctx *context.Context) {
	userID := ctx.ParamsID("userID")
	err := models.ActivateUserWithOutCode(userID)
	if err != nil {
		ctx.BadRequestByError(err)
		return
	}
	ctx.JSONEmpty()
}

// AddRoleForUser 添加角色
// @tags 管理 - 用户管理
// @Summary 添加用户角色
// @Accept x-www-form-urlencoded
// @Success 200 {object} context.JSONResult
// @Param id path string true "用户编号"
// @Param role_id query []string false "角色编号"
// @Router /admin/user/{id}/role [post]
// @Security AdminKeyAuth
func AddRoleForUser(e *casbin.Enforcer, ctx *context.Context) {
	userID := ctx.ParamsID("userID")
	roleIDs := ctx.QueryStrings("role_id")
	_, err := e.DeleteRolesForUser(userID.Str())
	if err != nil {
		ctx.BadRequestByError(err)
		return
	}
	_, err = e.AddRolesForUser(userID.Str(), roleIDs)
	if err != nil {
		ctx.BadRequestByError(err)
		return
	}
	ctx.JSONEmpty()
}

// GetRoleForUser 获取用户角色
// @tags 管理 - 用户管理
// @Summary 获取用户角色
// @Accept x-www-form-urlencoded
// @Success 200 {object} context.JSONResult
// @Param id path string true "用户编号"
// @Router /admin/user/{id}/role [get]
// @Security AdminKeyAuth
func GetRoleForUser(e *casbin.Enforcer, ctx *context.Context) {
	userID := ctx.ParamsID("userID")
	roles, err := e.GetRolesForUser(userID.Str())
	if err != nil {
		ctx.BadRequestByError(err)
		return
	}
	ctx.JSONList(int64(len(roles)), "roles", roles)
}
