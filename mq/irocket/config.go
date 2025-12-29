package irocket

import "sync"

type Topic string

var (
	conf = RocketConfig{} // mq配置信息
	once = new(sync.Once)
)

// RocketConfig 队列配置
type RocketConfig struct {
	Endpoint  []string
	AccessKey string
	SecretKey string
	Topics    map[Topic]TopicConfig
}

// TopicConfig 队列主题
type TopicConfig struct {
	Topic string
	Group string
	Tag   string
}

func SetConfig(rc RocketConfig) {
	once.Do(func() {
		conf = rc
	})
}
