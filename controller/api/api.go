package api

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/gbrlsnchs/jwt/v3"
	"github.com/go-macaron/binding"
	"github.com/ihuanglei/authenticator/pkg/config"
	"github.com/ihuanglei/authenticator/pkg/consts"
	"github.com/ihuanglei/authenticator/pkg/context"
	"github.com/ihuanglei/authenticator/pkg/errors"
	"github.com/ihuanglei/authenticator/pkg/logger"
	"github.com/ihuanglei/authenticator/pkg/mapper/st"
	"github.com/simplexwork/common"
	"gopkg.in/macaron.v1"
)

// Router 设置路由
func Router(m *macaron.Macaron) {
	m.Group("/v1/api", func() {

		m.Group("/reg", func() {
			m.Post("/name", binding.Bind(st.RegisterNameForm{}), RegisterWithNameAndPassword)
			m.Post("/email", binding.Bind(st.RegisterEmailForm{}), RegisterWithEmailAndPassword)
			m.Post("/mobile", binding.Bind(st.RegisterMobileForm{}), RegisterWithMobileAndPassword)
			m.Post("/third", binding.Bind(st.RegisterWithThirdForm{}), RegisterWithThirdCode)
			m.Group("/weixinmp", func() {
				m.Post("/userinfo/:id", binding.Bind(st.WeiXinMPForm{}), RegisterWithWeiXinMP)
				m.Post("/mobile/:id", binding.Bind(st.WeiXinMPForm{}), RegisterWithWeiXinMPPhone)
			})
			m.Get("/activate", binding.Bind(st.ActivateUserForm{}), ActivateUser)
			m.Post("/activate/resend", binding.Bind(st.EmailForm{}), ReSendActivateCode)
		})

		m.Group("/code", func() {
			m.Post("/reg", binding.Bind(st.MobileForm{}), SendCodeWithReg)
			m.Post("/login", binding.Bind(st.MobileForm{}), SendCodeWithLogin)

			m.Group("", func() {
				m.Post("/password", SendCodeWithPassword)
				m.Group("/bind", func() {
					m.Post("/mobile", binding.Bind(st.MobileForm{}), SendCodeWithBindMobile)
					m.Post("/email", binding.Bind(st.EmailForm{}), SendCodeWithBindEmail)
				})
			}, Authorize)

			m.Group("/forgot", func() {
				m.Post("/email", binding.Bind(st.EmailForm{}), SendCodeByForgotPasswordWithEmail)
			})
		})

		m.Group("/forgot", func() {
			m.Post("/reset/email", binding.Bind(st.ResetPasswordWithEmailCodeForm{}), ResetPasswordByCodeWithEmail)
		})

		m.Group("/login", func() {
			m.Post("/", binding.Bind(st.LoginForm{}), Login)
			m.Post("/mobile", binding.Bind(st.LoginWithMobileAndCodeForm{}), LoginByMobile)
			m.Group("/th", func() {
				m.Get("/:id", binding.Bind(st.LoginWithThirdForm{}), RedirectURLForThird)
				m.Post("/:id", binding.Bind(st.LoginWithThirdCodeForm{}), LoginByThirdCode)
				m.Post("/weixinmp/:id", binding.Bind(st.LoginWithWeiXinMPCodeForm{}), LoginByWeiXinMPCode)
			})
		})

		m.Group("/profile", func() {
			m.Get("/", Info)
			m.Group("/update", func() {
				m.Post("/avatar", UpdateAvatar)
				m.Post("/nickname", UpdateNickname)
				m.Post("/gender", UpdateGender)

				m.Group("/password", func() {
					m.Post("/old", binding.Bind(st.UpdatePasswordWithOldPasswordForm{}), UpdatePasswordWithOldPassword)
					m.Post("/mobile", binding.Bind(st.UpdatePasswordWithCodeForm{}), UpdatePasswordWithCode)
				})
				m.Group("/bind", func() {
					m.Post("/mobile", binding.Bind(st.UpdateMobileWithCodeForm{}), UpdateMobileWithCode)
					m.Post("/mobile/weixinmp/:id", binding.Bind(st.WeiXinMPForm{}), UpdateMobileWithWeiXinMP)
					m.Post("/email", binding.Bind(st.UpdateEmailWithCodeForm{}), UpdateEmailWidthCode)
				})
			})

		}, Authorize)
	})
}

// Authorize 登录认证
func Authorize(config *config.Config, ctx *context.Context) {
	authorizations := strings.Split(ctx.Req.Header.Get(consts.HeaderAuthorizationKey), " ")
	if len(authorizations) != 2 || authorizations[0] != "Authenticator" || authorizations[1] == "" {
		ctx.JSONAuth(errors.ErrNotLogin.Error())
		return
	}
	authCode := authorizations[1]
	hs256 := jwt.NewHS256([]byte(config.Server.JWTSecret))
	var p jwt.Payload
	now := time.Now()
	iatValidator := jwt.IssuedAtValidator(now)
	expValidator := jwt.ExpirationTimeValidator(now)
	audValidator := jwt.AudienceValidator(jwt.Audience{"authenticator"})
	verifyOption := jwt.ValidatePayload(&p, iatValidator, expValidator, audValidator)
	_, err := jwt.Verify([]byte(authCode), hs256, &p, verifyOption)
	if err != nil {
		logger.Debug(err)
		ctx.JSONAuth(errors.ErrAuthExpired.Error())
		return
	}
	sessionUser := new(context.SessionUser)
	err = json.Unmarshal([]byte(p.Subject), sessionUser)
	if err != nil {
		logger.Debug(err)
		ctx.JSONAuth(errors.ErrAuthInvalidData.Error())
		return
	}
	sessionUser.UserID = common.StrToID(sessionUser.UserStrID)
	ctx.SessionUser = sessionUser
}
