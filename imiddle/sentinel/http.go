package sentinel

import (
	"fmt"
	"net/http"

	"github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/base"

	"github.com/Verthandii/spring/ierrors"
	"github.com/Verthandii/spring/ilogger"
	"github.com/Verthandii/spring/transport/ihttp"
)

func HTTP() ihttp.Handler {
	return func(c *ihttp.Context) {
		ctx := c.Context
		resource := fmt.Sprintf("%s:%s", ctx.Request.Method, ctx.FullPath())
		ilogger.DebugwCtx(ctx.Request.Context(), "限流组件", ilogger.Any("resource", resource))
		entry, blockErr := api.Entry(
			resource,
			api.WithResourceType(base.ResTypeWeb),
			api.WithTrafficType(base.Inbound),
		)
		if blockErr != nil {
			ilogger.WarnwCtx(ctx.Request.Context(), "触发限流", ilogger.Any("resource", resource))
			abortWithHttpCodeErr(c, http.StatusTooManyRequests)
			return
		}
		defer entry.Exit()

		ctx.Next()
	}
}

func abortWithHttpCodeErr(c *ihttp.Context, code int) {
	c.FailureWithCode(code, ierrors.New(code, "too many requests"), nil)
}
