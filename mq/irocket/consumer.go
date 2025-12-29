package irocket

import (
	"context"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"

	"github.com/Verthandii/spring/ilogger"
)

var (
	// ConsumerGroup 消费者分组
	ConsumerGroup = make(map[string]rocketmq.PushConsumer)
	// SubscribeTopics 订阅Topic集合
	SubscribeTopics = make(map[string]rocketmq.PushConsumer)
)

// InitConsumerGroup 初始化消费者分组
func InitConsumerGroup(name Topic, option ...consumer.Option) {
	topic, ex := conf.Topics[name]
	if !ex {
		ilogger.Error("TopicConfig not found", "Operate", "FindTopic", "Name", name)
		return
	}
	option = append(option, consumer.WithGroupName(topic.Group))
	option = append(option, consumer.WithInterceptor(WithRecoverInterceptor()))
	option = append(option, consumer.WithNsResolver(primitive.NewPassthroughResolver(conf.Endpoint)))
	option = append(option, consumer.WithCredentials(primitive.Credentials{AccessKey: conf.AccessKey, SecretKey: conf.SecretKey}))
	c, err := rocketmq.NewPushConsumer(option...)
	if err != nil {
		ilogger.Error(err.Error(), "Operate", "NewPushConsumer", "TopicConfig", topic)
		return
	}
	ConsumerGroup[topic.Group] = c
}

// StartConsumerGroup 开启消费者分组
func StartConsumerGroup() {
	for _, cg := range ConsumerGroup {
		if err := cg.Start(); err != nil {
			ilogger.Error(err.Error(), "Operate", "Start")
			return
		}
	}
}

// CloseConsumerGroup 关闭消费者分组
func CloseConsumerGroup() {
	if len(SubscribeTopics) == 0 {
		return
	}

	ilogger.Info("【RocketMQ】取消订阅所有Topic")
	for topic, con := range SubscribeTopics {
		_ = con.Unsubscribe(topic)
	}
}

// Subscribe 消息订阅
func Subscribe(name Topic, callback func(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error)) {
	topic, ex := conf.Topics[name]
	if !ex {
		ilogger.Error("【RocketMQ】TopicConfig not found", "Operate", "FindTopic", "Name", name)
		return
	}
	// 获取消费者实例
	c, ex := ConsumerGroup[topic.Group]
	if !ex {
		ilogger.Error("【RocketMQ】Group not found", "Operate", "FindGroup", "Name", name)
		return
	}

	interceptor := chainInterceptors(defaultConsumerInterceptors...)
	// 开始订阅消息
	err := c.Subscribe(topic.Topic, consumer.MessageSelector{}, func(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
		return interceptor(ctx, &topic, msgs, msgHandlerWrapper(callback))
	})
	if err != nil {
		ilogger.Error(err.Error(), "Operate", "Subscribe", "TopicConfig", topic)
		return
	}
	SubscribeTopics[topic.Topic] = c
	ilogger.Info("【RocketMQ】Subscribe", "Group", topic.Group, "Topic", topic.Topic)
}

func msgHandlerWrapper(handler func(context.Context, ...*primitive.MessageExt) (consumer.ConsumeResult, error)) mqInvoker[[]*primitive.MessageExt, consumer.ConsumeResult] {
	return func(ctx context.Context, msgs []*primitive.MessageExt) (consumer.ConsumeResult, error) {
		return handler(ctx, msgs...)
	}
}

// WithRecoverInterceptor 捕获异常并重试消息
func WithRecoverInterceptor() primitive.Interceptor {
	return func(ctx context.Context, req, reply interface{}, next primitive.Invoker) error {
		defer func() {
			if err := recover(); err != nil {
				if r, ok := reply.(*consumer.ConsumeResultHolder); ok {
					r.ConsumeResult = consumer.ConsumeRetryLater
				}
				ilogger.Error("【RocketMQ】 consumer panic", "err", err, "msgs", req)
			}
		}()
		return next(ctx, req, reply)
	}
}

func WithMetricsInterceptor() primitive.Interceptor {
	return func(ctx context.Context, req, reply interface{}, next primitive.Invoker) error {
		return next(ctx, req, reply)
	}
}
