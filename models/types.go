package models

import (
	"time"

	"github.com/ihuanglei/authenticator/pkg/consts"
	"github.com/simplexwork/common"
)

//  用户
type user struct {
	// 递增主键
	ID int64 `xorm:"id PK AUTOINCR COMMENT('主键')"`
	// 用户编号 业务主键
	UserID common.ID `xorm:"BIGINT NOT NULL UNIQUE 'user_id' COMMENT('用户编号')"`
	// 用户名
	Name string `xorm:"VARCHAR(32) NOT NULL UNIQUE 'name' COMMENT('用户名')"`
	// 邮箱
	Email string `xorm:"VARCHAR(32) NOT NULL UNIQUE 'email' COMMENT('邮箱')"`
	// 手机
	Mobile string `xorm:"varchar(32) NOT NULL UNIQUE 'mobile' COMMENT('手机')"`
	// 密码
	Password string `xorm:"VARCHAR(64) NOT NULL 'password' COMMENT('密码')"`
	// 密码密钥
	Salt string `xorm:"VARCHAR(15) NOT NULL 'salt' COMMENT('密码密钥')"`
	// 注册方式
	Mode consts.Mode `xorm:"TINYINT NOT NULL 'mode' COMMENT('注册方式')"`
	// 状态
	Status consts.Status `xorm:"TINYINT NOT NULL DEFAULT 1 'status' COMMENT('状态')"`
	// 错误次数
	Error int `xorm:"TINYINT NOT NULL DEFAULT 0 'error' COMMENT('错误次数')"`
	// 最后一次登录错误时间
	LastErrorTime common.DateTime `xorm:"NOT NULL 'last_error_time' COMMENT('最后一次登录错误时间')"`
	// 禁用
	Forbidden consts.Forbidden `xorm:"TINYINT NOT NULL INDEX 'forbidden' COMMENT('禁用')"`
	// 禁用时间
	ForbiddenTime common.DateTime `xorm:"NOT NULL 'forbidden_time' COMMENT('禁用时间')"`
	// 激活
	Activate consts.Activate `xorm:"TINYINT NOT NULL DEFAULT -1 INDEX 'activate' COMMENT('激活')"`
	// 激活码
	ActivateCode string `xorm:"VARCHAR(32) NOT NULL 'activate_code' COMMENT('激活码')"`
	// 激活时间
	ActivateTime common.DateTime `xorm:"NOT NULL 'activate_time' COMMENT('激活时间')"`
	// 创建时间
	CreateTime common.DateTime `xorm:"NOT NULL 'create_time' COMMENT('创建时间')"`
	// 修改时间
	UpdateTime common.DateTime `xorm:"NOT NULL 'update_time' COMMENT('修改时间')"`
}

func (u *user) canLogin(numb int) bool {
	if numb < 3 {
		numb = 3
	}
	if u.Error >= numb {
		if time.Now().Sub(time.Time(u.LastErrorTime)).Minutes() < 2 {
			return false
		}
		return false
	}
	return true
}

func (u *user) IsActivated() bool {
	return u.Activate == consts.Activated
}

func (u *user) IsForbidden() bool {
	return u.Forbidden == consts.UnAvailable
}

func (u *user) IsDelete() bool {
	return u.Status == consts.Delete
}

// 用户详细信息
type userInfo struct {
	// 递增主键
	ID int64 `xorm:"id PK AUTOINCR COMMENT('主键')"`
	// 用户编号 业务主键
	UserID common.ID `xorm:"BIGINT NOT NULL UNIQUE 'user_id' COMMENT('用户编号')"`
	// 昵称
	Nickname string `xorm:"VARCHAR(60) NOT NULL INDEX 'nickname' COMMENT('昵称')"`
	// 头像
	Avatar string `xorm:"VARCHAR(255) NOT NULL 'avatar' COMMENT('头像')"`
	// 性别
	Gender consts.Gender `xorm:"TINYINT NOT NULL INDEX 'gender' COMMENT('性别')"`
	// qq
	QQ string `xorm:"VARCHAR(20) NOT NULL 'qq' COMMENT('qq')"`
	// weixin
	WeiXin string `xorm:"VARCHAR(40) NOT NULL 'weixin' COMMENT('微信')"`
	// 省
	Province string `xorm:"VARCHAR(30) NOT NULL 'province' COMMENT('省')"`
	// 市
	City string `xorm:"VARCHAR(30) NOT NULL 'city' COMMENT('市')"`
	// 区
	County string `xorm:"VARCHAR(30) NOT NULL 'county' COMMENT('区')"`
	// 注册ip
	IP string `xorm:"VARCHAR(30) NOT NULL 'ip' COMMENT('注册ip')"`
}

