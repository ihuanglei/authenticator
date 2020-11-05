package api

import (
	"fmt"

	"github.com/ihuanglei/authenticator/models"
	"github.com/ihuanglei/authenticator/pkg/context"
	"github.com/ihuanglei/authenticator/pkg/errors"
	"github.com/ihuanglei/authenticator/pkg/mapper/st"
	"github.com/simplexwork/cache"
)

// ResetPasswordByCodeWithEmail 忘记密码邮件验证码修改密码
// @tags 前端 - 重置密码
// @Summary 忘记密码重置
// @Accept x-www-form-urlencoded
// @Success 200 {object} context.JSONResult
// @Param email formData string false "邮箱"
// @Param password formData string false "新密码"
// @Param code formData string false "邮件验证码"
// @Router /api/forgot/reset/email [post]
func ResetPasswordByCodeWithEmail(form st.ResetPasswordWithEmailCodeForm, ctx *context.Context, cache cache.Cache) {
	key := fmt.Sprintf(codeKeyByForgotPwdWithEmail, form.Email)
	tmpCode, err := cache.GetString(key)
	if err != nil || tmpCode != form.Code {
		ctx.BadRequestByError(errors.ErrCode)
		return
	}
	user, err := models.GetUserByEmail(form.Email)
	if err != nil {
		ctx.BadRequestByError(err)
		return
	}
	err = models.UpdatePassword1ForUser(user.UserID, form.Password)
	if err != nil {
		ctx.BadRequestByError(err)
		return
	}
	cache.Del(key)
	ctx.JSONEmpty()
}
