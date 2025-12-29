package sentinel

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/base"
	"github.com/alibaba/sentinel-golang/util"
	"gopkg.in/yaml.v3"
)

func TestFlowControl(t *testing.T) {
	cs := `name: example-api
flowControls:
  - rule:
      resource: "POST:/example/api/v1/order/create"
      threshold: 100
  - rule:
      resource: "POST:/example/api/v1/activity/list"
      threshold: 1
      stat_interval_in_ms: 10000`
	cfg := &Config{}
	err := yaml.Unmarshal([]byte(cs), cfg)
	if err != nil {
		t.Error(err)
	}
	Init(cfg)

	for i := 0; i < 10; i++ {
		go func() {
			for {
				e, b := api.Entry("POST:/example/api/v1/activity/list", api.WithResourceType(base.ResTypeWeb),
					api.WithTrafficType(base.Inbound))
				if b != nil {
					// Blocked. We could get the block reason from the BlockError.
					fmt.Println(util.CurrentTimeMillis(), "blocked")
					time.Sleep(time.Duration(rand.Uint64()%10) * time.Millisecond)
				} else {
					// Passed, wrap the logic here.
					fmt.Println(util.CurrentTimeMillis(), "passed")
					// time.Sleep(time.Duration(rand.Uint64()%10) * time.Millisecond)

					// Be sure the entry is exited finally.
					e.Exit()
				}

			}
		}()
	}
	time.Sleep(time.Second * 5)
}
