package admin

import (
	"github.com/casbin/casbin/v2"
	"github.com/ihuanglei/authenticator/models"
	"github.com/ihuanglei/authenticator/pkg/context"
	"github.com/ihuanglei/authenticator/pkg/errors"
	"github.com/ihuanglei/authenticator/pkg/logger"
	"github.com/ihuanglei/authenticator/pkg/mapper/st"
	"github.com/simplexwork/common"
)

// GetResources 资源列表
// @tags 管理 - 权限管理
// @Summary 资源列表
// @Accept x-www-form-urlencoded
// @Success 200 {object} context.JSONResult
// @Param name query string false "名称"
// @Router /admin/authority/resource [get]
// @Security AdminKeyAuth
func GetResources(query st.ResourceQuery, ctx *context.Context) {
	count, resources, err := models.GetResources(query)
	if err != nil {
		ctx.BadRequestByError(err)
		return
	}
	ctx.JSONList(count, "resources", resources)
}

// CreateRole 创建角色
// @tags 管理 - 权限管理
// @Summary 创建角色
// @Accept x-www-form-urlencoded
// @Success 200 {object} context.JSONResult
// @Param name formData string false "名称"
// @Param res_id query []string false "资源id"
// @Router /admin/authority/role/create [post]
// @Security AdminKeyAuth
func CreateRole(e *casbin.Enforcer, ctx *context.Context) {
	name := ctx.QueryTrim("name")
	ids := ctx.QueryStrings("res_id")
	if len(ids) == 0 {
		ctx.BadRequestByError(errors.ErrArgument)
		return
	}
	var resIDs []common.ID
	for _, id := range ids {
		resIDs = append(resIDs, common.StrToID(id))
	}
	ress, err := models.GetResourcesByIDs(resIDs)
	if err != nil {
		ctx.BadRequestByError(err)
		return
	}
	roleID, err := models.CreateRole(name, ress)
	if err != nil {
		ctx.BadRequestByError(err)
		return
	}
	p := [][]string{}
	for _, res := range ress {
		rule := []string{roleID.Str(), res.URL, res.Method}
		p = append(p, rule)
	}
	_, err = e.AddPolicies(p)
	if err != nil {
		ctx.BadRequestByError(err)
		return
	}
	ctx.JSONEmpty()
}

// UpdateRole 更新角色
// @tags 管理 - 权限管理
// @Summary 更新角色
// @Accept x-www-form-urlencoded
// @Success 200 {object} context.JSONResult
// @Param id path string true "角色编号"
// @Param name formData string false "名称"
// @Param res_id query []string false "资源id"
// @Router /admin/authority/role/{id}/update [post]
// @Security AdminKeyAuth
func UpdateRole(e *casbin.Enforcer, ctx *context.Context) {
	roleID := ctx.ParamsID("roleID")
	name := ctx.QueryTrim("name")
	ids := ctx.QueryStrings("res_id")
	if len(ids) == 0 {
		ctx.BadRequestByError(errors.ErrArgument)
		return
	}
	var resIDs []common.ID
	for _, id := range ids {
		resIDs = append(resIDs, common.StrToID(id))
	}
	ress, err := models.GetResourcesByIDs(resIDs)
	if err != nil {
		ctx.BadRequestByError(err)
		return
	}
	err = models.UpdataRole(roleID, name, ress)
	if err != nil {
		ctx.BadRequestByError(err)
		return
	}
	p := [][]string{}
	for _, res := range ress {
		rule := []string{roleID.Str(), res.URL, res.Method}
		p = append(p, rule)
	}
	_, err = e.RemoveFilteredPolicy(0, roleID.Str())
	if err != nil {
		ctx.BadRequestByError(err)
		return
	}
	_, err = e.AddPolicies(p)
	if err != nil {
		ctx.BadRequestByError(err)
		return
	}
	ctx.JSONEmpty()
}

// DeleteRole 删除角色
// @tags 管理 - 权限管理
// @Summary 删除角色
// @Accept x-www-form-urlencoded
// @Success 200 {object} context.JSONResult
// @Param id path string true "角色编号"
// @Router /admin/authority/role/{id}/delete [post]
// @Security AdminKeyAuth
func DeleteRole(e *casbin.Enforcer, ctx *context.Context) {
	roleID := ctx.ParamsID("roleID")
	err := models.DeleteRole(roleID)
	if err != nil {
		ctx.BadRequestByError(err)
		return
	}
	_, err = e.RemoveFilteredPolicy(0, roleID.Str())
	if err != nil {
		logger.Errorln(err)
	}
	ctx.JSONEmpty()
}

// GetRoles 角色列表
// @tags 管理 - 权限管理
// @Summary 角色列表
// @Accept x-www-form-urlencoded
// @Success 200 {object} context.JSONResult
// @Param name query string false "名称"
// @Router /admin/authority/role [get]
// @Security AdminKeyAuth
func GetRoles(query st.RoleQuery, e *casbin.Enforcer, ctx *context.Context) {
	count, roles, err := models.GetRoles(query)
	if err != nil {
		ctx.BadRequestByError(err)
		return
	}
	ctx.JSONList(count, "roles", roles)
}

// GetRoleResources 角色资源
// @tags 管理 - 权限管理
// @Summary 角色资源
// @Accept x-www-form-urlencoded
// @Success 200 {object} context.JSONResult
// @Param id path string true "角色编号"
// @Router /admin/authority/role/{id}/resource [get]
// @Security AdminKeyAuth
func GetRoleResources(ctx *context.Context) {
	roleID := ctx.ParamsID("roleID")
	count, resources, err := models.GetRoleResourceByID(roleID)
	if err != nil {
		ctx.BadRequestByError(err)
		return
	}
	ctx.JSONList(count, "resources", resources)
}
