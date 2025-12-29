package main

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/Verthandii/spring/metrics"
)

var http_requests_total = metrics.NewCounter(prometheus.CounterOpts{
	Name: "http_requests_total",
	Help: "Total number of HTTP requests"},
	[]string{"path", "code"})

func main() {
	metrics.Init(&metrics.Config{Path: "/metrics", Addr: ":8080"})
	go func() {
		for {
			http_requests_total.WithLabelValues("/example", "200").Add(1)
			time.Sleep(time.Second)
		}
	}()
	select {}
}
