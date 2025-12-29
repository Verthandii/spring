package ihttp

import (
	"fmt"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/Verthandii/spring/env"
	"github.com/Verthandii/spring/metrics"
)

const serverSubNamespace = "http_server"

var (
	metricServerReqCnt = metrics.NewCounter(prometheus.CounterOpts{
		Namespace: metrics.Namespace,
		Subsystem: serverSubNamespace,
		Name:      "requests_total",
		Help:      "How many HTTP requests processed, partitioned by status code and HTTP method.",
	}, []string{"service", "code", "status", "method", "host", "url"})

	metricServerReqDur = metrics.NewSummary(prometheus.SummaryOpts{
		Namespace:  metrics.Namespace,
		Subsystem:  serverSubNamespace,
		Name:       "requests_duration_seconds",
		Help:       "The HTTP request latencies in seconds.",
		Objectives: metrics.DefaultObjectives,
	}, []string{"service", "code", "status", "method", "host", "url"})
)

func serverMetrics(c *Context) {
	start := time.Now()
	c.Next()

	yStatus := c.Writer.Status()
	yCodeString := fmt.Sprintf("%v", yStatus)
	code := strconv.Itoa(c.Writer.Status())

	elapsed := float64(time.Since(start)) / float64(time.Second)

	url := c.FullPath()
	service := env.GetServer()
	metricServerReqDur.WithLabelValues(service, code, yCodeString, c.Request.Method, c.Request.Host, url).Observe(elapsed)
	metricServerReqCnt.WithLabelValues(service, code, yCodeString, c.Request.Method, c.Request.Host, url).Inc()
}
