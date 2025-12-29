package metrics

type Config struct {
	Path   string `yaml:"path"`
	Addr   string `yaml:"addr"`
	Enable bool   `yaml:"enable"`
}

func NewDefaultConfig() *Config {
	return DefaultConfig
}