//  第三方登录
type userThird struct {
	ID int64 `xorm:"id PK AUTOINCR COMMENT('主键')"`
	// 用户编号
	UserID common.ID `xorm:"BIGINT NOT NULL INDEX 'user_id' COMMENT('用户编号')"`
	// 第三方类型 qq , weibo, weixin
	Type string `xorm:"VARCHAR(15) NOT NULL INDEX(type+openid) 'type' COMMENT('第三方类型')"`
	// 第三方唯一编号
	OpenID string `xorm:"VARCHAR(64) NOT NULL INDEX(type+openid) 'open_id' COMMENT('第三方唯一编号')"`
	// 状态
	Status consts.Status `xorm:"TINYINT NOT NULL DEFAULT 1 'status' COMMENT('状态')"`
	// 创建时间
	CreateTime common.DateTime `xorm:"NOT NULL 'create_time' COMMENT('创建时间')"`
	// 修改时间
	UpdateTime common.DateTime `xorm:"NOT NULL 'update_time' COMMENT('修改时间')"`
}

// 登录历史
type userLogin struct {
	ID int64 `xorm:"id PK AUTOINCR COMMENT('主键')"`
	// 用户编号
	UserID common.ID `xorm:"BIGINT NOT NULL INDEX 'user_id' COMMENT('用户编号')"`
	// 登录ip
	IP string `xorm:"VARCHAR(30) NOT NULL 'ip' COMMENT('登录ip')"`
	// 国家
	Country string `xorm:"VARCHAR(30) NOT NULL 'country' DEFAULT '-' COMMENT('国家')"`
	// 省
	Province string `xorm:"VARCHAR(30) NOT NULL 'province' DEFAULT '-' COMMENT('省')"`
	// 城市
	City string `xorm:"VARCHAR(30) NOT NULL 'city' DEFAULT '-' COMMENT('城市')"`
	// 区
	Region string `xorm:"VARCHAR(30) NOT NULL 'region' DEFAULT '-' COMMENT('区')"`
	// 维度
	Lat string `xorm:"VARCHAR(20) NOT NULL 'lat' DEFAULT '-' COMMENT('维度')"`
	// 经度
	Lng string `xorm:"VARCHAR(20) NOT NULL 'lng' DEFAULT '-' COMMENT('经度')"`
	//
	Geohash string `xorm:"VARCHAR(20) NOT NULL 'geohash' DEFAULT '-' COMMENT('经度')"`
	// 创建时间
	CreateTime common.DateTime `xorm:"NOT NULL 'create_time' COMMENT('创建时间')"`
}

type userAddress struct {
	ID int64 `xorm:"id PK AUTOINCR COMMENT('主键')"`
	// 地址编号
	AddressID common.ID `xorm:"BIGINT NOT NULL UNIQUE 'address_id' COMMENT('地址编号')"`
	// 用户编号
	UserID common.ID `xorm:"BIGINT NOT NULL INDEX 'user_id' COMMENT('用户编号')"`
	// 姓名
	Name string `xorm:"varchar(30) NOT NULL 'name' COMMENT('姓名')"`
	// 手机
	Mobile string `xorm:"varchar(18) NOT NULL 'mobile' COMMENT('手机')"`
	// 省
	Province string `xorm:"VARCHAR(30) NOT NULL 'province' COMMENT('省')"`
	// 市
	City string `xorm:"VARCHAR(30) NOT NULL 'city' COMMENT('市')"`
	// 区
	County string `xorm:"VARCHAR(30) NOT NULL 'county' COMMENT('区')"`
	// 地址
	Address string `xorm:"VARCHAR(255) NOT NULL 'address' COMMENT('地址')"`
	// 邮编
	Zip string `xorm:"VARCHAR(15) NOT NULL 'zip' COMMENT('邮编')"`
	// 默认
	Def string `xorm:"TINYINT NOT NULL DEFAULT 1 'def' COMMENT('默认')"`
	// 状态
	Status consts.Status `xorm:"TINYINT NOT NULL DEFAULT 1 'status' COMMENT('状态')"`
	// 创建时间
	CreateTime common.DateTime `xorm:"NOT NULL 'create_time' COMMENT('创建时间')"`
	// 修改时间
	UpdateTime common.DateTime `xorm:"NOT NULL 'update_time' COMMENT('修改时间')"`
}

