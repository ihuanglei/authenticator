package api

import (
	"fmt"

	"github.com/ihuanglei/authenticator/models"
	"github.com/ihuanglei/authenticator/pkg/consts"
	"github.com/ihuanglei/authenticator/pkg/context"
	"github.com/ihuanglei/authenticator/pkg/errors"
	"github.com/ihuanglei/authenticator/pkg/mapper/st"
	"github.com/ihuanglei/authenticator/pkg/third"
	"github.com/ihuanglei/authenticator/pkg/third/weixin"
	"github.com/simplexwork/cache"
	"github.com/simplexwork/common"
)

// UpdateAddress 更新地址
// func UpdateAddress(ctx *context.Context) {
// 	province := ctx.QueryTrim("province")
// 	city := ctx.QueryTrim("city")
// 	county := ctx.QueryTrim("county")
// 	address := ctx.QueryTrim("address")
// 	err := models.UpdateAddress(ctx.UserID, province, city, county, address)
// 	if e, ok := err.(models.Err); ok {
// 		ctx.BadRequestByCode(e.Code(), e.Error())
// 		return
// 	} else if err != nil {
// 		ctx.Error(err)
// 		return
// 	}
// 	ctx.JSONEmpty()
// }

// Info 用户信息
// @tags 前端 - 用户信息
// @Summary 用户信息
// @Accept x-www-form-urlencoded
// @Success 200 {object} context.JSONResult
// @Router /api/profile [get]
// @Security ApiKeyAuth
func Info(ctx *context.Context) {
	user, err := models.GetUserByID(ctx.UserID)
	if err != nil {
		ctx.BadRequestByError(err)
		return
	}
	userInfo, err := models.GetUserInfoByID(ctx.UserID)
	if err != nil {
		ctx.BadRequestByError(err)
		return
	}
	ret := map[string]interface{}{}
	if user.Mobile == user.UserID.Str() {
		ret["mobile"] = ""
	} else {
		ret["mobile"] = common.HideMobile(user.Mobile)
	}
	if user.Email == user.UserID.Str() {
		ret["email"] = ""
	} else {
		ret["email"] = user.Email
	}
	if user.Name == user.UserID.Str() {
		ret["name"] = ""
	} else {
		ret["name"] = user.Name
	}
	ret["nickname"] = userInfo.Nickname
	ret["avatar"] = userInfo.Avatar
	ret["gender"] = userInfo.Gender.Str()
	ret["province"] = userInfo.Province
	ret["city"] = userInfo.City
	ret["county"] = userInfo.County
	ret["weixin"] = userInfo.WeiXin
	ret["qq"] = userInfo.QQ
	ctx.JSON(ret)
}

// UpdateAvatar 修改头像
// @tags 前端 - 用户信息
// @Summary 修改头像
// @Accept x-www-form-urlencoded
// @Success 200 {object} context.JSONResult
// @Param avatar formData string false "头像地址"
// @Router /api/profile/update/avatar [post]
// @Security ApiKeyAuth
func UpdateAvatar(ctx *context.Context) {
	avatar := ctx.QueryTrim("avatar")
	if err := models.UpdateAvatarForUser(ctx.UserID, avatar); err != nil {
		ctx.BadRequestByError(err)
		return
	}
	ctx.JSONEmpty()
}

// UpdateNickname 修改昵称
// @tags 前端 - 用户信息
// @Summary 修改昵称
// @Accept x-www-form-urlencoded
// @Success 200 {object} context.JSONResult
// @Param nickname formData string false "昵称"
// @Router /api/profile/update/nickname [post]
// @Security ApiKeyAuth
func UpdateNickname(ctx *context.Context) {
	nickname := ctx.QueryTrim("nickname")
	if err := models.UpdateNicknameForUser(ctx.UserID, nickname); err != nil {
		ctx.BadRequestByError(err)
		return
	}
	ctx.JSONEmpty()
}

// UpdateGender 修改性别
// @tags 前端 - 用户信息
// @Summary 修改性别
// @Accept x-www-form-urlencoded
// @Success 200 {object} context.JSONResult
// @Param gender formData string false "性别" Enums(male, female)
// @Router /api/profile/update/gender [post]
// @Security ApiKeyAuth
func UpdateGender(ctx *context.Context) {
	gender := consts.NewGender(ctx.QueryTrim("gender"))
	if err := models.UpdateGenderForUser(ctx.UserID, gender); err != nil {
		ctx.BadRequestByError(err)
		return
	}
	ctx.JSONEmpty()
}

// UpdatePasswordWithOldPassword 根据原密码修改密码
// @tags 前端 - 用户信息
// @Summary 根据原密码修改密码
// @Accept x-www-form-urlencoded
// @Success 200 {object} context.JSONResult
// @Param old_password formData string false "原密码"
// @Param password formData string false "新密码"
// @Router /api/profile/update/password/old [post]
// @Security ApiKeyAuth
func UpdatePasswordWithOldPassword(form st.UpdatePasswordWithOldPasswordForm, ctx *context.Context) {
	if form.Password == form.OldPassword {
		ctx.BadRequestByError(errors.ErrSamePassword)
		return
	}
	err := models.UpdatePasswordForUser(ctx.UserID, form.OldPassword, form.Password)
	if err != nil {
		ctx.BadRequestByError(err)
		return
	}
	ctx.JSONEmpty()
}

