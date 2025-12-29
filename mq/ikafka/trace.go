package ikafka

import (
	"github.com/IBM/sarama"
)

func GetTraceId(msg *sarama.ConsumerMessage) string {
	for _, header := range msg.Headers {
		if string(header.Key) == "trace_id" {
			return string(header.Value)
		}
	}
	return ""
}
