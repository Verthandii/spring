package spring

import (
	"errors"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"runtime"

	"go.uber.org/automaxprocs/maxprocs"

	_ "github.com/Verthandii/spring/encoding/form"
	_ "github.com/Verthandii/spring/encoding/json"
	_ "github.com/Verthandii/spring/encoding/xml"
	_ "github.com/Verthandii/spring/encoding/yaml"
	"github.com/Verthandii/spring/env"
	"github.com/Verthandii/spring/ilogger"
	"github.com/Verthandii/spring/metrics"
	_ "github.com/Verthandii/spring/third/otel"
	"github.com/Verthandii/spring/utils/signal"
)

func init() {
	if env.IsLocal() {
		fmt.Println(" __   _____  ____ _____  _    ____    ____  _   _ ____  ____  _        _  _____ \n \\ \\ / / _ \\/ ___|_   _|/ \\  |  _ \\  |  _ \\| | | | __ )|  _ \\| |      / \\|_   _|\n  \\ V / | | \\___ \\ | | / _ \\ | |_) | | |_) | | | |  _ \\| |_) | |     / _ \\ | |  \n   | || |_| |___) || |/ ___ \\|  _ <  |  __/| |_| | |_) |  __/| |___ / ___ \\| |  \n   |_| \\___/|____/ |_/_/   \\_\\_| \\_\\ |_|    \\___/|____/|_|   |_____/_/   \\_\\_|  \n                                                          v2.0.0 powered by cz&wz")
	}
	initMaxProcs()
}

func initMaxProcs() {
	_, _ = maxprocs.Set(maxprocs.Logger(func(string, ...any) {}))
	ilogger.Info("【初始化】协程并发量", "runtime.GOMAXPROCS(-1)", runtime.GOMAXPROCS(-1))
}

type Server interface {
	signal.Closer
	Run() error
}

type Option func(app *App)

func WithMetrics(enable bool) Option {
	return func(app *App) {
		app.enableMetrics = enable
	}
}

func WithPprof(pprof bool) Option {
	return func(app *App) {
		app.enablePprof = pprof
	}
}

func WithPProfAddr(addr string) Option {
	return func(app *App) {
		app.pprofAddr = addr
	}
}

type App struct {
	svr []Server

	enableMetrics bool
	// enablePprof 默认从环境变量【APP_DEBUG】读取
	enablePprof bool
	// pprofAddr 默认值为【:6789】
	pprofAddr string
}

func New(svr []Server, opts ...Option) (*App, error) {
	if len(svr) == 0 {
		return nil, errors.New("svr is nil")
	}

	app := &App{
		svr:           svr,
		enableMetrics: true,
		enablePprof:   env.EnableDebug(),
		pprofAddr:     ":6789",
	}

	for _, opt := range opts {
		opt(app)
	}

	if app.enableMetrics {
		metrics.Init(metrics.NewDefaultConfig())
	}

	return app, nil
}

func (app *App) Run() {
	if app.enablePprof {
		go func() {
			if err := http.ListenAndServe(app.pprofAddr, nil); err != nil {
				ilogger.Error("Pprof 异常", "err", err)
			}
		}()
	}
	for _, svr := range app.svr {
		go func(svr Server) {
			if err := svr.Run(); err != nil {
				ilogger.Warn("【App】启动过程出现异常", "err", err)
			}
		}(svr)
	}
	closers := make([]signal.Closer, 0)
	for _, svr := range app.svr {
		closers = append(closers, svr)
	}
	signal.ListenSignalAndShutdown(closers...)
}
