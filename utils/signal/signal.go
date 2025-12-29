package signal

import (
	"context"
	"net/http"
	"time"

	"github.com/Verthandii/spring/env"
	"github.com/Verthandii/spring/ilogger"
	"github.com/Verthandii/spring/mq/irocket"
)

type Closer interface {
	Close(ctx context.Context)
}

// ListenSignalAndCloseServer 信号检测并关闭服务
func ListenSignalAndCloseServer(fn func()) {
	ilogger.Info("【Signal】收到关闭信号", "signal", ListenSignal().String())
	fn()
}

// ListenSignalAndShutdown 信号检测并平滑关闭服务
func ListenSignalAndShutdown(closers ...Closer) {
	ilogger.Info("【Signal】收到关闭信号", "signal", ListenSignal().String())
	shutdown(closers)
}

var closed bool

func shutdown(closers []Closer) {
	closed = true
	irocket.CloseConsumerGroup()

	// sleep 作用: 3秒一次的健康检查，3次判定为失败
	// 当容器重启时，需要让健康检查挂掉，不在接收新的请求
	if !env.IsLocal() {
		time.Sleep(time.Second * 15)
	}

	// 创建一个 10 秒的超时上下文
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	for _, closer := range closers {
		closer.Close(ctx)
	}
	// 预留5秒给接口中带有goroutine的业务处理时间
	if !env.IsLocal() {
		time.Sleep(time.Second * 5)
	}
	ilogger.Info("【Signal】服务已全部关闭")
}

func PingCode() (int, string) {
	if closed {
		return http.StatusServiceUnavailable, "Restart"
	}
	return http.StatusOK, "Pong"
}
