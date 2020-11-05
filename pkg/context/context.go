package context

import (
	"net/http"
	"runtime"
	"strings"
	"time"

	"github.com/ihuanglei/authenticator/pkg/config"
	"github.com/ihuanglei/authenticator/pkg/errors"
	"github.com/ihuanglei/authenticator/pkg/logger"
	"github.com/simplexwork/common"

	"gopkg.in/macaron.v1"
)

// SessionUser 会话中的用户
type SessionUser struct {
	// UserID 用户ID
	UserID    common.ID `json:"-"`
	UserStrID string    `json:"user_id"`
	Nickname  string    `json:"nickname"`
	Avatar    string    `json:"avatar"`
}

// JSONResult .
type JSONResult struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg,omitempty"`
	Data interface{} `json:"data,omitempty"`
} //@name Result

// Context .
type Context struct {
	*macaron.Context
	*SessionUser
	StartTime time.Time
	IP        string
	Secret    string
	Expire    int64
}

// ParamsID .
func (c *Context) ParamsID(name string) common.ID {
	return common.Int64ToID(c.ParamsInt64(name))
}

// QueryID .
func (c *Context) QueryID(name string) common.ID {
	return common.Int64ToID(c.QueryInt64(name))
}

// NotFound .
func (c *Context) NotFound() {
	c.Context.JSON(http.StatusNotFound, &JSONResult{Code: http.StatusNotFound, Msg: "Not Found"})
}

// InternalServerError .
func (c *Context) InternalServerError() {
	c.Context.JSON(http.StatusInternalServerError, &JSONResult{Code: http.StatusInternalServerError})
}

// Error .
func (c *Context) Error(err error) {
	var buf [1024]byte
	n := runtime.Stack(buf[:], false)
	logger.Error(string(buf[:n]))
	c.Context.JSON(http.StatusOK, &JSONResult{Code: http.StatusInternalServerError, Msg: err.Error()})
}

// BadRequest .
func (c *Context) BadRequest(msg string) {
	c.Context.JSON(http.StatusOK, &JSONResult{Code: http.StatusBadRequest, Msg: msg})
}

// BadRequestByCode .
func (c *Context) BadRequestByCode(code int, msg string) {
	c.Context.JSON(http.StatusOK, &JSONResult{Code: code, Msg: msg})
}

// BadRequestByError .
func (c *Context) BadRequestByError(err error) {
	if e, ok := err.(errors.Error); ok {
		c.BadRequestByCode(e.Code(), e.Error())
	} else if err != nil {
		c.Error(err)
	}
}

// JSON .
func (c *Context) JSON(i interface{}) {
	c.Context.JSON(http.StatusOK, &JSONResult{Code: http.StatusOK, Data: i})
}

// JSONList .
func (c *Context) JSONList(count int64, listName string, list interface{}) {
	c.Context.JSON(http.StatusOK, &JSONResult{
		Code: http.StatusOK,
		Data: map[string]interface{}{
			"count":  count,
			listName: list,
		},
	})
}

// JSONByCode .
func (c *Context) JSONByCode(code int, i interface{}) {
	c.Context.JSON(http.StatusOK, &JSONResult{Code: code, Data: i})
}

// JSONEmpty .
func (c *Context) JSONEmpty() {
	c.Context.JSON(http.StatusOK, &JSONResult{Code: http.StatusOK})
}

// JSONAuth 未认证
func (c *Context) JSONAuth(msg string) {
	c.Context.JSON(http.StatusOK, &JSONResult{Code: http.StatusUnauthorized, Data: msg})
}

// AccessDenied 无权限访问
func (c *Context) AccessDenied() {
	c.Context.JSON(http.StatusOK, &JSONResult{Code: http.StatusForbidden, Msg: "access denied"})
}

// Redirect 跳转
func (c *Context) Redirect(url string) {
	c.Context.Redirect(url, http.StatusFound)
}

// Contexter init
func Contexter() macaron.Handler {
	return func(config *config.Config, ctx *macaron.Context) {
		c := &Context{
			Context:   ctx,
			StartTime: time.Now(),
			IP:        ip(ctx.Req.Request),
			Secret:    config.Server.JWTSecret,
			Expire:    config.Server.Expire,
		}
		ctx.Resp.Header().Set("Access-Control-Allow-Origin", "*")
		if common.IsEmpty(ctx.Req.Header.Get("Content-Type")) {
			ctx.Req.Header.Set("Content-Type", "form-urlencoded")
		}
		ctx.Map(c)
	}
}

func ip(r *http.Request) string {
	ip := r.Header.Get("X-Forwarded-For")
	if strings.Contains(ip, "127.0.0.1") || ip == "" {
		ip = r.Header.Get("X-real-ip")
	}
	if ip == "" {
		ip = r.RemoteAddr
		if strings.Contains(ip, ":") {
			ip = strings.Split(ip, ":")[0]
		}
	}
	if ip == "" {
		ip = "127.0.0.1"
	}
	return ip
}
