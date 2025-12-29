package irocket

import (
	"context"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"

	"github.com/Verthandii/spring/ilogger"
)

// Producer 全局-生产者实例
var Producer *iProducer

// iProducer 自定义-生产者
type iProducer struct {
	cfgMap     map[Topic]TopicConfig
	Producer   rocketmq.Producer
	sendSyncFn producerSendSyncFn
}

// SendSync 发送同步消息
// 建议使用 SendSyncCtx，可以打印链路信息
func (i *iProducer) SendSync(msg *primitive.Message) error {
	cfg := i.cfgMap[Topic(msg.Topic)]
	_, err := i.sendSyncFn(context.Background(), &cfg, msg)
	if err != nil {
		ilogger.Error(err.Error(), "Operate", "iProducer SendSync")
		return err
	}
	return nil
}

// Deprecated: SendAsync 发送异步消息，如果传入的callback为nil，会调用默认回调函数 defaultCallback，否者请传入需要回调的函数
// 建议使用 SendAsyncCtx，可以打印链路信息
func (i *iProducer) SendAsync(msg *primitive.Message, callback func(ctx context.Context, result *primitive.SendResult, err error)) error {
	if callback == nil {
		callback = defaultCallback
	}
	err := i.Producer.SendAsync(context.Background(), callback, msg)
	if err != nil {
		ilogger.Error(err.Error(), "Operate", "iProducer SendAsync")
		return err
	}
	return nil
}

// SendSyncCtx 发送同步消息
func (i *iProducer) SendSyncCtx(ctx context.Context, msg *primitive.Message) (*primitive.SendResult, error) {
	cfg := i.cfgMap[Topic(msg.Topic)]
	res, err := i.sendSyncFn(ctx, &cfg, msg)
	if err != nil {
		ilogger.ErrorwCtx(ctx, err.Error(), "Operate", "iProducer SendSync")
		return nil, err
	}
	return res, nil
}

// SendAsyncCtx 发送异步消息，如果传入的callback为nil，会调用默认回调函数 defaultCallback，否者请传入需要回调的函数
func (i *iProducer) SendAsyncCtx(ctx context.Context, msg *primitive.Message, callback func(ctx context.Context, result *primitive.SendResult, err error)) error {
	if callback == nil {
		callback = defaultCallback
	}
	err := i.Producer.SendAsync(ctx, callback, msg)
	if err != nil {
		ilogger.ErrorwCtx(ctx, err.Error(), "Operate", "iProducer SendAsync")
		return err
	}
	return nil
}

// defaultCallback 异步消息-默认回调函数
func defaultCallback(ctx context.Context, result *primitive.SendResult, err error) {
	if err != nil {
		ilogger.ErrorwCtx(ctx, err.Error(), "Operate", "defaultCallback", "Result", result)
	}
}

// InitIProducer 初始化生产者
func InitIProducer(rc RocketConfig) {
	// 初始化配置
	SetConfig(rc)
	// 创建生产者
	pd, err := rocketmq.NewProducer(
		producer.WithNsResolver(primitive.NewPassthroughResolver(conf.Endpoint)),
		producer.WithRetry(3),
		producer.WithCredentials(primitive.Credentials{
			AccessKey: conf.AccessKey,
			SecretKey: conf.SecretKey,
		}),
	)
	if err != nil {
		ilogger.Fatal("InitIProducer rocketmq.NewProducer", "err", err)
	}
	// 启动生产者
	err = pd.Start()
	if err != nil {
		ilogger.Fatal("InitIProducer pd.Start", "err", err)
	}
	ilogger.Info("【RocketMQ】生产者启动成功")
	// 分配实例
	Producer = &iProducer{
		Producer:   pd,
		cfgMap:     rc.Topics,
		sendSyncFn: sendSyncWarpper(defaultProducerInterceptors, pd.SendSync),
	}
}
