package irocket

import (
	"context"

	"github.com/apache/rocketmq-client-go/v2/primitive"
)

type producerSendSyncFn func(ctx context.Context, cfg *TopicConfig, msg ...*primitive.Message) (*primitive.SendResult, error)

type mqInvoker[T any, R any] func(ctx context.Context, msg T) (R, error)

// mq 拦截器，producer 和 consumer 通用
type mqIntceptor[T any, R any] func(ctx context.Context, cfg *TopicConfig, msg T, invoker mqInvoker[T, R]) (R, error)

// SendSync 同步消息的包装器
func sendSyncWarpper(interceptors []mqIntceptor[[]*primitive.Message, *primitive.SendResult],
	invoker func(context.Context, ...*primitive.Message) (*primitive.SendResult, error)) producerSendSyncFn {

	interceptor := chainInterceptors(interceptors...)
	return func(ctx context.Context, cfg *TopicConfig, msg ...*primitive.Message) (*primitive.SendResult, error) {
		return interceptor(ctx, cfg, msg, func(ctx context.Context, msg []*primitive.Message) (*primitive.SendResult, error) {
			return invoker(ctx, msg...)
		})
	}
}

// 串联所有的拦截器
func chainInterceptors[T any, R any](interceptors ...mqIntceptor[T, R]) mqIntceptor[T, R] {
	n := len(interceptors)
	if n > 1 {
		lastI := n - 1
		return func(ctx context.Context, cfg *TopicConfig, msg T, finalInvoker mqInvoker[T, R]) (R, error) {
			var chainInvoker mqInvoker[T, R]
			var curI int
			chainInvoker = func(ctx context.Context, msg T) (R, error) {
				// 最后一个拦截器是最终的目标invoker
				if curI == lastI {
					return finalInvoker(ctx, msg)
				}
				// 循环执行下一个拦截器
				curI++
				result, err := interceptors[curI](ctx, cfg, msg, chainInvoker)
				curI--
				return result, err
			}
			return interceptors[0](ctx, cfg, msg, chainInvoker)
		}
	}
	return func(ctx context.Context, cfg *TopicConfig, msg T, invoker mqInvoker[T, R]) (R, error) {
		return interceptors[0](ctx, cfg, msg, invoker)
	}
}
