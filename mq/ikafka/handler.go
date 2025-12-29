package ikafka

import (
	"github.com/IBM/sarama"

	"github.com/Verthandii/spring/ilogger"
)

type ConsumerGroupHandler func(*sarama.ConsumerMessage) error

func (ConsumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (ConsumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

func (h ConsumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		if err := h(msg); err == nil {
			sess.MarkMessage(msg, "")
		} else {
			ilogger.ErrorwCtx(sess.Context(), "消息处理失败",
				"err", err,
				"topic", msg.Topic,
				"value", string(msg.Value),
			)
		}
	}
	return nil
}
