package cors

import (
	"context"

	"github.com/Verthandii/spring/transport/igrpc"
)

func (m *CORS) GRPC() igrpc.Middleware {
	return func(next igrpc.Handler) igrpc.Handler {
		return func(ctx context.Context, req any) (reply any, err error) {
			tr, ok := igrpc.FromServerContext(ctx)
			if !ok {
				return nil, nil
			}
			for _, origin := range m.origin {
				tr.ReplyHeader().Append("Access-Control-Allow-Origin", origin)
			}
			tr.ReplyHeader().Set("Access-Control-Allow-Credentials", m.credentials)
			tr.ReplyHeader().Set("Access-Control-Allow-Headers", m.headers)
			tr.ReplyHeader().Set("Access-Control-Allow-Methods", m.methods)
			return next(ctx, req)
		}
	}
}
