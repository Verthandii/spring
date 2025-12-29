package ilogger

import (
	"context"
	"testing"

	"go.uber.org/zap"
)

type _ctxKey string

func TestLogger(t *testing.T) {
	ck := (_ctxKey)("trace_id")
	logger := NewZapLogger(WithValuer("trace", func(ctx context.Context) interface{} {
		return ctx.Value(ck)
	}))
	SetLogger(logger)

	// 测试是否能正确传递上下文
	// 测试是否正确打印调用函数
	ctx := context.WithValue(context.Background(), ck, "tracevalue1111111")
	logger.InfowCtx(ctx, "test", Any("123", 123))
	InfowCtx(ctx, "test", Any("123", 123), zap.Int("456", 456), "Operate", "RequestLog")
	Info("test", "Operate", "RequestLog")
}
