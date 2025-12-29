package ilogger

import (
	"context"
	"errors"
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/Verthandii/spring/env"
)

// ZapLogger Debug Info Warn 输出到【标准输出】
// Error 输出到【标准错误】
// 由于在请求 Redis 或者 MySQL 时经常出现客户端断开连接
// 导致产生 context.Canceled 错误，但没有输出此类错误的必要，Error 在发现此类错误时会使用 Warn 代替
type ZapLogger struct {
	stdout       *zap.SugaredLogger
	stderr       *zap.SugaredLogger
	builtinPairs []interface{}
}

func (z *ZapLogger) log(level zapcore.Level, msg string, kvs ...interface{}) {
	switch level {
	case zap.DebugLevel:
		z.stdout.Debugw(msg, kvs...)
	case zap.InfoLevel:
		z.stdout.Infow(msg, kvs...)
	case zap.WarnLevel:
		z.stdout.Warnw(msg, kvs...)
	case zap.ErrorLevel:
		z.stderr.Errorw(msg, kvs...)
	case zap.FatalLevel:
		z.stderr.Fatalw(msg, kvs...)
	default:
		z.stdout.Warnw(msg, kvs...)
	}
}

func (z *ZapLogger) logf(level zapcore.Level, msg string, kvs ...any) {
	switch level {
	case zap.DebugLevel:
		z.stdout.Debugf(msg, kvs...)
	case zap.InfoLevel:
		z.stdout.Infof(msg, kvs...)
	case zap.WarnLevel:
		z.stdout.Warnf(msg, kvs...)
	case zap.ErrorLevel:
		z.stderr.Errorf(msg, kvs...)
	case zap.FatalLevel:
		z.stderr.Fatalf(msg, kvs...)
	default:
		z.stdout.Warnf(msg, kvs...)
	}
}

func (z *ZapLogger) logwCtx(ctx context.Context, level zapcore.Level, msg string, keyvals ...interface{}) {
	kvs := make([]interface{}, 0, len(_builtinPairs)+len(z.builtinPairs)+len(keyvals))
	kvs = append(kvs, z.builtinPairs...)
	kvs = append(kvs, _builtinPairs...)
	kvs = append(kvs, keyvals...)
	if len(kvs) > 0 {
		bindValues(ctx, kvs)
	}

	switch level {
	case zap.DebugLevel:
		z.stdout.Debugw(msg, kvs...)
	case zap.InfoLevel:
		z.stdout.Infow(msg, kvs...)
	case zap.WarnLevel:
		z.stdout.Warnw(msg, kvs...)
	case zap.ErrorLevel:
		z.stderr.Errorw(msg, kvs...)
	case zap.FatalLevel:
		z.stderr.Fatalw(msg, kvs...)
	default:
		z.stdout.Warnw(msg, kvs...)
	}
}

func (z *ZapLogger) Debug(msg string, kvs ...any) {
	z.log(zap.DebugLevel, msg, kvs...)
}

func (z *ZapLogger) Info(msg string, kvs ...any) {
	z.log(zap.InfoLevel, msg, kvs...)
}

func (z *ZapLogger) Warn(msg string, kvs ...any) {
	z.log(zap.WarnLevel, msg, kvs...)
}

func (z *ZapLogger) Error(msg string, kvs ...any) {
	for _, kv := range kvs {
		if err, ok := kv.(error); ok && errors.Is(err, context.Canceled) {
			z.Warn(msg, kvs...)
			return
		}
	}
	z.log(zap.ErrorLevel, msg, kvs...)
}

func (z *ZapLogger) Fatal(msg string, kvs ...any) {
	z.log(zap.FatalLevel, msg, kvs...)
}

func (z *ZapLogger) Debugf(template string, args ...any) {
	z.logf(zap.DebugLevel, template, args...)
}

func (z *ZapLogger) Infof(template string, args ...any) {
	z.logf(zap.InfoLevel, template, args...)
}

func (z *ZapLogger) Warnf(template string, args ...any) {
	z.logf(zap.WarnLevel, template, args...)
}

func (z *ZapLogger) Errorf(template string, args ...any) {
	for _, kv := range args {
		if err, ok := kv.(error); ok && errors.Is(err, context.Canceled) {
			z.Warnf(template, args...)
			return
		}
	}
	z.logf(zap.ErrorLevel, template, args...)
}

func (z *ZapLogger) Fatalf(template string, args ...any) {
	z.logf(zap.FatalLevel, template, args...)
}

func (z *ZapLogger) DebugwCtx(ctx context.Context, msg string, keyvals ...interface{}) {
	if z == nil {
		return
	}

	z.logwCtx(ctx, zapcore.DebugLevel, msg, keyvals...)
}

func (z *ZapLogger) InfowCtx(ctx context.Context, msg string, keyvals ...interface{}) {
	if z == nil {
		return
	}

	z.logwCtx(ctx, zapcore.InfoLevel, msg, keyvals...)
}

func (z *ZapLogger) WarnwCtx(ctx context.Context, msg string, keyvals ...interface{}) {
	if z == nil {
		return
	}

	z.logwCtx(ctx, zapcore.WarnLevel, msg, keyvals...)
}

func (z *ZapLogger) ErrorwCtx(ctx context.Context, msg string, keyvals ...interface{}) {
	if z == nil {
		return
	}

	z.logwCtx(ctx, zapcore.ErrorLevel, msg, keyvals...)
}

func (z *ZapLogger) FatalwCtx(ctx context.Context, msg string, keyvals ...interface{}) {
	if z == nil {
		return
	}

	z.logwCtx(ctx, zapcore.FatalLevel, msg, keyvals...)
}

func (z *ZapLogger) addValuer(key string, v Valuer) {
	z.builtinPairs = append(z.builtinPairs, key, v)
}

func NewZapLogger(opts ...Option) Logger {
	fields := []zap.Field{
		zap.String(env.AppEnv, env.GetEnv()),
		zap.String(env.AppVersion, env.GetVersion()),
		zap.String(env.AppServer, env.GetServer()),
		zap.String(env.AppRegion, env.GetLRegion()),
	}
	options := []zap.Option{
		zap.AddCaller(),
		zap.AddCallerSkip(3),
		zap.Development(),
		zap.AddStacktrace(zap.ErrorLevel),
	}
	if len(fields) > 0 {
		options = append(options, zap.Fields(fields...))
	}

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	stdoutCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout)),
		zapLogLevel(env.LogLevel()),
	)

	stderrCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stderr)),
		zap.ErrorLevel,
	)

	logger := &ZapLogger{
		stdout: zap.New(stdoutCore, options...).Sugar(),
		stderr: zap.New(stderrCore, options...).Sugar(),
	}
	for _, opt := range opts {
		opt.apply(logger)
	}

	return logger
}

func zapLogLevel(level string) zapcore.Level {
	switch strings.ToLower(level) {
	case "debug":
		return zap.DebugLevel
	case "info":
		return zap.InfoLevel
	case "warn":
		return zap.WarnLevel
	case "error":
		return zap.ErrorLevel
	default:
		return zap.InfoLevel
	}
}
