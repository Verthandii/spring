package irocket

import (
	"context"
	"time"

	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"

	"github.com/Verthandii/spring/third/otel"
)

// defaultConsumerInterceptors 默认的消费者拦截器
var defaultConsumerInterceptors = []mqIntceptor[[]*primitive.MessageExt, consumer.ConsumeResult]{
	metricsConsumerInterceptor,
	setTraceIdConsumerInterceptor,
}

const (
	msgConsumerSuccess = "0" // 消费成功
	msgConsumerFailed  = "1" // 消费失败
)

func metricsConsumerInterceptor(ctx context.Context, cfg *TopicConfig,
	msg []*primitive.MessageExt, invoker mqInvoker[[]*primitive.MessageExt, consumer.ConsumeResult]) (consumer.ConsumeResult, error) {

	startAt := time.Now()
	result, err := invoker(ctx, msg)
	elapsed := time.Since(startAt)
	errorType := msgConsumerSuccess
	if err != nil {
		errorType = msgConsumerFailed
	}
	consumeDur.WithLabelValues(cfg.Topic, cfg.Group, cfg.Tag, errorType).Observe(float64(elapsed.Milliseconds()))
	return result, err
}

func setTraceIdConsumerInterceptor(ctx context.Context, cfg *TopicConfig,
	msg []*primitive.MessageExt, invoker mqInvoker[[]*primitive.MessageExt, consumer.ConsumeResult]) (consumer.ConsumeResult, error) {

	ctx = otel.SetRandTraceId(ctx)
	return invoker(ctx, msg)
}
