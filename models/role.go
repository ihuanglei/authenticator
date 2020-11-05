package models

import (
	"github.com/ihuanglei/authenticator/pkg/consts"
	"github.com/ihuanglei/authenticator/pkg/convert"
	"github.com/ihuanglei/authenticator/pkg/errors"
	"github.com/ihuanglei/authenticator/pkg/mapper/st"
	"github.com/simplexwork/common"
	"xorm.io/builder"
)

// CreateRole 创建角色
func CreateRole(name string, resource []*st.ResourceDto) (common.ID, error) {
	has, err := HasRoleByName(name)
	if err != nil {
		return 0, err
	}
	if has {
		return 0, errors.ErrRoleExist
	}
	var roleRess = make([]*roleResource, len(resource))
	err = convert.Map(&resource, &roleRess)
	if err != nil {
		return 0, err
	}
	return createRole(name, roleRess)
}

// UpdataRole 更新角色
func UpdataRole(roleID common.ID, name string, resource []*st.ResourceDto) error {
	role, err := getRoleByID(roleID)
	if err != nil {
		return err
	}
	if role == nil {
		return errors.ErrRoleNotFound
	}
	if role.Name != name {
		has, err := HasRoleByName(name)
		if err != nil {
			return err
		}
		if has {
			return errors.ErrRoleExist
		}
	} else {
		name = ""
	}
	var roleRess = make([]*roleResource, len(resource))
	err = convert.Map(&resource, &roleRess)
	if err != nil {
		return err
	}
	return updateRole(roleID, name, roleRess)
}

// GetRoles .
func GetRoles(query st.RoleQuery) (int64, []*st.RoleDto, error) {
	page := query.Page
	limit := query.Limit
	cond := builder.And(builder.Eq{"status": consts.Normal})
	if common.Trim(query.Name) != "" {
		cond = cond.And(builder.Like{"name", query.Name + "%"})
	}
	count, roles, err := getRoles(cond, page, limit)
	if err != nil {
		return 0, nil, err
	}
	var roleDtos = make([]*st.RoleDto, len(roles))
	err = convert.Map(&roles, &roleDtos)
	if err != nil {
		return 0, nil, err
	}
	return count, roleDtos, nil
}

// GetRoleResourceByID .
func GetRoleResourceByID(roleID common.ID) (int64, []*st.RoleResourceDto, error) {
	count, roleResource, err := getRoleResources(builder.Eq{"role_id": roleID})
	if err != nil {
		return 0, nil, err
	}
	var roleResourceDtos = make([]*st.RoleResourceDto, len(roleResource))
	err = convert.Map(&roleResource, &roleResourceDtos)
	if err != nil {
		return 0, nil, err
	}
	return count, roleResourceDtos, nil
}

// HasRoleByName .
func HasRoleByName(name string) (bool, error) {
	count, err := getRoleCount(builder.Eq{"name": name, "status": consts.Normal})
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// HasRoleByID .
func HasRoleByID(roleID common.ID) (bool, error) {
	count, err := getRoleCount(builder.Eq{"role_id": roleID})
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// DeleteRole .
func DeleteRole(roleID common.ID) error {
	has, err := HasRoleByID(roleID)
	if err != nil {
		return err
	}
	if !has {
		return errors.ErrRoleNotFound
	}
	return deleteRole(roleID)
}
