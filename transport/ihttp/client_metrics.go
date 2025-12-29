package ihttp

import (
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/Verthandii/spring/metrics"
)

const clientSubNamespace = "http_client"

var (
	metricClientReqDur = metrics.NewHistogram(prometheus.HistogramOpts{
		Namespace: metrics.Namespace,
		Subsystem: clientSubNamespace,
		Name:      "requests_duration_ms",
		Help:      "http client requests duration(ms).",
		Buckets:   []float64{0.25, 0.5, 1, 2, 5, 10, 25, 50, 100, 250, 500, 1000, 2000, 5000, 10000, 15000},
	}, []string{"endpoint", "method", "path"})

	metricClientReqCodeTotal = metrics.NewCounter(prometheus.CounterOpts{
		Namespace: metrics.Namespace,
		Subsystem: clientSubNamespace,
		Name:      "requests_code_total",
		Help:      "http client requests code count.",
	}, []string{"endpoint", "method", "path", "code"})
)

func roundTripperMetrics(next http.RoundTripper) promhttp.RoundTripperFunc {
	return func(r *http.Request) (*http.Response, error) {
		start := time.Now()

		resp, err := next.RoundTrip(r)

		var statusCode int
		if err != nil || resp == nil {
			statusCode = http.StatusInternalServerError
		} else {
			statusCode = resp.StatusCode
		}
		code := strconv.Itoa(statusCode)
		method := r.Method
		endpoint := r.URL.Host
		path := r.URL.Path

		metricClientReqDur.WithLabelValues(endpoint, method, path).Observe(float64(time.Since(start).Milliseconds()))
		metricClientReqCodeTotal.WithLabelValues(endpoint, method, path, code).Inc()

		return resp, err
	}
}
