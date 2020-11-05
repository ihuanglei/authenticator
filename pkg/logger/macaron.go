package logger

import (
	"net/http"
	"time"

	macaron "gopkg.in/macaron.v1"
)

// MacaronLogger .
func MacaronLogger() macaron.Handler {
	return func(ctx *macaron.Context) {
		start := time.Now()
		Infof("[WEB] Started %s %s for %s", ctx.Req.Method, ctx.Req.RequestURI, ctx.RemoteAddr())
		rw := ctx.Resp.(macaron.ResponseWriter)
		ctx.Next()
		Infof("[WEB] Completed %s %s %v %s in %v", ctx.Req.Method, ctx.Req.RequestURI, rw.Status(), http.StatusText(rw.Status()), time.Since(start))
	}
}
