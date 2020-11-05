package st

import (
	"github.com/ihuanglei/authenticator/pkg/consts"
	"github.com/simplexwork/common"
)

// UserLoginDto 登录信息
type UserLoginDto struct {
	IP         string          `json:"ip"`
	Country    string          `json:"country"`
	Province   string          `json:"province"`
	Region     string          `json:"region"`
	City       string          `json:"city"`
	Lat        string          `json:"lat"`
	Lng        string          `json:"lng"`
	CreateTime common.DateTime `json:"create_time"`
}

// UserDto 用户
type UserDto struct {
	UserID     common.ID        `json:"user_id"`
	Name       string           `json:"name"`
	Email      string           `json:"email"`
	Mobile     string           `json:"mobile"`
	Status     consts.Status    `json:"status"`
	Error      int              `json:"error"`
	Forbidden  consts.Forbidden `json:"forbidden"`
	Activate   consts.Activate  `json:"activate"`
	CreateTime common.DateTime  `json:"create_time"`
}

// UserInfoDto 用户详情
type UserInfoDto struct {
	UserID   common.ID     `json:"user_id"`
	Nickname string        `json:"nickname"`
	Avatar   string        `json:"avatar"`
	Gender   consts.Gender `json:"gender"`
	QQ       string        `json:"qq"`
	WeiXin   string        `json:"weixin"`
	Province string        `json:"province"`
	City     string        `json:"city"`
	County   string        `json:"county"`
}

// RegisterDto 注册信息
type RegisterDto struct {
	Email     string
	Mobile    string `json:"mobile"`
	LoginName string
	Password  string
	Nickname  string `json:"nickname"`
	IP        string
	// 三方注册的数据
	TP       string `json:"tp"`
	OpenID   string `json:"open_id"`
	Avatar   string `json:"avatar"`
	Gender   string `json:"gender"`
	Province string `json:"province"`
	City     string `json:"city"`
}

// AddressDto 地址
type AddressDto struct {
	AddressID  common.ID       `json:"address_id"`
	Name       string          `json:"name"`
	Mobile     string          `json:"mobile"`
	Province   string          `json:"province"`
	City       string          `json:"city"`
	County     string          `json:"county"`
	Address    string          `json:"address"`
	Zip        string          `json:"zip"`
	Def        string          `json:"def"`
	CreateTime common.DateTime `json:"create_time"`
	UpdateTime common.DateTime `json:"update_time"`
}

// LoginDto 登录
type LoginDto struct {
	LoginName string
	Password  string
	Mobile    string
	IP        string
	OpenID    string
	Type      string
}

// DictDto 字典
type DictDto struct {
	// 字典编号
	DictID common.ID `json:"dict_id"`
	// 类型
	Cate string `json:"cate"`
	// 名称
	Name string `json:"name"`
	// 内容
	Value string `json:"value"`
	// 业务类型
	TP string `json:"tp"`
}

// ResourceDto 资源
type ResourceDto struct {
	ID string `json:"id"`
	// 名称
	Name string `json:"name"`
	// 资源
	URL string `json:"url"`
	// 方法
	Method string `json:"method"`
}

// RoleDto 角色
type RoleDto struct {
	// 编号
	RoleID common.ID `json:"role_id"`
	// 名称
	Name string `json:"name"`
}

// RoleResourceDto 角色资源
type RoleResourceDto struct {
	// 名称
	Name string `json:"name"`
	// 资源
	URL string `json:"url"`
	// 方法
	Method string `json:"method"`
}
