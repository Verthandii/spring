package igrpc

import (
	"context"

	"google.golang.org/grpc/metadata"
)

// Transport is a gRPC transport.
type Transport struct {
	operation   string
	reqHeader   metadata.MD
	replyHeader metadata.MD
}

type (
	serverTransportKey struct{}
)

// NewServerContext returns a new Context that carries value.
func NewServerContext(ctx context.Context, tr *Transport) context.Context {
	return context.WithValue(ctx, serverTransportKey{}, tr)
}

// FromServerContext returns the Transport value stored in ctx, if any.
func FromServerContext(ctx context.Context) (tr *Transport, ok bool) {
	tr, ok = ctx.Value(serverTransportKey{}).(*Transport)
	return
}

// Operation returns the transport operation.
func (tr *Transport) Operation() string {
	return tr.operation
}

// RequestHeader returns the request header.
func (tr *Transport) RequestHeader() metadata.MD {
	return tr.reqHeader
}

// ReplyHeader returns the reply header.
func (tr *Transport) ReplyHeader() metadata.MD {
	return tr.replyHeader
}
