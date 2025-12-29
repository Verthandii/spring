package ihttp

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Verthandii/spring/env"
	"github.com/Verthandii/spring/ilogger"
	"github.com/Verthandii/spring/utils"
	"github.com/Verthandii/spring/utils/signal"
)

// SiteIndex ingress glb 心跳检查
func SiteIndex(c *Context) {
	c.Success("Hello world")
}

// SitePing ingress glb 心跳检查
func SitePing(c *Context) {
	code, str := signal.PingCode()
	c.String(code, str)
}

// SiteEnvInfo 服务器信息
func SiteEnvInfo(c *Context) {
	c.Success(map[string]interface{}{
		"IP":        c.ClientIP(),
		"Env":       env.GetEnv(),
		"Server":    env.GetServer(),
		"Version":   env.GetVersion(),
		"BuildTime": utils.DateTime(env.GetBuildAt()),
		"Time":      utils.DateTime(time.Now().Unix()),
	})
}

// SiteTestLog 测试日志
func SiteTestLog(c *Context) {
	ilogger.InfowCtx(c.Request.Context(), "这是一条测试Info消息", "Operate", "Test1")
	ilogger.WarnwCtx(c.Request.Context(), "这是一条测试Warn消息", "Operate", "Test2")
	ilogger.ErrorwCtx(c.Request.Context(), "这是一条测试Error消息", "Operate", "Test3")
	num := 1
	num = 0
	num = 1 / num
}

// SiteSetTimeZone 设置当前容器时区
func SiteSetTimeZone(c *Context) {
	if env.IsStagingOrProd() {
		c.String(http.StatusBadRequest, "该环境不允许此操作")
		return
	}

	var bind struct {
		Time string `json:"Time"`
	}
	if err := c.ShouldBind(&bind); err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("参数有误:%s", err))
		return
	}

	layout := "2006-01-02 15:04:05"
	bjTZ := time.FixedZone("Asia/Shanghai", 3600*8)
	if bind.Time == "" {
		time.Local = bjTZ
		c.String(http.StatusOK, fmt.Sprintf("系统时间已还原,当前时间:%s", utils.DateTime(time.Now().Unix())))
		return
	}

	t, err := time.ParseInLocation(layout, bind.Time, bjTZ)
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("时间格式错误:%s", err))
		return
	}
	fixed := t.Unix() - time.Now().Unix() + 3600*8
	time.Local = time.FixedZone("Asia/Shanghai", int(fixed))
	c.String(http.StatusOK, fmt.Sprintf("系统时间已调整,当前时间:%s", utils.DateTime(time.Now().Unix())))
}
