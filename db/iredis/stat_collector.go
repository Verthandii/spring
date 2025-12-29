package iredis

import (
	"fmt"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/redis/go-redis/v9"

	"github.com/Verthandii/spring/metrics"
)

// statCollector redis 状态采集器
type statCollector struct {
	lock        sync.Mutex
	clients     []*statClient
	hitDesc     *prometheus.Desc // Number of times a connection was found in the pool
	missDesc    *prometheus.Desc // Number of times a connection was not found in the pool
	timeoutDesc *prometheus.Desc // Number of times a timeout occurred when looking for a connection in the pool
	totalDesc   *prometheus.Desc // Total number of connections
	idleDesc    *prometheus.Desc // Number of idle connections
	staleDesc   *prometheus.Desc // Number of stale connections
	maxDesc     *prometheus.Desc // Maximum number of connections
}

type statClient struct {
	addr     string
	dbName   string
	poolSize int
	inner    *RedisCli
}

var connLabels = []string{"addr", "db_name"}

func newStatCollector() *statCollector {
	fqName := func(name string) string {
		return fmt.Sprintf("%s_%s_%s", metrics.Namespace, subNamespace, name)
	}

	collector := &statCollector{
		hitDesc: prometheus.NewDesc(
			fqName("pool_hit_total"),
			"Number of times a connection was found in the pool",
			connLabels, nil,
		),
		missDesc: prometheus.NewDesc(
			fqName("pool_miss_total"),
			"Number of times a connection was not found in the pool",
			connLabels, nil,
		),
		timeoutDesc: prometheus.NewDesc(
			fqName("pool_timeout_total"),
			"Number of times a timeout occurred when looking for a connection in the pool",
			connLabels, nil,
		),
		totalDesc: prometheus.NewDesc(
			fqName("pool_conn_total_current"),
			"Current number of connections in the pool",
			connLabels, nil,
		),
		idleDesc: prometheus.NewDesc(
			fqName("pool_conn_idle_current"),
			"Current number of idle connections in the pool",
			connLabels, nil,
		),
		staleDesc: prometheus.NewDesc(
			fqName("pool_conn_stale_total"),
			"Number of times a connection was removed from the pool because it was stale",
			connLabels, nil,
		),
		maxDesc: prometheus.NewDesc(
			fqName("pool_conn_max"),
			"Max number of connections in the pool",
			connLabels, nil,
		),
	}
	prometheus.MustRegister(collector)
	return collector
}

func (s *statCollector) Describe(descs chan<- *prometheus.Desc) {
	descs <- s.hitDesc
	descs <- s.missDesc
	descs <- s.timeoutDesc
	descs <- s.totalDesc
	descs <- s.idleDesc
	descs <- s.staleDesc
	descs <- s.maxDesc
}

func (s *statCollector) Collect(metrics chan<- prometheus.Metric) {
	s.lock.Lock()
	defer s.lock.Unlock()

	for _, client := range s.clients {
		addr, dbName := client.addr, client.dbName
		var stats *redis.PoolStats
		if client.inner.Client != nil {
			stats = client.inner.Client.PoolStats()
		} else {
			stats = client.inner.Cluster.PoolStats()
		}
		metrics <- prometheus.MustNewConstMetric(
			s.hitDesc,
			prometheus.CounterValue,
			float64(stats.Hits),
			addr,
			dbName,
		)
		metrics <- prometheus.MustNewConstMetric(
			s.missDesc,
			prometheus.CounterValue,
			float64(stats.Misses),
			addr,
			dbName,
		)
		metrics <- prometheus.MustNewConstMetric(
			s.timeoutDesc,
			prometheus.CounterValue,
			float64(stats.Timeouts),
			addr,
			dbName,
		)
		metrics <- prometheus.MustNewConstMetric(
			s.totalDesc,
			prometheus.GaugeValue,
			float64(stats.TotalConns),
			addr,
			dbName,
		)
		metrics <- prometheus.MustNewConstMetric(
			s.idleDesc,
			prometheus.GaugeValue,
			float64(stats.IdleConns),
			addr,
			dbName,
		)
		metrics <- prometheus.MustNewConstMetric(
			s.staleDesc,
			prometheus.CounterValue,
			float64(stats.StaleConns),
			addr,
			dbName,
		)
		metrics <- prometheus.MustNewConstMetric(
			s.maxDesc,
			prometheus.GaugeValue,
			float64(client.poolSize),
			addr,
			dbName,
		)
	}
}

func (s *statCollector) registerClient(client *statClient) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.clients = append(s.clients, client)
}
