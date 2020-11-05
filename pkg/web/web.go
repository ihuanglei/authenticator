package web

import (
	"fmt"
	"net/http"
	"os"

	"github.com/ihuanglei/authenticator/controller/admin"
	"github.com/ihuanglei/authenticator/controller/api"
	"github.com/ihuanglei/authenticator/pkg/authzer"
	"github.com/ihuanglei/authenticator/pkg/config"
	"github.com/ihuanglei/authenticator/pkg/consts"
	"github.com/ihuanglei/authenticator/pkg/context"
	"github.com/ihuanglei/authenticator/pkg/logger"

	"github.com/simplexwork/cache"
	"github.com/simplexwork/common"
	"gopkg.in/macaron.v1"
)

// Run start http
func Run(config *config.Config) {

	cacheType := cache.Memory
	if config.Cache == "redis" {
		cacheType = cache.Redis
	}

	cacheOption := cache.Option{
		Type: cacheType,
		Redis: cache.RedisOption{
			Host:     config.Redis.Host,
			Port:     config.Redis.Port,
			Password: config.Redis.Password,
			DB:       config.Redis.DB,
		},
		Memory: cache.MemoryOption{
			Size: config.Memory.Size,
		},
	}

	cache := cache.Cacher(&cacheOption)

	m := macaron.New()

	m.Use(macaron.Recovery())
	m.Use(logger.MacaronLogger())
	m.Use(macaron.Renderer())
	m.Use(context.Contexter())

	// 注入
	m.Map(cache)
	m.Map(authzer.NewAuthzer())
	m.Map(config)

	m.NotFound(func(ctx *context.Context) {
		ctx.NotFound()
	})
	m.InternalServerError(func(ctx *context.Context) {
		ctx.InternalServerError()
	})

	// 解决跨域访问
	m.Options("/*", func(ctx *context.Context) {
		ctx.Resp.Header().Set("Access-Control-Allow-Headers", fmt.Sprintf("%s,%s", consts.HeaderAuthorizationKey, consts.HeaderAuthorizationAdminKey))
		ctx.Resp.Header().Set("Access-Control-Allow-Methods", "POST,GET")
	})

	api.Router(m)
	admin.Router(m)

	// IP PORT
	host := config.Server.Host
	if len(host) == 0 {
		host = os.Getenv("AUTH_HOST")
		if len(host) == 0 {
			host = "0.0.0.0"
		}
	}
	port := config.Server.Port
	if port == 0 {
		port = common.StrToInt(os.Getenv("AUTH_PORT"), 9990)
	}

	addr := fmt.Sprintf("%s:%v", host, port)

	logger.Infof("[WEB] Listening on %s", addr)
	logger.Fatalln(http.ListenAndServe(addr, m))
}
