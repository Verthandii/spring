package metrics

// https://go-kratos.dev/docs/component/metrics/
// Counter 最简单的计数器，对外提供了Inc, Add两个方法
type Counter interface {
	WithLabelValues(labelValues ...string) Counter
	Add(delta float64)
	Inc()
}

// Gauge 是个状态指示器，用于记录服务当前的状态，状态值可以随着时间增加或减少。通常用于监控服务当前的cpu使用率，内存使用量等。
type Gauge interface {
	WithLabelValues(labelValues ...string) Gauge
	Set(value float64)
	Add(delta float64)
	Inc()
}

// Histogram 与 Summary 的区别在于，一个是服务端计算，一个是客户端计算，这两个计算都是基于概率采样的。
// Histogram 直方图用于记录不同分桶的数量。比如不同请求耗时区间的请求数，用于指示将指标保存到了多个分桶，因此Histogram几乎无开销。
type Histogram interface {
	WithLabelValues(labelValues ...string) Histogram
	Observe(value float64)
}

// Summary则记录了不同分位的值，基于概率采样计算，比如90% 99% 分位耗时，由于需要进行额外的计算，因此对于服务有一定的开销。
type Summary interface {
	WithLabelValues(labelValues ...string) Summary
	Observe(value float64)
}
