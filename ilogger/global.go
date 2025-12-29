package ilogger

import (
	"context"
	"sync"
)

// globalLogger is designed as a global logger in current process.
var global = &loggerAppliance{}

// loggerAppliance is the proxy of `Logger` to
// make logger change will affect all sub-logger.
type loggerAppliance struct {
	lock sync.Mutex
	Logger
}

func (a *loggerAppliance) SetLogger(in Logger) {
	a.lock.Lock()
	defer a.lock.Unlock()
	a.Logger = in
}

// SetLogger should be called before any other log call.
func SetLogger(logger Logger) {
	global.SetLogger(logger)
}

// GetLogger returns global logger appliance as logger in current process.
func GetLogger() Logger {
	return global
}

func Debug(msg string, kvs ...any) {
	global.Debug(msg, kvs...)
}

func Info(msg string, kvs ...any) {
	global.Info(msg, kvs...)
}

func Warn(msg string, kvs ...any) {
	global.Warn(msg, kvs...)
}

func Error(msg string, kvs ...any) {
	global.Error(msg, kvs...)
}

func Fatal(msg string, kvs ...any) {
	global.Fatal(msg, kvs...)
}

func Debugf(template string, args ...any) {
	global.Debugf(template, args...)
}

func Infof(template string, args ...any) {
	global.Infof(template, args...)
}

func Warnf(template string, args ...any) {
	global.Warnf(template, args...)
}

func Errorf(template string, args ...any) {
	global.Errorf(template, args...)
}

func Fatalf(template string, args ...any) {
	global.Fatalf(template, args...)
}

func DebugwCtx(ctx context.Context, msg string, fields ...any) {
	global.DebugwCtx(ctx, msg, fields...)
}

func InfowCtx(ctx context.Context, msg string, fields ...any) {
	global.InfowCtx(ctx, msg, fields...)
}

func WarnwCtx(ctx context.Context, msg string, fields ...any) {
	global.WarnwCtx(ctx, msg, fields...)
}

func ErrorwCtx(ctx context.Context, msg string, fields ...any) {
	global.ErrorwCtx(ctx, msg, fields...)
}

func FatalwCtx(ctx context.Context, msg string, fields ...any) {
	global.FatalwCtx(ctx, msg, fields...)
}
