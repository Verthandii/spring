package ilogger

import (
	"context"

	"go.uber.org/zap"
)

var (
	defaultStdLogger = NewZapLogger()
)

func init() {
	SetLogger(defaultStdLogger)
}

type Logger interface {
	Debug(msg string, kvs ...any)
	Debugf(msg string, kvs ...any)
	Info(msg string, kvs ...any)
	Infof(msg string, kvs ...any)
	Warn(msg string, kvs ...any)
	Warnf(msg string, kvs ...any)
	Error(msg string, kvs ...any)
	Errorf(msg string, kvs ...any)
	Fatal(msg string, kvs ...any)
	Fatalf(msg string, kvs ...any)
	DebugwCtx(ctx context.Context, msg string, kvs ...any)
	InfowCtx(ctx context.Context, msg string, kvs ...any)
	WarnwCtx(ctx context.Context, msg string, kvs ...any)
	ErrorwCtx(ctx context.Context, msg string, kvs ...any)
	FatalwCtx(ctx context.Context, msg string, kvs ...any)

	addValuer(key string, v Valuer)
}

var Any = zap.Any

// Valuer is returns a log value.
type Valuer func(ctx context.Context) interface{}

func bindValues(ctx context.Context, keyvals []interface{}) {
	for i := 0; i < len(keyvals); i++ {
		if v, ok := keyvals[i].(Valuer); ok {
			keyvals[i] = v(ctx)
		}
	}
}
