package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var enabled = false

func IsEnabled() bool {
	return enabled
}

func Init(cfg *Config) {
	if cfg.Enable {
		enabled = true
		http.Handle(cfg.Path, promhttp.Handler())
		go http.ListenAndServe(cfg.Addr, nil)
	}
}

func update(fn func()) {
	if !enabled {
		return
	}

	fn()
}
