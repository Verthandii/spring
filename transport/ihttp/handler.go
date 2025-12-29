package ihttp

import (
	"github.com/gin-gonic/gin"
)

type Handler func(c *Context)

func (s *Server) H(h Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		h(newContext(s, c))
	}
}
