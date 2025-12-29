package igrpc

import (
	"context"

	"google.golang.org/grpc"
	grpcmd "google.golang.org/grpc/metadata"
)

// unaryServerInterceptor is a gRPC unary server interceptor
func (s *Server) unaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		var cancel func()
		md, _ := grpcmd.FromIncomingContext(ctx)
		replyHeader := grpcmd.MD{}
		tr := &Transport{
			operation:   info.FullMethod,
			reqHeader:   md,
			replyHeader: replyHeader,
		}
		ctx = NewServerContext(ctx, tr)
		if s.timeout > 0 {
			ctx, cancel = context.WithTimeout(ctx, s.timeout)
			defer cancel()
		}

		h := func(ctx context.Context, req any) (any, error) {
			return handler(ctx, req)
		}
		if next := s.middleware.Match(tr.Operation()); len(next) > 0 {
			h = Chain(next...)(h)
		}
		reply, err := h(ctx, req)
		if len(replyHeader) > 0 {
			_ = grpc.SetHeader(ctx, replyHeader)
		}
		return reply, err
	}
}

// wrappedStream is rewrite grpc stream's context
type wrappedStream struct {
	grpc.ServerStream
	ctx context.Context
}

func NewWrappedStream(ctx context.Context, stream grpc.ServerStream) grpc.ServerStream {
	return &wrappedStream{
		ServerStream: stream,
		ctx:          ctx,
	}
}

func (w *wrappedStream) Context() context.Context {
	return w.ctx
}

// streamServerInterceptor is a gRPC stream server interceptor
func (s *Server) streamServerInterceptor() grpc.StreamServerInterceptor {
	return func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		ctx := ss.Context()
		md, _ := grpcmd.FromIncomingContext(ctx)
		replyHeader := grpcmd.MD{}
		ctx = NewServerContext(ctx, &Transport{
			operation:   info.FullMethod,
			reqHeader:   md,
			replyHeader: replyHeader,
		})
		ws := NewWrappedStream(ctx, ss)
		err := handler(srv, ws)
		if len(replyHeader) > 0 {
			_ = grpc.SetHeader(ctx, replyHeader)
		}
		return err
	}
}
