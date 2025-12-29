package log

import (
	"context"
	"net/http"
	"time"

	"google.golang.org/grpc/metadata"

	"github.com/Verthandii/spring/ierrors"
	"github.com/Verthandii/spring/ilogger"
	"github.com/Verthandii/spring/transport/igrpc"
)

func GRPC() igrpc.Middleware {
	return func(next igrpc.Handler) igrpc.Handler {
		return func(ctx context.Context, req any) (reply any, err error) {
			var (
				code          = http.StatusOK
				reason        string
				operation     string
				requestHeader metadata.MD
				replyHeader   metadata.MD
			)
			startTime := time.Now()
			tr, ok := igrpc.FromServerContext(ctx)
			if ok {
				operation = tr.Operation()
				requestHeader = tr.RequestHeader()
				replyHeader = tr.ReplyHeader()
			}
			reply, err = next(ctx, req)
			if err != nil {
				code, reason = ierrors.Info(err)
			}
			ilogger.DebugwCtx(ctx, "log middleware",
				"operation", operation,
				"request", req,
				"request_header", requestHeader,
				"reply", reply,
				"reply_header", replyHeader,
				"code", code,
				"reason", reason,
				"latency", time.Since(startTime).String(),
			)
			return
		}
	}
}
