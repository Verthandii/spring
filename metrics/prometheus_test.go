package metrics

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	mg "github.com/prometheus/client_model/go"
	"github.com/stretchr/testify/assert"
)

func TestCounter(t *testing.T) {
	c := NewCounter(prometheus.CounterOpts{Name: "test_counter", Help: "test counter"}, []string{"method", "code"})

	c.WithLabelValues("POST", "404").Add(1)
	value := make(chan prometheus.Metric, 1)
	cObj := c.(*counter)
	cObj.cv.Collect(value)
	m := <-value
	data := mg.Metric{}
	m.Write(&data)
	t.Logf("%+v", &data)
	assert.Equal(t, 1.0, data.GetCounter().GetValue())
}

func TestGuage(t *testing.T) {
	g := NewGauge(prometheus.GaugeOpts{Name: "test_gauge", Help: "test gauge"}, []string{"dog_name"})
	g.WithLabelValues("xifan").Set(3)
	value := make(chan prometheus.Metric, 1)
	gObj := g.(*gauge)
	gObj.gv.Collect(value)
	m := <-value
	data := mg.Metric{}
	m.Write(&data)
	t.Logf("%+v", &data)
	assert.Equal(t, 3.0, data.GetGauge().GetValue())
}

func TestHistogram(t *testing.T) {
	h := NewHistogram(prometheus.HistogramOpts{Name: "test_histogram", Help: "test histogram"}, []string{"host"})
	h.WithLabelValues("127.0.0.1").Observe(10.0)
	value := make(chan prometheus.Metric, 1)
	hObj := h.(*histogram)
	hObj.hv.Collect(value)
	m := <-value
	data := mg.Metric{}
	m.Write(&data)
	t.Logf("%+v", &data)
	assert.Equal(t, 10.0, *data.GetHistogram().SampleSum)
	assert.Equal(t, uint64(1), *data.GetHistogram().SampleCount)
}

func TestSummary(t *testing.T) {
	s := NewSummary(prometheus.SummaryOpts{Name: "test_summary", Help: "test summary"}, []string{"host"})
	s.WithLabelValues("127.0.0.1").Observe(10.0)
	value := make(chan prometheus.Metric, 1)
	sObj := s.(*summary)
	sObj.sv.Collect(value)
	m := <-value
	data := mg.Metric{}
	m.Write(&data)
	t.Logf("%+v", &data)
	assert.Equal(t, 10.0, *data.GetSummary().SampleSum)
	assert.Equal(t, uint64(1), *data.GetSummary().SampleCount)
}
