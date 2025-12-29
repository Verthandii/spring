package ihttp

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Context struct {
	*gin.Context

	srv *Server
}

func newContext(srv *Server, c *gin.Context) *Context {
	return &Context{
		Context: c,
		srv:     srv,
	}
}

func (c *Context) Success(data ...any) {
	var data0 any
	if len(data) > 0 {
		data0 = data[0]
	}
	c.srv.successEc(c.Writer, c.Request, data0)
}

func (c *Context) Failure(err error, data ...any) {
	var data0 any
	if len(data) > 0 {
		data0 = data[0]
	}
	c.Abort()
	_ = c.Error(err)
	c.srv.errorEc(c.Writer, c.Request, err, data0, http.StatusOK)
}

func (c *Context) FailureWithCode(statusCode int, err error, data ...any) {
	var data0 any
	if len(data) > 0 {
		data0 = data[0]
	}
	c.Abort()
	_ = c.Error(err)
	c.srv.errorEc(c.Writer, c.Request, err, data0, statusCode)
}

// 定义日志数据的键名，避免魔法字符串
const logDataKey = "yostar-pubplat_log_data"

// SetLogData 设置日志数据
func (c *Context) SetLogData(key string, value any) {
	// 尝试获取现有的日志数据
	existingData, exists := c.Get(logDataKey)
	var logData map[string]any
	if exists {
		// 如果存在，转换为map类型
		logData = existingData.(map[string]any)
	} else {
		// 不存在则创建新的map
		logData = make(map[string]any)
	}
	// 设置新的键值对
	logData[key] = value
	// 将更新后的map存回context
	c.Set(logDataKey, logData)
}

// GetLogData 获取所有日志数据
func (c *Context) GetLogData() map[string]any {
	data, exists := c.Get(logDataKey)
	if !exists {
		return make(map[string]any)
	}
	// 类型断言确保返回的是map
	if logData, ok := data.(map[string]any); ok {
		return logData
	}
	return make(map[string]any)
}
