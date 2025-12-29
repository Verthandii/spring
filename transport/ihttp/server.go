package ihttp

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"

	"github.com/Verthandii/spring/encoding"
	"github.com/Verthandii/spring/env"
	"github.com/Verthandii/spring/ierrors"
	"github.com/Verthandii/spring/ilogger"
	"github.com/Verthandii/spring/third/otel"
	_ "github.com/Verthandii/spring/third/otel"
)

func upperRespPool() sync.Pool {
	return sync.Pool{
		New: func() any {
			return &upperResponse{}
		},
	}
}

func lowerRespPool() sync.Pool {
	return sync.Pool{
		New: func() any {
			return &lowerResponse{}
		},
	}
}

type Option func(s *Server)

func Addr(addr string) Option {
	return func(s *Server) {
		s.addr = addr
	}
}

func Timeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.timeout = timeout
	}
}

func WithEncodeResponseFunc(enc SuccessEncodeFunc) Option {
	return func(s *Server) {
		s.successEc = enc
	}
}

func WithEncodeErrorFunc(ene ErrorEncodeFunc) Option {
	return func(s *Server) {
		s.errorEc = ene
	}
}

func WithUpperResp() Option {
	return func(s *Server) {
		s.respPool = upperRespPool()
	}
}

func WithLowerResp() Option {
	return func(s *Server) {
		s.respPool = lowerRespPool()
	}
}

func WithMetrics(enable bool) Option {
	return func(s *Server) {
		s.enableMetrics = enable
	}
}

func WithSiteRouter(enable bool) Option {
	return func(s *Server) {
		s.enableSiteRouter = enable
	}
}

func WithSiteRouterPrefix(prefix string) Option {
	return func(s *Server) {
		s.siteRouterPrefix = prefix
	}
}

type Server struct {
	RouterGroup
	Engine    *gin.Engine
	srv       *http.Server
	addr      string
	timeout   time.Duration
	successEc SuccessEncodeFunc // c.Success encoder
	errorEc   ErrorEncodeFunc   // c.Failure encoder
	respPool  sync.Pool         // 默认大写

	enableMetrics    bool   // 是否开启 metrics，默认开启
	enableSiteRouter bool   // 是否开启 site 路由，默认开启
	siteRouterPrefix string // site 路由前缀，默认空
}

type RouterGroup struct {
	*gin.RouterGroup
	srv *Server
}

func setTID(c *gin.Context) {
	h := c.Writer.Header()
	h["TID"] = []string{otel.ID(c)}
}

func NewServer(opts ...Option) *Server {
	if env.IsStagingOrProd() {
		gin.DefaultWriter = io.Discard
	}
	g := gin.New()
	g.ContextWithFallback = true
	g.Use(otelgin.Middleware(env.GetServer()), setTID)
	srv := &Server{
		RouterGroup:      RouterGroup{RouterGroup: &g.RouterGroup},
		Engine:           g,
		srv:              nil,
		addr:             ":8000",
		timeout:          1 * time.Second,
		respPool:         upperRespPool(),
		enableMetrics:    true,
		enableSiteRouter: true,
		siteRouterPrefix: "",
	}
	srv.RouterGroup.srv = srv
	srv.errorEc = srv.DefaultErrorEncoder
	srv.successEc = srv.DefaultSuccessEncoder

	for _, opt := range opts {
		opt(srv)
	}

	if srv.enableMetrics {
		srv.Use(serverMetrics)
	}

	srv.srv = &http.Server{
		Addr:         srv.addr,
		Handler:      http.AllowQuerySemicolons(srv.Engine),
		WriteTimeout: srv.timeout,
	}

	return srv
}

func (s *Server) Run() error {
	if s.enableSiteRouter {
		s.initSiteRouter()
	}
	ilogger.Info("【HTTP】启动成功", "addr", s.addr)
	return s.srv.ListenAndServe()
}

func (s *Server) Close(ctx context.Context) {
	ilogger.Info("【HTTP】已关闭")
	_ = s.srv.Shutdown(ctx)
}

func (s *Server) DefaultSuccessEncoder(w http.ResponseWriter, req *http.Request, data any) {
	if data == nil {
		data = _empty
	}

	w.WriteHeader(http.StatusOK)
	h := w.Header()
	h["Content-Type"] = []string{"application/json; charset=utf-8"}

	resp := s.respPool.Get().(response)
	resp.setCode(http.StatusOK)
	resp.setMsg("OK")
	resp.setData(data)

	codec := encoding.GetCodec("json")
	respBytes, err := codec.Marshal(resp)
	if err != nil {
		_, _ = w.Write([]byte(err.Error()))
		resp.reset()
		s.respPool.Put(resp)
		return
	}

	_, _ = w.Write(respBytes)
	resp.reset()
	s.respPool.Put(resp)
}

