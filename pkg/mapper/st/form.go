package st

import (
	"github.com/ihuanglei/authenticator/pkg/context"
	"github.com/ihuanglei/authenticator/pkg/errors"
	"github.com/simplexwork/common"

	"github.com/go-macaron/binding"
	"gopkg.in/macaron.v1"
)

type formErrors map[string]error

func (fe formErrors) Val(filed string) error {
	if v, ok := fe[filed]; ok {
		return v
	}
	return errors.ErrArgument
}

var (
	_errors = formErrors{
		"LOGINNAME":           errors.ErrName,
		"CODE":                errors.ErrCode,
		"MOBILE":              errors.ErrMobile,
		"EMAIL":               errors.ErrEmail,
		"PASSWORD":            errors.ErrPassword,
		"OLDPASSWORD":         errors.ErrPassword,
		"ACTIVATECODE":        errors.ErrActiveCode,
		"WEIXINMPCODE":        errors.ErrWeiXinMPCode,
		"WEIXINMPKEY":         errors.ErrWeiXinMPKey,
		"WEIXINENCRYPTEDDATA": errors.ErrWeiXinEncryptedData,
		"WEIXINIV":            errors.ErrWeiXinIV,
	}
)

// FormError .
type FormError struct {
}

// DictForm .
type DictForm struct {
	FormError
	Name  string `form:"name" binding:"Required"`
	Cate  string `form:"cate" binding:"Required"`
	Value string `form:"value" binding:"Required"`
	TP    string `form:"tp"`
}

// RegisterNameForm 用户名注册表单
type RegisterNameForm struct {
	FormError
	LoginName string `form:"login_name" binding:"Required"`
	Password  string `form:"password" binding:"Required;Password"`
}

// RegisterWithThirdForm 三方登录表单
type RegisterWithThirdForm struct {
	FormError
	Code string `form:"code" binding:"Required"`
}

// WeiXinMPForm 微信小程序表单
type WeiXinMPForm struct {
	FormError
	WeiXinMPKey         string `form:"key" binding:"Required"`
	WeiXinEncryptedData string `form:"encrypted_data" binding:"Required"`
	WeiXinIV            string `form:"iv" binding:"Required"`
}

// LoginForm 用户名邮箱手机号和密码登录表单
type LoginForm struct {
	FormError
	LoginName string `form:"login_name" binding:"Required"`
	Password  string `form:"password" binding:"Required;Password"`
}

// LoginWithThirdForm 第三方code登录表单
type LoginWithThirdForm struct {
	FormError
	State string `form:"state"`
}

// LoginWithThirdCodeForm 第三方code登录表单
type LoginWithThirdCodeForm struct {
	LoginWithThirdForm
	Code   string `form:"code" binding:"Required"`
	OpenID string
	Type   string
}

// LoginWithWeiXinMPCodeForm 小程序code登录表单
type LoginWithWeiXinMPCodeForm struct {
	FormError
	WeiXinMPCode string `form:"code" binding:"Required"`
	OpenID       string
	Type         string
	OnlySession  bool `form:"only_session"`
}

// UpdatePasswordWithOldPasswordForm 验证码修改密码表单
type UpdatePasswordWithOldPasswordForm struct {
	FormError
	Password    string `form:"password" binding:"Required;Password"`
	OldPassword string `form:"old_password" binding:"Required;Password"`
}

// UpdatePasswordWithCodeForm 验证码修改密码表单
type UpdatePasswordWithCodeForm struct {
	FormError
	Password string `form:"password" binding:"Required;Password"`
	Code     string `form:"code" binding:"Required;Size(6)"`
}

// ************ 手机号相关表单

// MobileForm 手机号表单
type MobileForm struct {
	FormError
	Mobile string `form:"mobile" binding:"Required;Mobile"`
}

// RegisterMobileForm 手机注册表单
type RegisterMobileForm struct {
	MobileForm
	Password string `form:"password" binding:"Required;Password"`
	Code     string `form:"code" binding:"Required;Size(6)"`
}

// LoginWithMobileAndCodeForm 手机号验证码登录表单
type LoginWithMobileAndCodeForm struct {
	MobileForm
	Code string `form:"code" binding:"Required;Size(6)"`
}

// UpdateMobileWithCodeForm 手机验证码更新手机号表单
type UpdateMobileWithCodeForm struct {
	MobileForm
	Code string `form:"code" binding:"Required;Size(6)"`
}

// ************ 邮件相关表单

// ActivateUserForm 邮件激活表单
type ActivateUserForm struct {
	FormError
	ActivateCode string `form:"s" binding:"Required"`
}

// EmailForm 邮件表单
type EmailForm struct {
	FormError
	Email string `form:"email" binding:"Required;Email"`
}

// RegisterEmailForm 邮箱注册表单
type RegisterEmailForm struct {
	EmailForm
	Password string `form:"password" binding:"Required;Password"`
}

// ResetPasswordWithEmailCodeForm 忘记密码邮件验证码修改密码表单
type ResetPasswordWithEmailCodeForm struct {
	EmailForm
	Password string `form:"password" binding:"Required;Password"`
	Code     string `form:"code" binding:"Required;Size(6)"`
}

// UpdateEmailWithCodeForm 邮件验证码更新邮箱表单
type UpdateEmailWithCodeForm struct {
	EmailForm
	Code string `form:"code" binding:"Required;Size(6)"`
}

func (from FormError) Error(mctx *macaron.Context, errs binding.Errors) {
	if len(errs) > 0 {
		ctx := context.Context{Context: mctx}
		err := errs[0]
		if len(err.Fields()) > 0 {
			err := _errors.Val(common.ToUpper(err.Fields()[0]))
			ctx.BadRequestByError(err)
		} else {
			ctx.BadRequestByError(errors.ErrArgument)
		}
	}
}

func init() {
	binding.AddRule(&binding.Rule{
		IsMatch: func(rule string) bool {
			return rule == "Mobile"
		},
		IsValid: func(errs binding.Errors, name string, v interface{}) (bool, binding.Errors) {
			mobile, ok := v.(string)
			if !ok {
				return false, errs
			}
			if !common.IsMobile(mobile) {
				errs.Add([]string{name}, "Mobile", errors.ErrMobile.Error())
				return false, errs
			}
			return true, errs
		},
	})

	binding.AddRule(&binding.Rule{
		IsMatch: func(rule string) bool {
			return rule == "Password"
		},
		IsValid: func(errs binding.Errors, name string, v interface{}) (bool, binding.Errors) {
			password, ok := v.(string)
			if !ok {
				return false, errs
			}
			if !common.IsSimplePassword(password) {
				errs.Add([]string{name}, "Password", errors.ErrPassword.Error())
				return false, errs
			}
			return true, errs
		},
	})
}
