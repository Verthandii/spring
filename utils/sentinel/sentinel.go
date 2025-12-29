package sentinel

import (
	"github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/config"
	"github.com/alibaba/sentinel-golang/core/flow"

	"github.com/Verthandii/spring/ilogger"
)

func Init(cfg *Config) {
	if cfg == nil {
		ilogger.Warn("sentinel config is nil")
		return
	}
	if cfg.Name == "" {
		ilogger.Warn("sentinel name is empty")
		return
	}
	scfg := config.NewDefaultConfig()
	scfg.Sentinel.App.Name = cfg.Name
	err := api.InitWithConfig(scfg)
	if err != nil {
		ilogger.Warn("sentinel init failed", ilogger.Any("err", err))
		return
	}
	rules := []*flow.Rule{}
	for _, flowRule := range cfg.FlowControls {
		if flowRule.Rule != nil {
			rules = append(rules, &flow.Rule{
				Resource:         flowRule.Rule.Resource,
				Threshold:        flowRule.Rule.Threshold,
				StatIntervalInMs: flowRule.Rule.StatIntervalInMs,
			})
		}
	}
	// 注册流控
	ok, err := flow.LoadRules(rules)
	if err != nil {
		ilogger.Warn("sentinel flow rule init failed", ilogger.Any("err", err))
		return
	}
	if !ok {
		ilogger.Warn("sentinel flow rule init failed")
		return
	}
}
