package recovery

import (
	"fmt"
	"runtime"

	"github.com/Verthandii/spring/ilogger"
	"github.com/Verthandii/spring/transport/ihttp"
)

type HTTPFn func(ctx *ihttp.Context, r any)

type HTTPOption func(m *HTTPRecovery)

func WithHTTPFn(fn HTTPFn) HTTPOption {
	return func(r *HTTPRecovery) {
		r.fn = fn
	}
}

type HTTPRecovery struct {
	fn HTTPFn
}

func NewHTTP(opts ...HTTPOption) ihttp.Handler {
	m := &HTTPRecovery{}
	m.fn = m.defaultHTTPFn
	for _, opt := range opts {
		opt(m)
	}
	return m.handler()
}

func (m *HTTPRecovery) handler() ihttp.Handler {
	return func(ctx *ihttp.Context) {
		defer func() {
			if r := recover(); r != nil {
				m.fn(ctx, r)
			}
		}()
		ctx.Next()
	}
}

func (m *HTTPRecovery) defaultHTTPFn(c *ihttp.Context, rerr any) {
	buf := make([]byte, 64<<10)
	n := runtime.Stack(buf, false)
	buf = buf[:n]
	ilogger.ErrorwCtx(c.Request.Context(), "服务器内部错误", "panic message", fmt.Errorf("%v: %s", rerr, buf))
}