// UpdatePasswordWithCode 根据手机验证码修改密码
// @tags 前端 - 用户信息
// @Summary 根据手机验证码修改密码
// @Accept x-www-form-urlencoded
// @Success 200 {object} context.JSONResult
// @Param password formData string false "新密码"
// @Param code formData string false "验证码"
// @Router /api/profile/update/password/mobile [post]
// @Security ApiKeyAuth
func UpdatePasswordWithCode(form st.UpdatePasswordWithCodeForm, ctx *context.Context, cache cache.Cache) {
	user, err := models.GetUserByID(ctx.UserID)
	if err != nil {
		ctx.BadRequestByError(err)
		return
	}
	key := fmt.Sprintf(codeKeyWithChangePassword, user.Mobile)
	tmpCode, err := cache.GetString(key)
	if err != nil || tmpCode != form.Code {
		ctx.BadRequestByError(errors.ErrCode)
		return
	}
	err = models.UpdatePassword1ForUser(ctx.UserID, form.Password)
	if err != nil {
		ctx.BadRequestByError(err)
		return
	}
	cache.Del(key)
	ctx.JSONEmpty()
}

// UpdateMobileWithCode 根据手机验证码绑定或更新手机号
// @tags 前端 - 用户信息
// @Summary 根据手机验证码绑定或更新手机号
// @Accept x-www-form-urlencoded
// @Success 200 {object} context.JSONResult
// @Param mobile formData string false "手机号"
// @Param code formData string false "验证码"
// @Router /api/profile/update/bind/mobile [post]
// @Security ApiKeyAuth
func UpdateMobileWithCode(form st.UpdateMobileWithCodeForm, ctx *context.Context, cache cache.Cache) {
	key := fmt.Sprintf(codeKeyWithBindMobile, form.Mobile)
	tmpCode, err := cache.GetString(key)
	if err != nil || tmpCode != form.Code {
		ctx.BadRequestByError(errors.ErrCode)
		return
	}
	has, err := models.HasUserByMobile(form.Mobile)
	if err != nil {
		ctx.Error(err)
		return
	}
	if has {
		ctx.BadRequestByError(errors.ErrUserMobileExist)
		return
	}
	err = models.UpdateMobileForUser(ctx.UserID, form.Mobile)
	if err != nil {
		ctx.BadRequestByError(err)
		return
	}
	cache.Del(key)
	ctx.JSONEmpty()
}

// UpdateMobileWithWeiXinMP 微信小程序更新手机号
// @tags 前端 - 用户信息
// @Summary 微信小程序更新手机号
// @Accept x-www-form-urlencoded
// @Success 200 {object} context.JSONResult
// @Param id path string true "第三方编号"
// @Param key formData string false "调用微信小程序登录接口(/login/th/weixinmp/{id})返回的内容"
// @Param encrypted_data formData string false "微信小程序通过getPhoneNumber获取的encryptedData"
// @Param iv formData string false "微信小程序通过getPhoneNumber获取的encryptedData"
// @Router /api/profile/update/bind/mobile/weixinmp/{id} [post]
// @Security ApiKeyAuth
func UpdateMobileWithWeiXinMP(form st.WeiXinMPForm, ctx *context.Context, cache cache.Cache) {
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
	err = models.UpdateMobileForUser(ctx.UserID, mobile)
	if err != nil {
		ctx.BadRequestByError(err)
		return
	}
	ctx.JSONEmpty()
}

// UpdateEmailWidthCode 根据邮箱验证码绑定或更新邮箱
// @tags 前端 - 用户信息
// @Summary 根据邮箱验证码绑定或更新邮箱
// @Accept x-www-form-urlencoded
// @Success 200 {object} context.JSONResult
// @Param email formData string false "邮箱"
// @Param code formData string false "验证码"
// @Router /api/profile/update/bind/email [post]
// @Security ApiKeyAuth
func UpdateEmailWidthCode(form st.UpdateEmailWithCodeForm, ctx *context.Context, cache cache.Cache) {
	key := fmt.Sprintf(codeKeyWithBindEmail, form.Email)
	tmpCode, err := cache.GetString(key)
	if err != nil || tmpCode != form.Code {
		ctx.BadRequestByError(errors.ErrCode)
		return
	}
	has, err := models.HasUserByEmail(form.Email)
	if err != nil {
		ctx.Error(err)
		return
	}
	if has {
		ctx.BadRequestByError(errors.ErrUserEmailExist)
		return
	}
	err = models.UpdateEmailForUser(ctx.UserID, form.Email)
	if err != nil {
		ctx.BadRequestByError(err)
		return
	}
	cache.Del(key)
	ctx.JSONEmpty()
}
