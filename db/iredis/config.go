package iredis

type Config struct {
	Host           string
	Port           int
	DB             int
	Password       string
	MinIdleConns   int
	PoolSize       int
	DialTimeout    int
	ReadTimeout    int
	WriteTimeout   int
	IdleTimeout    int
	PoolTimeout    int
	DisableMetrics bool
}
