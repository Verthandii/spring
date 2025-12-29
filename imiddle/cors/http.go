package cors

import (
	"net/http"

	"github.com/Verthandii/spring/transport/ihttp"
	"github.com/Verthandii/spring/utils"
)

func (m *CORS) HTTP() ihttp.Handler {
	return func(c *ihttp.Context) {
		if utils.InSlice("*", m.origin) {
			allowOrigin := "*"
			origin := c.GetHeader("Origin")
			if m.originReflection && origin != "" {
				allowOrigin = origin
				c.Writer.Header().Add("Vary", "Origin")
			}
			c.Writer.Header().Add("Access-Control-Allow-Origin", allowOrigin)
		} else {
			// 无 * 设置对应的域名
			curOrigin := c.GetHeader("Origin")
			if utils.InSlice(curOrigin, m.origin) {
				// 为了不暴露过多的域名，有符合条件的就直接写上
				c.Writer.Header().Add("Access-Control-Allow-Origin", curOrigin)
			} else {
				// 没有符合条件的，写所有
				for _, origin := range m.origin {
					c.Writer.Header().Add("Access-Control-Allow-Origin", origin)
				}
			}
		}

		if m.credentials == "true" {
			// 该字段只能设置为 true
			c.Writer.Header().Set("Access-Control-Allow-Credentials", m.credentials)
		}
		c.Writer.Header().Set("Access-Control-Allow-Headers", m.headers)
		c.Writer.Header().Set("Access-Control-Allow-Methods", m.methods)
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusOK)
		}
	}
}
