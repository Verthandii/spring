package iorm

import (
	"github.com/prometheus/client_golang/prometheus"

	"github.com/Verthandii/spring/metrics"
)

const subNamespace = "mysql"

var (
	metricReqDur = metrics.NewHistogram(prometheus.HistogramOpts{
		Namespace: metrics.Namespace,
		Subsystem: subNamespace,
		Name:      "request_duration_ms",
		Help:      "The duration of requests in ms.",
		Buckets:   []float64{0.25, 0.5, 1, 1.5, 2, 3, 5, 10, 25, 50, 100, 250, 500, 1000, 2000, 5000, 10000, 15000},
	}, []string{"addr", "db_name", "table", "operation"})

	_              prometheus.Collector = (*statCollector)(nil)
	_statCollector                      = newStatCollector()
)
