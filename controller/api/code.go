package api

import (
	"fmt"
	"time"

	"github.com/ihuanglei/authenticator/models"
	"github.com/ihuanglei/authenticator/pkg/context"
	"github.com/ihuanglei/authenticator/pkg/errors"
	"github.com/ihuanglei/authenticator/pkg/mapper/st"
	"github.com/simplexwork/cache"
	"github.com/simplexwork/common"
)

const (
	// 登录用户修改密码验证码
	codeKeyWithChangePassword = "__code_change_password_%v"
	// 登录用户绑定手机号
	codeKeyWithBindMobile = "__code_bind_mobile_%v"
	// 登录用户绑定邮箱
	codeKeyWithBindEmail = "__code_bind_email_%v"

	// 忘记密码验证过你吗
	codeKeyByForgotPwdWithEmail = "__code_forgot_email_%v"

	// 手机号登录验证码
	codeKeyWithLogin = "__code_login_%v"
	// 手机号注册验证码
	codeKeyWithReg = "__code_reg_%v"
)

func saveCodeAndSendWithMobile(t, mobile string, ex int, cache cache.Cache) error {
	key := fmt.Sprintf(t, mobile)
	exp := time.Minute * time.Duration(ex)
	code := common.RandomNumber(6)
	err := cache.Set(key, code, exp)
	if err != nil {
		return err
	}
	go sendMessageWithMobile(mobile, code)
	return nil
}

func saveCodeAndSendWithEmail(t, email string, ex int, cache cache.Cache) error {
	key := fmt.Sprintf(t, email)
	exp := time.Minute * time.Duration(ex)
	code := common.RandomNumber(6)
	err := cache.Set(key, code, exp)
	if err != nil {
		return err
	}
	switch t {
	case codeKeyWithBindEmail:
		go sendBindMessageWithEmail(email, code)
	case codeKeyByForgotPwdWithEmail:
		go sendForgotMessageWithEmail(email, code)
	}
	return nil
}

// SendCodeWithReg 注册验证码
// @tags 前端 - 手机验证码
// @Summary 手机注册验证码
// @Accept x-www-form-urlencoded
// @Success 200 {object} context.JSONResult "验证码"
// @Param mobile formData string false "手机号"
// @Router /api/code/reg [post]
func SendCodeWithReg(form st.MobileForm, ctx *context.Context, cache cache.Cache) {
	has, err := models.HasUserByMobile(form.Mobile)
	if err != nil {
		ctx.Error(err)
		return
	}
	if has {
		ctx.BadRequestByError(errors.ErrUserExist)
		return
	}
	if err := saveCodeAndSendWithMobile(codeKeyWithReg, form.Mobile, 5, cache); err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSONEmpty()
}

// SendCodeWithLogin 登录验证码
// @tags 前端 - 手机验证码
// @Summary 手机登录验证码
// @Accept x-www-form-urlencoded
// @Success 200 {object} context.JSONResult "验证码"
// @Param mobile formData string false "手机号"
// @Router /api/code/login [post]
func SendCodeWithLogin(form st.MobileForm, ctx *context.Context, cache cache.Cache) {
	if err := models.CheckUserByMobile(form.Mobile); err != nil {
		ctx.BadRequestByError(err)
		return
	}
	if err := saveCodeAndSendWithMobile(codeKeyWithLogin, form.Mobile, 5, cache); err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSONEmpty()
}

// SendCodeWithPassword 更新密码验证码(用户已登录)
// @tags 前端 - 手机验证码
// @Summary 更新密码验证码(用户已登录)
// @Accept x-www-form-urlencoded
// @Success 200 {object} context.JSONResult "验证码"
// @Router /api/code/password [post]
// @Security ApiKeyAuth
func SendCodeWithPassword(ctx *context.Context, cache cache.Cache) {
	user, err := models.GetUserByID(ctx.UserID)
	if err != nil {
		ctx.BadRequestByError(err)
		return
	}
	if user.Mobile == user.UserID.Str() {
		ctx.BadRequestByError(errors.ErrUserMobileNotBind)
		return
	}
	if err := saveCodeAndSendWithMobile(codeKeyWithChangePassword, user.Mobile, 5, cache); err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSONEmpty()
}

// SendCodeWithBindMobile 绑定或更新手机号验证码(用户已登录)
// @tags 前端 - 手机验证码
// @Summary 绑定或更新手机号验证码(用户已登录)
// @Accept x-www-form-urlencoded
// @Success 200 {object} context.JSONResult "验证码"
// @Param mobile formData string false "手机号"
// @Router /api/code/bind/mobile [post]
// @Security ApiKeyAuth
func SendCodeWithBindMobile(form st.MobileForm, ctx *context.Context, cache cache.Cache) {
	has, err := models.HasUserByMobile(form.Mobile)
	if err != nil {
		ctx.Error(err)
		return
	}
	if has {
		ctx.BadRequestByError(errors.ErrUserMobileExist)
		return
	}
	if err := saveCodeAndSendWithMobile(codeKeyWithBindMobile, form.Mobile, 5, cache); err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSONEmpty()
}

// SendCodeWithBindEmail 绑定或更新邮箱验证码(用户已登录)
// @tags 前端 - 邮件验证码
// @Summary 绑定或更新邮箱验证码(用户已登录)
// @Accept x-www-form-urlencoded
// @Success 200 {object} context.JSONResult "验证码"
// @Param email formData string false "邮箱"
// @Router /api/code/bind/email [post]
// @Security ApiKeyAuth
func SendCodeWithBindEmail(form st.EmailForm, ctx *context.Context, cache cache.Cache) {
	has, err := models.HasUserByEmail(form.Email)
	if err != nil {
		ctx.Error(err)
		return
	}
	if has {
		ctx.BadRequestByError(errors.ErrUserEmailExist)
		return
	}
	if err := saveCodeAndSendWithEmail(codeKeyWithBindEmail, form.Email, 5, cache); err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSONEmpty()
}

// SendCodeByForgotPasswordWithEmail 重置密码验证码
// @tags 前端 - 邮件验证码
// @Summary 忘记密码重置验证码
// @Accept x-www-form-urlencoded
// @Success 200 {object} context.JSONResult "验证码"
// @Param email formData string false "邮箱"
// @Router /api/code/forgot/email [post]
// @Security ApiKeyAuth
func SendCodeByForgotPasswordWithEmail(form st.EmailForm, ctx *context.Context, cache cache.Cache) {
	has, err := models.HasUserByEmail(form.Email)
	if err != nil {
		ctx.Error(err)
		return
	}
	if !has {
		ctx.BadRequestByError(errors.ErrUserNotExist)
		return
	}
	if err := saveCodeAndSendWithEmail(codeKeyByForgotPwdWithEmail, form.Email, 5, cache); err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSONEmpty()
}
