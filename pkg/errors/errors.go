package errors

// Error .
type Error [2]interface{}

// Code 错误码
func (err Error) Code() int {
	if v, ok := err[0].(int); ok {
		return v
	}
	return -100
}

// Error 错误内容
func (err Error) Error() string {
	if v, ok := err[1].(string); ok {
		return v
	}
	return "unknown"
}

//
var (
	ErrUnknown = Error{10500, "未知错误"}

	ErrNotLogin        = Error{10001, "未登录"}
	ErrAuthExpired     = Error{10002, "登录信息已过期"}
	ErrAuthInvalidData = Error{10003, "无效的登录信息"}

	ErrUserExist           = Error{10100, "用户已存在"}
	ErrUserNotExist        = Error{10101, "用户不存在"}
	ErrUserNameExist       = Error{10102, "用户名已存在"}
	ErrUserEmailExist      = Error{10103, "邮箱已存在"}
	ErrUserMobileExist     = Error{10104, "手机号已存在"}
	ErrUserMobileNotBind   = Error{10105, "未绑定手机号"}
	ErrUserForbidden       = Error{10106, "用户已被禁用"}
	ErrUserLocked          = Error{10107, "用户已被锁定"}
	ErrUserDelete          = Error{10108, "用户已注销"}
	ErrUserUnActivated     = Error{10109, "用户未激活"}
	ErrUserAlreadyActivate = Error{10110, "用户已激活"}
	ErrUserAlreadyBind     = Error{10111, "用户已经绑定"}
	ErrAddressNotFound     = Error{10112, "地址不存在"}
	ErrDictNotFound        = Error{10113, "字典中数据不存在"}

	ErrArgument        = Error{10400, "参数错误"}
	ErrPassword        = Error{10401, "密码长度必须为6-20位"}
	ErrInvalidPassword = Error{10402, "密码错误"}
	ErrSamePassword    = Error{10403, "新密码不能和原密码一致"}
	ErrEmail           = Error{10404, "邮箱格式错误"}
	ErrName            = Error{10405, "用户名长度必须为5-20位，字母、数字、下划线组合，不允许纯数字"}
	ErrMobile          = Error{10406, "手机号格式错误"}
	ErrCode            = Error{10407, "验证码错误或已过期"}
	ErrActiveCode      = Error{10408, "无效的激活码"}
	ErrAuthenticator   = Error{10409, "认证参数错误"}
	ErrThirdCode       = Error{10410, "无效的第三方认证令牌"}
	ErrAvatar          = Error{10411, "头像地址不能为空"}
	ErrNickname        = Error{10412, "昵称长度必须为1-15个字"}

	ErrWeiXinMPCode        = Error{10501, "微信小程序临时登录凭证错误"}
	ErrWeiXinMPKey         = Error{10502, "调用微信小程序登录返回的key不存在或错误"}
	ErrWeiXinEncryptedData = Error{10503, "微信小程序用户信息加密数据错误"}
	ErrWeiXinIV            = Error{10504, "微信小程序加密算法的初始向量错误"}

	ErrRoleNotFound = Error{10600, "角色不存在"}
	ErrRoleExist    = Error{10601, "角色已存在"}
)
