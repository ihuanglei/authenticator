package models

import (
	"github.com/ihuanglei/authenticator/pkg/consts"
	"github.com/ihuanglei/authenticator/pkg/errors"
	"github.com/simplexwork/common"
	"xorm.io/builder"
)

// 新增角色
func createRole(name string, roleRess []*roleResource) (common.ID, error) {
	roleID, err := _IDWorker.Next()
	if err != nil {
		return 0, err
	}

	var role role
	role.RoleID = roleID
	role.Name = name
	role.CreateTime = common.Now()
	role.UpdateTime = role.CreateTime
	role.Status = consts.Normal

	session := _Engine.NewSession()
	defer session.Close()
	if err := session.Begin(); err != nil {
		return 0, err
	}

	_, err = session.Insert(&role)
	if err != nil {
		return 0, err
	}
	for _, v := range roleRess {
		v.RoleID = roleID
		v.CreateTime = role.CreateTime
	}
	_, err = session.InsertMulti(roleRess)
	if err != nil {
		return 0, err
	}
	return roleID, session.Commit()
}

// 更新角色
func updateRole(roleID common.ID, name string, roleRess []*roleResource) error {
	now := common.Now()
	role := &role{Name: name, UpdateTime: common.Now()}
	for _, v := range roleRess {
		v.RoleID = roleID
		v.CreateTime = now
	}
	cols := []string{"update_time"}
	if !common.IsEmpty(name) {
		cols = append(cols, "name")
	}

	session := _Engine.NewSession()
	defer session.Close()
	if err := session.Begin(); err != nil {
		return err
	}
	_, err := session.Cols(cols...).Where("role_id = ?", roleID).Update(role)
	if err != nil {
		return err
	}
	var roleResource roleResource
	_, err = session.Where("role_id = ?", roleID).Delete(&roleResource)
	if err != nil {
		return err
	}
	_, err = session.InsertMulti(roleRess)
	if err != nil {
		return err
	}
	return session.Commit()
}

// 删除角色
func deleteRole(roleID common.ID) error {
	role := &role{Status: consts.Delete, UpdateTime: common.Now()}
	session := _Engine.NewSession()
	defer session.Close()
	if err := session.Begin(); err != nil {
		return err
	}
	_, err := session.Cols("status", "update_time").Where("role_id = ?", roleID).Update(role)
	if err != nil {
		return err
	}
	return session.Commit()
}

// 获取角色信息
func getRoleByID(roleID common.ID) (*role, error) {
	return getRole(builder.Eq{"role_id": roleID})
}

// 获取角色信息
func getRole(cond builder.Cond) (*role, error) {
	role := new(role)
	has, err := _Engine.Where(cond).Get(role)
	if err != nil {
		return nil, err
	} else if !has {
		return nil, errors.ErrRoleNotFound
	}
	return role, nil
}

// 角色数量
func getRoleCount(cond builder.Cond) (int64, error) {
	role := new(role)
	return _Engine.Where(cond).Count(role)
}

// 获取角色列表
func getRoles(cond builder.Cond, page, limit int) (int64, []*role, error) {
	if limit <= 0 {
		limit = consts.PageSize
	}
	if page <= 0 {
		page = 1
	}
	start := (page - 1) * limit
	var roles = make([]*role, 0)
	count, err := _Engine.Omit("id").Desc("create_time").Where(cond).Limit(limit, start).FindAndCount(&roles)
	if err != nil {
		return 0, nil, err
	}
	return count, roles, nil
}

// 获取角色资源
func getRoleResources(cond builder.Cond) (int64, []*roleResource, error) {
	var roleResources = make([]*roleResource, 0)
	count, err := _Engine.Where(cond).FindAndCount(&roleResources)
	if err != nil {
		return 0, nil, err
	}
	return count, roleResources, nil
}
