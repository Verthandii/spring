package sentinel

type Config struct {
	Name string `json:"name" yaml:"name"`
	// Enabled      bool           `json:"enabled" yaml:"enabled"`
	// Type         int32          `json:"type" yaml:"type"`
	FlowControls []*FlowControl `json:"flowControls" yaml:"flowControls"`
}

type FlowControl struct {
	Rule *FlowRule `json:"rule" yaml:"rule"`
}

type FlowRule struct {
	Resource         string  `json:"resource" yaml:"resource"`
	Threshold        float64 `json:"threshold" yaml:"threshold"`
	StatIntervalInMs uint32  `json:"stat_interval_in_ms" yaml:"statIntervalInMs"`
}
