package metrics

import "github.com/prometheus/client_golang/prometheus"

const Namespace = "pubplat"

// 默认误差
var DefaultObjectives = map[float64]float64{0.5: 0.05, 0.95: 0.005, 0.99: 0.001}

// 默认桶
var DefHistogramBuckets = prometheus.DefBuckets

// 默认配置
var DefaultConfig = &Config{
	Addr:   ":7000", // 跟运维商定的7000端口
	Path:   "/metrics",
	Enable: true,
}