func (s *Server) DefaultErrorEncoder(w http.ResponseWriter, req *http.Request, err error, data any, statusCode int) {
	if data == nil {
		data = _empty
	}

	w.WriteHeader(statusCode)
	h := w.Header()
	h["Content-Type"] = []string{"application/json; charset=utf-8"}

	code, msg := ierrors.Info(err)

	resp := s.respPool.Get().(response)
	resp.setCode(code)
	resp.setMsg(msg)
	resp.setData(data)

	codec := encoding.GetCodec("json")
	respBytes, err := codec.Marshal(resp)
	if err != nil {
		_, _ = w.Write([]byte(err.Error()))
		resp.reset()
		s.respPool.Put(resp)
		return
	}

	_, _ = w.Write(respBytes)
	resp.reset()
	s.respPool.Put(resp)
}

func (s *Server) Use(handlers ...Handler) {
	s.Engine.Use(s.convertHandler(handlers)...)
}

// convertHandler 把ihttp的Handler转成成gin的Handler
func (g *RouterGroup) convertHandler(handlers []Handler) []gin.HandlerFunc {
	var ghs []gin.HandlerFunc
	for _, m := range handlers {
		ghs = append(ghs, g.srv.H(m))
	}
	return ghs
}

// Use attaches a global middleware to the router. i.e. the middleware attached through Use() will be
// included in the handlers chain for every single request. Even 404, 405, static files...
// For example, this is the right place for a logger or error management middleware.
func (g *RouterGroup) Use(handlers ...Handler) {
	g.RouterGroup.Use(g.srv.convertHandler(handlers)...)
}

// Handle registers a new request handle and middleware with the given path and method.
// The last handler should be the real handler, the other ones should be middleware that can and should be shared among different routes.
// See the example code in GitHub.
//
// For GET, POST, PUT, PATCH and DELETE requests the respective shortcut
// functions can be used.
//
// This function is intended for bulk loading and to allow the usage of less
// frequently used, non-standardized or custom methods (e.g. for internal
// communication with a proxy).
func (g *RouterGroup) Handle(httpMethod, relativePath string, handlers ...Handler) {
	g.RouterGroup.Handle(httpMethod, relativePath, g.srv.convertHandler(handlers)...)
}

// POST is a shortcut for router.Handle("POST", path, handlers).
func (g *RouterGroup) POST(relativePath string, handlers ...Handler) {
	g.RouterGroup.Handle(http.MethodPost, relativePath, g.srv.convertHandler(handlers)...)
}

// GET is a shortcut for router.Handle("GET", path, handlers).
func (g *RouterGroup) GET(relativePath string, handlers ...Handler) {
	g.RouterGroup.Handle(http.MethodGet, relativePath, g.srv.convertHandler(handlers)...)
}

// DELETE is a shortcut for router.Handle("DELETE", path, handlers).
func (g *RouterGroup) DELETE(relativePath string, handlers ...Handler) {
	g.RouterGroup.Handle(http.MethodDelete, relativePath, g.srv.convertHandler(handlers)...)
}

// PATCH is a shortcut for router.Handle("PATCH", path, handlers).
func (g *RouterGroup) PATCH(relativePath string, handlers ...Handler) {
	g.RouterGroup.Handle(http.MethodPatch, relativePath, g.srv.convertHandler(handlers)...)
}

// PUT is a shortcut for router.Handle("PUT", path, handlers).
func (g *RouterGroup) PUT(relativePath string, handlers ...Handler) {
	g.RouterGroup.Handle(http.MethodPut, relativePath, g.srv.convertHandler(handlers)...)
}

// OPTIONS is a shortcut for router.Handle("OPTIONS", path, handlers).
func (g *RouterGroup) OPTIONS(relativePath string, handlers ...Handler) {
	g.RouterGroup.Handle(http.MethodOptions, relativePath, g.srv.convertHandler(handlers)...)
}

// HEAD is a shortcut for router.Handle("HEAD", path, handlers).
func (g *RouterGroup) HEAD(relativePath string, handlers ...Handler) {
	g.RouterGroup.Handle(http.MethodHead, relativePath, g.srv.convertHandler(handlers)...)
}

// Any registers a route that matches all the HTTP methods.
// GET, POST, PUT, PATCH, HEAD, OPTIONS, DELETE, CONNECT, TRACE.
func (g *RouterGroup) Any(relativePath string, handlers ...Handler) {
	g.RouterGroup.Any(relativePath, g.srv.convertHandler(handlers)...)
}

// Match registers a route that matches the specified methods that you declared.
func (g *RouterGroup) Match(methods []string, relativePath string, handlers ...Handler) {
	g.RouterGroup.Match(methods, relativePath, g.srv.convertHandler(handlers)...)
}

// Group creates a new router group. You should add all the routes that have common middlewares or the same path prefix.
// For example, all the routes that use a common middleware for authorization could be grouped.
func (g *RouterGroup) Group(relativePath string, handlers ...Handler) *RouterGroup {
	return &RouterGroup{
		RouterGroup: g.RouterGroup.Group(relativePath, g.srv.convertHandler(handlers)...),
		srv:         g.srv,
	}
}

// initSiteRouter 初始化站点路由
func (g *RouterGroup) initSiteRouter() {
	g.GET("", SiteIndex)
	g.GET("ping", SitePing)
	t := g.Group(fmt.Sprintf("%s/tools", g.srv.siteRouterPrefix))
	{
		t.GET("test-log", SiteTestLog)
		t.GET("env-info", SiteEnvInfo)
		t.PUT("timezone", SiteSetTimeZone)
	}
}
