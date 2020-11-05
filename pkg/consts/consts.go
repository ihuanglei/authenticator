package consts

// HeaderAuthorizationKey header普通鉴权key
const HeaderAuthorizationKey = "X-ACMS-Authorization"

// HeaderAuthorizationAdminKey header管理后台鉴权key
const HeaderAuthorizationAdminKey = "X-AACMS-Authorization"

// LoginErrorCount 登录错误次数
const LoginErrorCount = 5

// PageSize 每页默认长度
const PageSize = 20

// Mode 注册方式
type Mode int

const (
	// Name 用户名
	Name Mode = 1
	// Email 邮箱
	Email Mode = 2
	// Mobile 手机号
	Mobile Mode = 3
	// Third 第三方
	Third Mode = 4
)

// Status 数据逻辑状态
type Status int

const (
	// Normal 正常
	Normal Status = 1
	// Delete 删除
	Delete Status = -1
)

// Str 返回值
func (st Status) Str() string {
	switch st {
	case Normal:
		return "normal"
	case Delete:
		return "delete"
	}
	return "normal"
}

// Val 实际值
func (st Status) Val() int {
	return int(st)
}

// MarshalText json格式返回
func (st Status) MarshalText() ([]byte, error) {
	return []byte(st.Str()), nil
}

// Activate 激活
type Activate int

const (
	// Activated 激活
	Activated Activate = 1
	// UnActivated 未激活
	UnActivated Activate = -1
)

// Str 返回值
func (a Activate) Str() string {
	switch a {
	case Activated:
		return "activated"
	}
	return "unActivated"
}

// MarshalText json格式返回
func (a Activate) MarshalText() ([]byte, error) {
	return []byte(a.Str()), nil
}

// NewActivate 创建
func NewActivate(val string) Activate {
	if val == "activated" {
		return Activated
	}
	return UnActivated
}

// Forbidden 禁用状态
type Forbidden int

const (
	// Available 可用
	Available Forbidden = 1
	// UnAvailable 不可用
	UnAvailable Forbidden = -1
)

// Check 检查参数
// func (f Forbidden) Check() error {
// 	if f != Available && f != UnAvailable {
// 		return ErrArguments{Val: "参数错误"}
// 	}
// 	return nil
// }

// Str 返回值
func (f Forbidden) Str() string {
	switch f {
	case Available:
		return "available"
	case UnAvailable:
		return "unAvailable"
	}
	return "unknown"
}

// MarshalText json格式返回
func (f Forbidden) MarshalText() ([]byte, error) {
	return []byte(f.Str()), nil
}

// NewForbidden 创建
func NewForbidden(val string) Forbidden {
	if val == "available" {
		return Available
	}
	return UnAvailable
}

// Gender 性别
type Gender int

const (
	// Unknown 未知
	Unknown Gender = iota
	// Male 男
	Male
	// Female 女
	Female
)

// Str 返回值
func (g Gender) Str() string {
	switch g {
	case Male:
		return "male"
	case Female:
		return "female"
	}
	return "unknown"
}

// MarshalText json格式返回
func (g Gender) MarshalText() ([]byte, error) {
	return []byte(g.Str()), nil
}

// NewGender 性别
func NewGender(val string) Gender {
	if val == "male" {
		return Male
	} else if val == "female" {
		return Female
	} else {
		return Unknown
	}
}

// Query 查询
type Query struct {
	Page  int `form:"page"`
	Limit int `form:"limit"`
}
