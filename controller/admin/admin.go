package admin

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/casbin/casbin/v2"
	"github.com/gbrlsnchs/jwt/v3"
	"github.com/go-macaron/binding"
	"github.com/ihuanglei/authenticator/models"
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

	m.Group("/v1/admin", func() {

		m.Group("/authority", func() {
			m.Group("/role", func() {
				m.Get("/", binding.Bind(st.RoleQuery{}), GetRoles)
				m.Post("/create", CreateRole)
				m.Post("/:roleID/update", UpdateRole)
				m.Post("/:roleID/delete", DeleteRole)
				m.Get("/:roleID/resource", GetRoleResources)
			})
			m.Get("/resource", binding.Bind(st.ResourceQuery{}), GetResources)
		})

		m.Group("/user", func() {
			m.Get("/", binding.Bind(st.UserQuery{}), GetUsers)
			m.Get("/:userID", GetUser)
			m.Post("/:userID/forbidden/:forbidden", Forbidden)
			m.Post("/:userID/password", ChangePassword)
			m.Post("/:userID/reset", ResetLoginError)
			m.Post("/:userID/activate", ActivateUser)
			m.Post("/:userID/role", AddRoleForUser)
			m.Get("/:userID/role", GetRoleForUser)
			m.Get("/:userID/login", binding.Bind(st.EmptyQuery{}), GetUserLogins)
		})

		m.Group("/dict", func() {
			m.Get("/", GetDictByCate)
			m.Get("/one", GetOneDict)
			m.Post("/create", binding.Bind(st.DictForm{}), CreateDict)
			m.Post("/:dictID/update", binding.Bind(st.DictForm{}), UpdateDict)
			m.Post("/one/update", binding.Bind(st.DictForm{}), UpdateOneDict)
			m.Post("/:dictID/del", DelDict)
		})

	}, Authorize)
}

// Authorize 登录认证及权限管理
func Authorize(enforce *casbin.Enforcer, config *config.Config, ctx *context.Context) {
	authorizations := strings.Split(ctx.Req.Header.Get(consts.HeaderAuthorizationAdminKey), " ")
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
	_, err = models.GetUserByID(sessionUser.UserID)
	if err != nil {
		logger.Error(err)
		ctx.JSONAuth(errors.ErrAuthInvalidData.Error())
	}
	ctx.SessionUser = sessionUser

	method := common.ToLower(ctx.Req.Method)
	path := ctx.Req.URL.Path
	ok, err := enforce.Enforce(ctx.UserStrID, path, method)
	if err != nil {
		ctx.Error(err)
		return
	}
	if !ok {
		ctx.AccessDenied()
		return
	}
}
