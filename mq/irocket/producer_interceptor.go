package irocket

import (
	"context"
	"time"

	"github.com/apache/rocketmq-client-go/v2/primitive"
)

// defaultProducerInterceptors 默认的生产者拦截器
var defaultProducerInterceptors = []mqIntceptor[[]*primitive.Message, *primitive.SendResult]{
	metricsProducerInterceptor,
}

func metricsProducerInterceptor(ctx context.Context, cfg *TopicConfig, msgs []*primitive.Message,
	invoker mqInvoker[[]*primitive.Message, *primitive.SendResult]) (*primitive.SendResult, error) {

	startAt := time.Now()
	result, err := invoker(ctx, msgs)
	elapsed := time.Since(startAt)

	producerDur.WithLabelValues(cfg.Topic, cfg.Group, cfg.Tag).Observe(float64(elapsed.Milliseconds()))
	return result, err
}
