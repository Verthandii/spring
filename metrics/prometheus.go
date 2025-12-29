package metrics

import "github.com/prometheus/client_golang/prometheus"

type counter struct {
	cv  *prometheus.CounterVec
	lvs []string
}

func NewCounter(opts prometheus.CounterOpts, labelNames []string) Counter {
	c := &counter{
		cv: prometheus.NewCounterVec(opts, labelNames),
	}
	prometheus.MustRegister(c.cv)
	return c
}

func (c *counter) WithLabelValues(labelValues ...string) Counter {
	return &counter{
		cv:  c.cv,
		lvs: append(c.lvs, labelValues...),
	}
}

func (c *counter) Add(delta float64) {
	update(func() {
		c.cv.WithLabelValues(c.lvs...).Add(delta)
	})
}

func (c *counter) Inc() {
	update(func() {
		c.cv.WithLabelValues(c.lvs...).Add(1)
	})
}

type gauge struct {
	gv  *prometheus.GaugeVec
	lvs []string
}

func NewGauge(opts prometheus.GaugeOpts, labelNames []string) Gauge {
	g := &gauge{
		gv: prometheus.NewGaugeVec(opts, labelNames),
	}
	prometheus.MustRegister(g.gv)
	return g
}

func (g *gauge) WithLabelValues(labelValues ...string) Gauge {
	return &gauge{
		gv:  g.gv,
		lvs: append(g.lvs, labelValues...),
	}
}

func (g *gauge) Set(value float64) {
	update(func() {
		g.gv.WithLabelValues(g.lvs...).Set(value)
	})
}

func (g *gauge) Add(delta float64) {
	update(func() {
		g.gv.WithLabelValues(g.lvs...).Add(delta)
	})
}

func (g *gauge) Inc() {
	update(func() {
		g.gv.WithLabelValues(g.lvs...).Add(1)
	})
}

type summary struct {
	sv  *prometheus.SummaryVec
	lvs []string
}

func NewSummary(opts prometheus.SummaryOpts, labelNames []string) Summary {
	s := &summary{
		sv: prometheus.NewSummaryVec(opts, labelNames),
	}
	prometheus.MustRegister(s.sv)
	return s
}

func (s *summary) WithLabelValues(labelValues ...string) Summary {
	return &summary{
		sv:  s.sv,
		lvs: append(s.lvs, labelValues...),
	}
}

func (s *summary) Observe(value float64) {
	update(func() {
		s.sv.WithLabelValues(s.lvs...).Observe(value)
	})
}

type histogram struct {
	hv  *prometheus.HistogramVec
	lvs []string
}

func NewHistogram(opts prometheus.HistogramOpts, labelNames []string) Histogram {
	h := &histogram{
		hv: prometheus.NewHistogramVec(opts, labelNames),
	}
	prometheus.MustRegister(h.hv)
	return h
}

func (h *histogram) WithLabelValues(labelValues ...string) Histogram {
	return &histogram{
		hv:  h.hv,
		lvs: append(h.lvs, labelValues...),
	}
}

func (h *histogram) Observe(value float64) {
	update(func() {
		h.hv.WithLabelValues(h.lvs...).Observe(value)
	})
}
