package api

import (
	"fmt"
	"strings"

	"github.com/ihuanglei/authenticator/models"
	"github.com/ihuanglei/authenticator/pkg/context"
	"github.com/ihuanglei/authenticator/pkg/convert"
	"github.com/ihuanglei/authenticator/pkg/errors"
	"github.com/ihuanglei/authenticator/pkg/mapper/st"
	"github.com/ihuanglei/authenticator/pkg/third"
	"github.com/ihuanglei/authenticator/pkg/third/weixin"
	"github.com/simplexwork/cache"
	"github.com/simplexwork/common"
)

// RegisterWithThirdCode 第三方QQ，微信，微博使用code注册
// @tags 前端 - 用户注册
// @Summary 第三方QQ，微信，微博使用code注册
// @Accept x-www-form-urlencoded
// @Success 200 {object} context.JSONResult{data=string} "令牌"
// @Param code formData string false "第三方授权后返回的code"
// @Router /api/reg/third [post]
func RegisterWithThirdCode(form st.RegisterWithThirdForm, cache cache.Cache, ctx *context.Context) {
	key := fmt.Sprintf(tokenThirdKey, form.Code)
	data, err := cache.Get(key)
	if err != nil {
		ctx.BadRequestByError(errors.ErrThirdCode)
		return
	}
	var registerDto st.RegisterDto
	if err := common.FromJSON([]byte(data), &registerDto); err != nil {
		ctx.Error(err)
		return
	}
	registerDto.IP = ctx.IP
	userID, err := models.CreateUserWithThird(&registerDto)
	if err != nil {
		ctx.BadRequestByError(err)
		return
	}
	token, err := getUserAndCreateJWTToken(userID, ctx.Secret, ctx.Expire)
	if err != nil {
		ctx.BadRequestByError(err)
		return
	}
	cache.Del(key)
	ctx.JSON(fmt.Sprintf("Authenticator %v", string(token)))
}

// RegisterWithWeiXinMP 微信小程序注册(通过用户信息)
// @tags 前端 - 用户注册
// @Summary 微信小程序注册(通过用户信息)
// @Accept x-www-form-urlencoded
// @Success 200 {object} context.JSONResult{data=string} "令牌"
// @Param id path string true "第三方编号"
// @Param key formData string false "调用微信小程序登录接口(/login/th/weixinmp/{id})返回的内容"
// @Param encrypted_data formData string false "微信小程序通过wx.getUserInfo获取的encryptedData"
// @Param iv formData string false "微信小程序通过wx.getUserInfo获取的iv"
// @Router /api/reg/weixinmp/userinfo/{id} [post]
func RegisterWithWeiXinMP(form st.WeiXinMPForm, cache cache.Cache, ctx *context.Context) {
	key := fmt.Sprintf(sessionKeyWithWeiXinMP, form.WeiXinMPKey)
	data, err := cache.Get(key)
	if err != nil {
		ctx.BadRequestByError(errors.ErrWeiXinMPKey)
		return
	}
	var thirdUserTmp third.User
	if err := common.FromJSON([]byte(data), &thirdUserTmp); err != nil {
		ctx.Error(err)
		return
	}
	id := ctx.Params("id")
	mp, err := getThird(id)
	if err != nil {
		ctx.BadRequestByError(err)
		return
	}
	v, ok := mp.(*weixin.MinPro)
	if !ok {
		ctx.BadRequestByError(errors.ErrDictNotFound)
		return
	}
	thirdUser, err := v.GetMPUser(thirdUserTmp.Ext, form.WeiXinIV, form.WeiXinEncryptedData)
	if err != nil {
		ctx.BadRequestByError(err)
		return
	}
	var registerDto st.RegisterDto
	registerDto.Avatar = thirdUser.Avatar
	registerDto.Nickname = thirdUser.Nickname
	registerDto.Province = thirdUser.Province
	registerDto.City = thirdUser.City
	registerDto.Gender = thirdUser.Gender
	registerDto.TP = thirdUser.TP
	registerDto.OpenID = thirdUser.OpenID
	registerDto.IP = ctx.IP
	userID, err := models.CreateUserWithThird(&registerDto)
	if err != nil {
		ctx.BadRequestByError(err)
		return
	}
	token, err := getUserAndCreateJWTToken(userID, ctx.Secret, ctx.Expire)
	if err != nil {
		ctx.BadRequestByError(err)
		return
	}
	// cache.Del(key)
	ctx.JSON(fmt.Sprintf("Authenticator %v", string(token)))
}

