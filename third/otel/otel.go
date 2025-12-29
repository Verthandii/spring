package otel

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"

	"github.com/Verthandii/spring/env"
	"github.com/Verthandii/spring/ilogger"
)

func init() {
	r, _ := resource.New(context.Background(),
		resource.WithFromEnv(),   // pull attributes from OTEL_RESOURCE_ATTRIBUTES and OTEL_SERVICE_NAME environment variables
		resource.WithProcess(),   // This option configures a set of Detectors that discover process information
		resource.WithOS(),        // This option configures a set of Detectors that discover OS information
		resource.WithContainer(), // This option configures a set of Detectors that discover container information
		resource.WithHost(),      // This option configures a set of Detectors that discover host information
		resource.WithAttributes(semconv.ServiceNameKey.String(env.GetServer())),
	)

	// 暂时没有链路后端，不需要exporter
	provider := sdktrace.NewTracerProvider(
		sdktrace.WithResource(r),
	)
	otel.SetTracerProvider(provider)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		// propagation.Baggage{},
	))

	// 暂时没有后端，不需要存储
	// provider.Shutdown(context.Background())
	ilogger.Register("TID", id)
}

// SetRandTraceId 设置自定义的 traceId
// 方便在日志中串起整个链路
func SetRandTraceId(ctx context.Context) context.Context {
	ctx, ospan := otel.GetTracerProvider().Tracer(env.GetServer()).Start(ctx, "ctxkit.SetRandTraceId")
	defer ospan.End()
	return ctx
}

// SetTraceId 设置自定义的 traceId
func SetTraceId(ctx context.Context, traceId string) context.Context {
	// 创建自定义 Trace ID
	tid, err := trace.TraceIDFromHex(traceId)
	if err != nil {
		return SetRandTraceId(ctx)
	}

	// 用自定义 SpanContext 创建新的 Context
	spanCtx := trace.ContextWithSpanContext(ctx, trace.NewSpanContext(trace.SpanContextConfig{
		TraceID:    tid,
		SpanID:     trace.SpanID{},
		TraceFlags: trace.FlagsSampled,
	}))

	// 使用上面创建的 Context 启动新的 Span
	ctx, span := otel.Tracer(env.GetServer()).Start(spanCtx, "default")
	defer span.End()
	return ctx
}

// ID 从 ctx 中获取 traceId
func ID(ctx context.Context) string {
	if span := trace.SpanContextFromContext(ctx); span.HasTraceID() {
		return span.TraceID().String()
	}
	return ""
}

func id(ctx context.Context) interface{} {
	return ID(ctx)
}

// SpanID returns a spanid valuer.
func SpanID() func(ctx context.Context) interface{} {
	return func(ctx context.Context) interface{} {
		if span := trace.SpanContextFromContext(ctx); span.HasSpanID() {
			return span.SpanID().String()
		}
		return ""
	}
}

// Clone 复制 preCtx 中对应 key 的值，移除父级 cancel。
func Clone(preCtx context.Context) context.Context {
	newCtx := context.Background()

	// 从 preCtx 开启一个子 span，来传递 traceId
	_, ospan := otel.GetTracerProvider().
		Tracer(env.GetServer()).
		Start(preCtx, "ctxkit.Clone")
	defer ospan.End()
	newCtx = trace.ContextWithSpan(newCtx, ospan)
	return newCtx
}
