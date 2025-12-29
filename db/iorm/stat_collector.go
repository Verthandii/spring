package iorm

import (
	"fmt"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"gorm.io/gorm"

	"github.com/Verthandii/spring/metrics"
)

var connLabels = []string{"addr", "db_name"}

// statCollector gorm状态采集器
type statCollector struct {
	lock    sync.Mutex
	clients []*statClient

	maxOpenConnections *prometheus.Desc // Maximum number of open connections to the database.
	openConnections    *prometheus.Desc // The number of established connections both in use and idle.
	inUseConnections   *prometheus.Desc // The number of connections currently in use.
	idleConnections    *prometheus.Desc // The number of idle connections.
	waitCount          *prometheus.Desc // The total number of connections waited for.
	waitDuration       *prometheus.Desc // The total time blocked waiting for a new connection.
	maxIdleClosed      *prometheus.Desc // The total number of connections closed due to maxIdleClosed.
	maxIdleTimeClosed  *prometheus.Desc // The total number of connections closed due to maxIdleTimeClosed.
	maxLifetimeClosed  *prometheus.Desc // The total number of connections closed due to maxLifetimeClosed.
}

type statClient struct {
	addr   string
	dbName string
	inner  *gorm.DB
}

func newStatCollector() *statCollector {
	fqName := func(name string) string {
		return fmt.Sprintf("%s_%s_%s", metrics.Namespace, subNamespace, name)
	}
	collector := &statCollector{
		maxOpenConnections: prometheus.NewDesc(
			fqName("max_open_connections"),
			"Maximum number of open connections to the database.",
			connLabels, nil,
		),
		openConnections: prometheus.NewDesc(
			fqName("open_connections"),
			"The number of established connections both in use and idle.",
			connLabels, nil,
		),
		inUseConnections: prometheus.NewDesc(
			fqName("in_use_connections"),
			"The number of connections currently in use.",
			connLabels, nil,
		),
		idleConnections: prometheus.NewDesc(
			fqName("idle_connections"),
			"The number of idle connections.",
			connLabels, nil,
		),
		waitCount: prometheus.NewDesc(
			fqName("wait_count_total"),
			"The total number of connections waited for.",
			connLabels, nil,
		),
		waitDuration: prometheus.NewDesc(
			fqName("wait_duration_seconds_total"),
			"The total time blocked waiting for a new connection.",
			connLabels, nil,
		),
		maxIdleClosed: prometheus.NewDesc(
			fqName("max_idle_closed_total"),
			"The total number of connections closed due to SetMaxIdleConns.",
			connLabels, nil,
		),
		maxIdleTimeClosed: prometheus.NewDesc(
			fqName("max_idle_time_closed_total"),
			"The total number of connections closed due to SetConnMaxIdleTime.",
			connLabels, nil,
		),
		maxLifetimeClosed: prometheus.NewDesc(
			fqName("max_lifetime_closed_total"),
			"The total number of connections closed due to SetConnMaxLifetime.",
			connLabels, nil,
		),
	}

	prometheus.MustRegister(collector)
	return collector
}

func (c *statCollector) Describe(descs chan<- *prometheus.Desc) {
	descs <- c.maxOpenConnections
	descs <- c.openConnections
	descs <- c.inUseConnections
	descs <- c.idleConnections
	descs <- c.waitCount
	descs <- c.waitDuration
	descs <- c.maxIdleClosed
	descs <- c.maxLifetimeClosed
	descs <- c.maxIdleTimeClosed
}

// Collect 由 prometheus 调用，导出 gorm 状态数据
func (c *statCollector) Collect(metrics chan<- prometheus.Metric) {
	c.lock.Lock()
	defer c.lock.Unlock()

	for _, client := range c.clients {
		addr := client.addr
		dbName := client.dbName
		db, err := client.inner.DB()
		if err != nil {
			continue
		}
		stats := db.Stats()

		metrics <- prometheus.MustNewConstMetric(c.maxOpenConnections, prometheus.GaugeValue,
			float64(stats.MaxOpenConnections), addr, dbName)
		metrics <- prometheus.MustNewConstMetric(c.openConnections, prometheus.GaugeValue,
			float64(stats.OpenConnections), addr, dbName)
		metrics <- prometheus.MustNewConstMetric(c.inUseConnections, prometheus.GaugeValue,
			float64(stats.InUse), addr, dbName)
		metrics <- prometheus.MustNewConstMetric(c.idleConnections, prometheus.GaugeValue,
			float64(stats.Idle), addr, dbName)
		metrics <- prometheus.MustNewConstMetric(c.waitCount, prometheus.CounterValue,
			float64(stats.WaitCount), addr, dbName)
		metrics <- prometheus.MustNewConstMetric(c.waitDuration, prometheus.CounterValue,
			stats.WaitDuration.Seconds(), addr, dbName)
		metrics <- prometheus.MustNewConstMetric(c.maxIdleClosed, prometheus.CounterValue,
			float64(stats.MaxIdleClosed), addr, dbName)
		metrics <- prometheus.MustNewConstMetric(c.maxLifetimeClosed, prometheus.CounterValue,
			float64(stats.MaxLifetimeClosed), addr, dbName)
		metrics <- prometheus.MustNewConstMetric(c.maxIdleTimeClosed, prometheus.CounterValue,
			float64(stats.MaxIdleTimeClosed), addr, dbName)
	}
}

func (c *statCollector) registerClient(client *statClient) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.clients = append(c.clients, client)
}
