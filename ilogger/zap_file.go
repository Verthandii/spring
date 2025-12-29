package ilogger

import (
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/Verthandii/spring/env"
)

// FileZapLogger 输出日志到指定的文件
type FileZapLogger struct {
	logger       *zap.SugaredLogger
	builtinPairs []interface{}
}

func (z *FileZapLogger) log(level zapcore.Level, msg string, kvs ...interface{}) {
	switch level {
	case zap.DebugLevel:
		z.logger.Debugw(msg, kvs...)
	case zap.InfoLevel:
		z.logger.Infow(msg, kvs...)
	case zap.WarnLevel:
		z.logger.Warnw(msg, kvs...)
	case zap.ErrorLevel:
		z.logger.Errorw(msg, kvs...)
	case zap.FatalLevel:
		z.logger.Fatalw(msg, kvs...)
	default:
		z.logger.Warnw(msg, kvs...)
	}
}

func (z *FileZapLogger) logf(level zapcore.Level, msg string, kvs ...any) {
	switch level {
	case zap.DebugLevel:
		z.logger.Debugf(msg, kvs...)
	case zap.InfoLevel:
		z.logger.Infof(msg, kvs...)
	case zap.WarnLevel:
		z.logger.Warnf(msg, kvs...)
	case zap.ErrorLevel:
		z.logger.Errorf(msg, kvs...)
	case zap.FatalLevel:
		z.logger.Fatalf(msg, kvs...)
	default:
		z.logger.Warnf(msg, kvs...)
	}
}

func (z *FileZapLogger) logwCtx(ctx context.Context, level zapcore.Level, msg string, keyvals ...interface{}) {
	kvs := make([]interface{}, 0, len(_builtinPairs)+len(z.builtinPairs)+len(keyvals))
	kvs = append(kvs, z.builtinPairs...)
	kvs = append(kvs, _builtinPairs...)
	kvs = append(kvs, keyvals...)
	if len(kvs) > 0 {
		bindValues(ctx, kvs)
	}

	switch level {
	case zap.DebugLevel:
		z.logger.Debugw(msg, kvs...)
	case zap.InfoLevel:
		z.logger.Infow(msg, kvs...)
	case zap.WarnLevel:
		z.logger.Warnw(msg, kvs...)
	case zap.ErrorLevel:
		z.logger.Errorw(msg, kvs...)
	case zap.FatalLevel:
		z.logger.Fatalw(msg, kvs...)
	default:
		z.logger.Warnw(msg, kvs...)
	}
}

func (z *FileZapLogger) Debug(msg string, kvs ...any) {
	z.log(zap.DebugLevel, msg, kvs...)
}

func (z *FileZapLogger) Info(msg string, kvs ...any) {
	z.log(zap.InfoLevel, msg, kvs...)
}

func (z *FileZapLogger) Warn(msg string, kvs ...any) {
	z.log(zap.WarnLevel, msg, kvs...)
}

func (z *FileZapLogger) Error(msg string, kvs ...any) {
	z.log(zap.ErrorLevel, msg, kvs...)
}

func (z *FileZapLogger) Fatal(msg string, kvs ...any) {
	z.log(zap.FatalLevel, msg, kvs...)
}

func (z *FileZapLogger) Debugf(template string, args ...any) {
	z.logf(zap.DebugLevel, template, args...)
}

func (z *FileZapLogger) Infof(template string, args ...any) {
	z.logf(zap.InfoLevel, template, args...)
}

func (z *FileZapLogger) Warnf(template string, args ...any) {
	z.logf(zap.WarnLevel, template, args...)
}

func (z *FileZapLogger) Errorf(template string, args ...any) {
	z.logf(zap.ErrorLevel, template, args...)
}

func (z *FileZapLogger) Fatalf(template string, args ...any) {
	z.logf(zap.FatalLevel, template, args...)
}

func (z *FileZapLogger) DebugwCtx(ctx context.Context, msg string, keyvals ...interface{}) {
	if z == nil {
		return
	}

	z.logwCtx(ctx, zapcore.DebugLevel, msg, keyvals...)
}

func (z *FileZapLogger) InfowCtx(ctx context.Context, msg string, keyvals ...interface{}) {
	if z == nil {
		return
	}

	z.logwCtx(ctx, zapcore.InfoLevel, msg, keyvals...)
}

func (z *FileZapLogger) WarnwCtx(ctx context.Context, msg string, keyvals ...interface{}) {
	if z == nil {
		return
	}

	z.logwCtx(ctx, zapcore.WarnLevel, msg, keyvals...)
}

func (z *FileZapLogger) ErrorwCtx(ctx context.Context, msg string, keyvals ...interface{}) {
	if z == nil {
		return
	}

	z.logwCtx(ctx, zapcore.ErrorLevel, msg, keyvals...)
}

func (z *FileZapLogger) FatalwCtx(ctx context.Context, msg string, keyvals ...interface{}) {
	if z == nil {
		return
	}

	z.logwCtx(ctx, zapcore.FatalLevel, msg, keyvals...)
}

func (z *FileZapLogger) addValuer(key string, v Valuer) {
	z.builtinPairs = append(z.builtinPairs, key, v)
}

// NewFileZapLogger .
func NewFileZapLogger(filename string, opts ...Option) Logger {
	fields := []zap.Field{
		zap.String(env.AppEnv, env.GetEnv()),
		zap.String(env.AppVersion, env.GetVersion()),
		zap.String(env.AppServer, env.GetServer()),
		zap.String(env.AppRegion, env.GetLRegion()),
	}
	options := []zap.Option{
		zap.AddCaller(),
		zap.Development(),
		zap.AddStacktrace(zap.ErrorLevel),
	}
	if len(fields) > 0 {
		options = append(options, zap.Fields(fields...))
	}

	writer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   filename, // 日志文件路径
		MaxSize:    10,       // 日志最大保存 M
		MaxAge:     3,        // 日志保留天数
		Compress:   false,    // 是否压缩
		MaxBackups: 10,
	})

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	core := zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), writer, zapLogLevel(env.LogLevel()))

	logger := &FileZapLogger{logger: zap.New(core, options...).Sugar()}
	for _, opt := range opts {
		opt.apply(logger)
	}

	return logger
}