// RegisterWithWeiXinMPPhone 微信小程序注册(通过手机号)
// @tags 前端 - 用户注册
// @Summary 微信小程序注册(通过手机号)
// @Accept x-www-form-urlencoded
// @Success 200 {object} context.JSONResult{data=string} "令牌"
// @Param id path string true "第三方编号"
// @Param key formData string false "调用微信小程序登录接口(/login/th/weixinmp/{id})返回的内容"
// @Param encrypted_data formData string false "微信小程序通过getPhoneNumber获取的encryptedData"
// @Param iv formData string false "微信小程序通过getPhoneNumber获取的encryptedData"
// @Router /api/reg/weixinmp/mobile/{id} [post]
func RegisterWithWeiXinMPPhone(form st.WeiXinMPForm, cache cache.Cache, ctx *context.Context) {
	key := fmt.Sprintf(sessionKeyWithWeiXinMP, form.WeiXinMPKey)
	data, err := cache.Get(key)
	if err != nil {
		ctx.BadRequestByError(errors.ErrWeiXinMPKey)
		return
	}
	var thirdUser third.User
	if err := common.FromJSON([]byte(data), &thirdUser); err != nil {
		ctx.Error(err)
		return
	}
	id := ctx.Params("id")
	mp, err := getThird(id)
	if err != nil {
		ctx.BadRequestByError(err)
		return
	}
	v, ok := mp.(*weixin.MinPro)
	if !ok {
		ctx.BadRequestByError(errors.ErrDictNotFound)
		return
	}
	mobile, err := v.GetMPUserPhone(thirdUser.Ext, form.WeiXinIV, form.WeiXinEncryptedData)
	if err != nil {
		ctx.BadRequestByError(err)
		return
	}
	var registerDto st.RegisterDto
	registerDto.TP = thirdUser.TP
	registerDto.OpenID = thirdUser.OpenID
	registerDto.Mobile = mobile
	registerDto.IP = ctx.IP
	userID, err := models.CreateUserWithThird(&registerDto)
	if err != nil {
		ctx.BadRequestByError(err)
		return
	}
	token, err := getUserAndCreateJWTToken(userID, ctx.Secret, ctx.Expire)
	if err != nil {
		ctx.BadRequestByError(err)
		return
	}
	// cache.Del(key)
	ctx.JSON(fmt.Sprintf("Authenticator %v", string(token)))
}

// RegisterWithNameAndPassword 用户名和密码注册
// @tags 前端 - 用户注册
// @Summary 用户名和密码注册
// @Accept x-www-form-urlencoded
// @Success 200 {object} context.JSONResult{data=string} "令牌"
// @Param login_name formData string false "用户名"
// @Param password formData string false "密码"
// @Router /api/reg/name [post]
func RegisterWithNameAndPassword(form st.RegisterNameForm, ctx *context.Context) {
	registerDto, err := registerForm2Dto(&form)
	if err != nil {
		ctx.Error(err)
		return
	}
	registerDto.IP = ctx.IP
	userID, err := models.CreateUserWithName(registerDto)
	if err != nil {
		ctx.BadRequestByError(err)
		return
	}
	token, err := getUserAndCreateJWTToken(userID, ctx.Secret, ctx.Expire)
	if err != nil {
		ctx.BadRequestByError(err)
		return
	}
	ctx.JSON(fmt.Sprintf("Authenticator %v", string(token)))
}

