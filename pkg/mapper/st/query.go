package st

import (
	"github.com/ihuanglei/authenticator/pkg/consts"
)

// UserQuery 用户搜索
type UserQuery struct {
	consts.Query
	Keyword    string   `form:"keyword"`
	Name       string   `form:"name"`
	Email      string   `form:"email"`
	Mobile     string   `form:"mobile"`
	Genders    []string `form:"gender"`
	Forbiddens []string `form:"forbidden"`
	Activates  []string `form:"activate"`
}

// ResourceQuery 资源搜索
type ResourceQuery struct {
	consts.Query
	Name string `form:"name"`
}

// RoleQuery 角色搜索
type RoleQuery struct {
	consts.Query
	Name string `form:"name"`
}

// EmptyQuery 搜索
type EmptyQuery struct {
	consts.Query
}
