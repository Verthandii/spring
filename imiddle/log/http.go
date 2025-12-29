package log

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/Verthandii/spring/env"
	"github.com/Verthandii/spring/ierrors"
	"github.com/Verthandii/spring/ilogger"
	"github.com/Verthandii/spring/it"
	"github.com/Verthandii/spring/transport/ihttp"
)

type wrapper struct {
	gin.ResponseWriter
	data []byte
}

func (w *wrapper) Write(data []byte) (int, error) {
	w.data = bytes.Clone(data)
	return w.ResponseWriter.Write(data)
}

func HTTP() ihttp.Handler {
	return func(c *ihttp.Context) {
		if env.DisableReqLog() {
			return
		}
		req := c.Request
		if req.URL.Path == "/ping" {
			return
		}
		t := time.Now()
		requestData, _ := c.GetRawData()
		c.Request.Body = io.NopCloser(bytes.NewBuffer(requestData))
		w := &wrapper{ResponseWriter: c.Writer, data: nil}
		c.Writer = w

		defer func() {
			if r := recover(); r != nil {
				c.Failure(ierrors.InternalServerError.Format(r))
				handleLog(c, req, requestData, w, t)
				panic(r)
			}
		}()
		c.Next()

		handleLog(c, req, requestData, w, t)
	}
}

// handleLog 处理日志打印
func handleLog(c *ihttp.Context, req *http.Request, requestData []byte, w *wrapper, t time.Time) {
	now := time.Now()
	var busCode int
	if e := c.Errors.Last(); e != nil {
		busCode = ierrors.Code(e.Err)
	}

	data := it.H{
		"Request":        beauty(req.Header, requestData),
		"RequestQuery":   req.URL.Query().Encode(),
		"RequestHeader":  c.Request.Header,
		"Response":       beauty(c.Writer.Header(), w.data),
		"ResponseHeader": w.Header(),
		"StatusCode":     w.Status(),
	}

	// 添加自定义日志数据到日志组件
	logData := c.GetLogData()
	for key, value := range logData {
		data[key] = value
	}

	kvs := []any{
		"Operate", "RequestLog",
		"URL", req.Method + "|" + req.URL.Path,
		"Data", data,
		"Time", it.H{"Req": t, "Res": now, "Duration": now.Sub(t).Milliseconds()},
	}
	if busCode == ierrors.GetInternalServerErrorCode() {
		ilogger.ErrorwCtx(c.Request.Context(), "Error", kvs...)
	} else {
		ilogger.InfowCtx(c.Request.Context(), "OK", kvs...)
	}
}

func beauty(header http.Header, data []byte) any {
	contentType := header.Get("Content-Type")
	if strings.Contains(contentType, "form-data") ||
		strings.Contains(contentType, "application/octet-stream") {
		return fmt.Sprintf("ignore file data(%d bytes)", len(data))
	}
	disposition := header.Get("Content-Disposition")
	if strings.Contains(disposition, "filename=") {
		return fmt.Sprintf("ignore file data(%d bytes)", len(data))
	}

	m := make(it.H)
	err := json.Unmarshal(data, &m)
	if err != nil {
		if len(data) > 1048576 {
			return fmt.Sprintf(">1M:len=(%d)", len(data))
		}
		return string(data)
	}
	return m
}