// RegisterWithEmailAndPassword 邮箱和密码注册
// @tags 前端 - 用户注册
// @Summary 邮箱和密码注册
// @Accept x-www-form-urlencoded
// @Success 200 {object} context.JSONResult{data=string} "令牌"
// @Param email formData string false "邮箱"
// @Param password formData string false "密码"
// @Router /api/reg/email [post]
func RegisterWithEmailAndPassword(form st.RegisterEmailForm, ctx *context.Context) {
	registerDto, err := registerForm2Dto(&form)
	if err != nil {
		ctx.Error(err)
		return
	}
	registerDto.IP = ctx.IP
	userID, activateCode, err := models.CreateUserWithEmail(registerDto)
	if err != nil {
		ctx.BadRequestByError(err)
		return
	}
	go sendActivateMessageWithEmail(userID, activateCode, registerDto.Email)
	ctx.JSONEmpty()
}

// RegisterWithMobileAndPassword 手机号和密码注册
// @tags 前端 - 用户注册
// @Summary 手机号和密码注册
// @Accept x-www-form-urlencoded
// @Success 200 {object} context.JSONResult{data=string} "令牌"
// @Param mobile formData string false "手机号"
// @Param password formData string false "密码"
// @Param code formData string false "验证码"
// @Router /api/reg/mobile [post]
func RegisterWithMobileAndPassword(form st.RegisterMobileForm, cache cache.Cache, ctx *context.Context) {
	key := fmt.Sprintf(codeKeyWithReg, form.Mobile)
	tmpCode, err := cache.GetString(key)
	if err != nil || tmpCode != form.Code {
		ctx.BadRequestByError(errors.ErrCode)
		return
	}
	registerDto, err := registerForm2Dto(&form)
	if err != nil {
		ctx.Error(err)
		return
	}
	registerDto.IP = ctx.IP
	userID, err := models.CreateUserWithMobile(registerDto)
	if err != nil {
		ctx.BadRequestByError(err)
		return
	}
	token, err := getUserAndCreateJWTToken(userID, ctx.Secret, ctx.Expire)
	if err != nil {
		ctx.BadRequestByError(err)
		return
	}
	cache.Del(key)
	ctx.JSON(fmt.Sprintf("Authenticator %v", string(token)))
}

// ActivateUser 邮箱注册激活用户
// @tags 前端 - 用户注册
// @Summary 邮箱注册激活用户
// @Accept x-www-form-urlencoded
// @Success 200 {object} context.JSONResult
// @Param s query string false "邮件激活码"
// @Router /api/reg/activate [get]
func ActivateUser(form st.ActivateUserForm, ctx *context.Context) {
	activateCodeWrap, err := common.Base64Decode(form.ActivateCode)
	if err != nil {
		ctx.BadRequestByError(errors.ErrActiveCode)
		return
	}
	vals := strings.Split(string(activateCodeWrap), "@")
	if len(vals) != 2 || vals[0] == "" || vals[1] == "" {
		ctx.BadRequestByError(errors.ErrActiveCode)
		return
	}
	err = models.ActivateUser(common.StrToID(vals[0]), vals[1])
	if err != nil {
		ctx.BadRequestByError(err)
		return
	}
	ctx.JSONEmpty()
}

// ReSendActivateCode 重发邮件激活码
// @tags 前端 - 用户注册
// @Summary 重发邮件激活码
// @Accept x-www-form-urlencoded
// @Success 200 {object} context.JSONResult
// @Param email formData string false "邮箱"
// @Router /api/reg/activate/resend [post]
func ReSendActivateCode(form st.EmailForm, ctx *context.Context) {
	user, err := models.GetUserByEmail(form.Email)
	if err != nil {
		ctx.BadRequestByError(err)
		return
	}
	activateCode, err := models.UpdateActivateCodeForUser(user.UserID)
	if err != nil {
		ctx.BadRequestByError(err)
		return
	}
	go sendActivateMessageWithEmail(user.UserID, activateCode, form.Email)
	ctx.JSONEmpty()
}

func registerForm2Dto(form interface{}) (*st.RegisterDto, error) {
	registerDto := new(st.RegisterDto)
	if err := convert.Map(form, registerDto); err != nil {
		return nil, err
	}
	return registerDto, nil
}
