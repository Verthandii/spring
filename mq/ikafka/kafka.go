package ikafka

import (
	"context"
	"crypto/tls"
	"time"

	"github.com/IBM/sarama"

	"github.com/Verthandii/spring/ilogger"
	"github.com/Verthandii/spring/third/otel"
)

type Kafka struct {
	producer sarama.SyncProducer
	client   sarama.Client
	cfg      *Config
}

type Config struct {
	Topic            string
	GroupId          string
	BootstrapServers []string
	Protocol         string
	Username         string
	Password         string
}

func NewClient(cfg *Config) (*Kafka, error) {
	scfg := buildConfig(cfg)
	return newKafkaClient(cfg, scfg)
}

func newKafkaClient(cfg *Config, scfg *sarama.Config) (*Kafka, error) {
	// 连接kafka client
	client, err := sarama.NewClient(cfg.BootstrapServers, scfg)
	if err != nil {
		return nil, err
	}

	syncProducer, err := sarama.NewSyncProducer(cfg.BootstrapServers, scfg)
	if err != nil {
		return nil, err
	}

	return &Kafka{
		producer: syncProducer,
		client:   client,
		cfg:      cfg,
	}, nil
}

// 构建kafka配置
func buildConfig(v *Config) *sarama.Config {
	cfg := sarama.NewConfig()
	cfg.Consumer.Offsets.Initial = sarama.OffsetOldest
	cfg.Producer.RequiredAcks = -1
	cfg.Producer.Return.Successes = true
	cfg.Producer.Partitioner = sarama.NewRandomPartitioner
	switch v.Protocol {
	case "plaintext":

	case "sasl_ssl":
		cfg.Net.SASL.Enable = true
		cfg.Net.SASL.User = v.Username
		cfg.Net.SASL.Password = v.Password
		cfg.Net.TLS.Enable = true
		cfg.Net.TLS.Config = &tls.Config{
			MinVersion:         tls.VersionTLS12,
			InsecureSkipVerify: true,
		}
	}
	return cfg
}

func (k *Kafka) SendMessage(ctx context.Context, value string, headers ...sarama.RecordHeader) error {
	msg := &sarama.ProducerMessage{
		Topic:     k.cfg.Topic,
		Value:     sarama.StringEncoder(value),
		Timestamp: time.Now(),
	}
	msg.Headers = append(msg.Headers, sarama.RecordHeader{
		Key:   []byte("trace_id"),
		Value: []byte(otel.ID(ctx)),
	})
	msg.Headers = append(msg.Headers, headers...)

	partition, offset, err := k.producer.SendMessage(msg)
	if err != nil {
		ilogger.ErrorwCtx(ctx, "发送kafka消息失败", "err", err, "topic", msg.Topic)
		return err
	}
	ilogger.DebugwCtx(ctx, "发送kafka消息成功",
		"topic", msg.Topic,
		"partition", partition,
		"offset", offset,
		"msg", msg,
	)
	return nil
}

func (k *Kafka) ConsumeGroup(ctx context.Context, fn ConsumerGroupHandler) error {
	consumerGroup, err := sarama.NewConsumerGroupFromClient(k.cfg.GroupId, k.client)
	if err != nil {
		return err
	}

	defer func() {
		if err := recover(); err != nil {
			ilogger.ErrorwCtx(ctx, "消费kafka信息发生panic", "panic", err)
		}
	}()

	defer func() {
		err := consumerGroup.Close()
		if err != nil {
			ilogger.ErrorwCtx(ctx, "close err", "err", err)
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			err := consumerGroup.Consume(ctx, []string{k.cfg.Topic}, fn)
			if err != nil {
				ilogger.ErrorwCtx(ctx, "消费kafka失败", "err", err)
			}
		}
	}
}