//  字典表
type dict struct {
	ID int64 `xorm:"id PK AUTOINCR COMMENT('主键')"`
	// 字典编号
	DictID common.ID `xorm:"BIGINT NOT NULL UNIQUE 'dict_id' COMMENT('字典编号')"`
	// 类型
	Cate string `xorm:"VARCHAR(10) NOT NULL INDEX 'cate' COMMENT('类型')"`
	// 业务类型
	TP string `xorm:"VARCHAR(30) NOT NULL INDEX 'tp' COMMENT('业务类型')"`
	// 名称
	Name string `xorm:"VARCHAR(30) NOT NULL INDEX 'name' COMMENT('名称')"`
	// 内容
	Value string `xorm:"TEXT NOT NULL 'value' COMMENT('内容')"`
	// 创建时间
	CreateTime common.DateTime `xorm:"NOT NULL 'create_time' COMMENT('创建时间')"`
	// 修改时间
	UpdateTime common.DateTime `xorm:"NOT NULL 'update_time' COMMENT('修改时间')"`
	// 状态
	Status consts.Status `xorm:"TINYINT NOT NULL DEFAULT 1 'status' COMMENT('状态')"`
}

type resource struct {
	ID int64 `xorm:"id PK AUTOINCR COMMENT('主键')"`
	// 名称
	Name string `xorm:"VARCHAR(30) NOT NULL INDEX 'name' COMMENT('名称')"`
	// 资源
	URL string `xorm:"VARCHAR(90) NOT NULL 'url' COMMENT('资源')"`
	// 方法
	Method string `xorm:"VARCHAR(10) NOT NULL 'method' COMMENT('方法')"`
	// 创建时间
	CreateTime common.DateTime `xorm:"NOT NULL DEFAULT current_timestamp() 'create_time' COMMENT('创建时间')"`
}

type role struct {
	ID int64 `xorm:"id PK AUTOINCR COMMENT('主键')"`
	// 角色编号
	RoleID common.ID `xorm:"BIGINT NOT NULL UNIQUE 'role_id' COMMENT('角色编号')"`
	// 名称
	Name string `xorm:"VARCHAR(30) NOT NULL INDEX 'name' COMMENT('名称')"`
	// 状态
	Status consts.Status `xorm:"TINYINT NOT NULL DEFAULT 1 'status' COMMENT('状态')"`
	// 创建时间
	CreateTime common.DateTime `xorm:"NOT NULL 'create_time' COMMENT('创建时间')"`
	// 修改时间
	UpdateTime common.DateTime `xorm:"NOT NULL 'update_time' COMMENT('修改时间')"`
}

type roleResource struct {
	ID int64 `xorm:"id PK AUTOINCR COMMENT('主键')"`
	// 角色编号
	RoleID common.ID `xorm:"BIGINT NOT NULL INDEX 'role_id' COMMENT('角色编号')"`
	// 名称
	Name string `xorm:"VARCHAR(30) NOT NULL INDEX 'name' COMMENT('名称')"`
	// 资源
	URL string `xorm:"VARCHAR(90) NOT NULL 'url' COMMENT('资源')"`
	// 方法
	Method string `xorm:"VARCHAR(10) NOT NULL 'method' COMMENT('方法')"`
	// 创建时间
	CreateTime common.DateTime `xorm:"NOT NULL 'create_time' COMMENT('创建时间')"`
}
