package irocket

import (
	"github.com/prometheus/client_golang/prometheus"

	"github.com/Verthandii/spring/metrics"
)

const subNamespace = "rocketmq"

var (
	producerDur = metrics.NewHistogram(prometheus.HistogramOpts{
		Namespace: metrics.Namespace,
		Subsystem: subNamespace,
		Name:      "produce_duration_ms",
		Help:      "Total number of rocketmq messages in seconds.",
		Buckets:   []float64{1, 2, 4, 10, 25, 50, 100, 250, 500, 1000, 2000, 5000, 10000, 15000},
	}, []string{"topic", "group", "tag"})
)
