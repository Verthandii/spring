package iredis

import (
	"context"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/redis/go-redis/v9"

	"github.com/Verthandii/spring/metrics"
)

const subNamespace = "redis"

var (
	// 请求耗时
	metricsReqDur = metrics.NewHistogram(prometheus.HistogramOpts{
		Namespace: metrics.Namespace,
		Subsystem: subNamespace,
		Name:      "requests_execute_duration_ms",
		Help:      "The request latencies in seconds.",
		Buckets:   []float64{0.25, 0.5, 1, 1.5, 2, 3, 5, 10, 25, 50, 100, 250, 500, 1000, 2000, 5000, 10000, 15000},
	}, []string{"addr", "db", "command"})

	_statCollector            = newStatCollector()
	_              redis.Hook = &durationHook{}
)

type durationHook struct {
	addr string
	host string
	port int
	db   string
}

func newDurationHook(addr string, host string, port int, db string) redis.Hook {
	return &durationHook{
		addr: addr,
		host: host,
		port: port,
		db:   db,
	}
}

func (h durationHook) DialHook(next redis.DialHook) redis.DialHook {
	return next
}

// ProcessHook 普通命令处理
func (h durationHook) ProcessHook(next redis.ProcessHook) redis.ProcessHook {
	return func(ctx context.Context, cmd redis.Cmder) error {
		start := time.Now()

		err := next(ctx, cmd)

		duration := time.Since(start)

		metricsReqDur.WithLabelValues(h.addr, h.db, cmd.Name()).Observe(float64(duration.Milliseconds()))
		return err
	}
}

// ProcessPipelineHook pipeline 命令处理
func (h durationHook) ProcessPipelineHook(next redis.ProcessPipelineHook) redis.ProcessPipelineHook {
	return func(ctx context.Context, cmds []redis.Cmder) error {
		if len(cmds) == 0 {
			return next(ctx, cmds)
		}

		start := time.Now()

		err := next(ctx, cmds)

		duration := time.Since(start)

		for _, cmd := range cmds {
			metricsReqDur.WithLabelValues(h.addr, h.db, cmd.Name()).Observe(float64(duration.Milliseconds()))
		}
		return err
	}
}
