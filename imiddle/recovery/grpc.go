package recovery

import (
	"context"
	"fmt"
	"runtime"

	"github.com/Verthandii/spring/ierrors"
	"github.com/Verthandii/spring/ilogger"
	"github.com/Verthandii/spring/transport/igrpc"
)

type GRPCFn func(ctx context.Context, req any, r any) (reply any, err error)

type GRPCOption func(m *GRPCRecovery)

func WithGRPCFn(fn GRPCFn) GRPCOption {
	return func(m *GRPCRecovery) {
		m.fn = fn
	}
}

type GRPCRecovery struct {
	fn GRPCFn
}

func NewGRPC(opts ...GRPCOption) igrpc.Middleware {
	m := &GRPCRecovery{
		fn: defaultGRPCFn,
	}

	for _, opt := range opts {
		opt(m)
	}

	return m.handler()
}

func (c *GRPCRecovery) handler() igrpc.Middleware {
	return func(next igrpc.Handler) igrpc.Handler {
		return func(ctx context.Context, req any) (reply any, err error) {
			defer func() {
				if r := recover(); r != nil {
					reply, err = c.fn(ctx, req, r)
				}
			}()
			return next(ctx, req)
		}
	}
}

func defaultGRPCFn(ctx context.Context, req any, rerr any) (reply any, err error) {
	buf := make([]byte, 64<<10)
	n := runtime.Stack(buf, false)
	buf = buf[:n]
	ilogger.ErrorwCtx(ctx, "recovery middleware", "error message", fmt.Errorf("%v: %+v\n%s\n", rerr, req, buf))
	err = ierrors.InternalServerError
	return
}
