package irocket

import (
	"github.com/apache/rocketmq-client-go/v2/primitive"

	"github.com/Verthandii/spring/ilogger"
)

// NewMessage 创建信息
func NewMessage(name Topic, data []byte, tag string) *primitive.Message {
	topic, ex := conf.Topics[name]
	if !ex {
		ilogger.Error("TopicConfig is Nil", "Operate", "FindTopic", "Name", name)
		return nil
	}
	msg := primitive.NewMessage(topic.Topic, data)
	if tag != "" {
		msg.WithTag(tag)
	}
	return msg
}
