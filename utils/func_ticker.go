package utils

import (
	"context"
	"time"
)

type FuncTicker struct {
	immediately bool
	duration    time.Duration
	fn          func()
}

type Option func(f *FuncTicker)

func WithImmediatelyOnce(immediatelyOnce bool) Option {
	return func(f *FuncTicker) {
		f.immediately = immediatelyOnce
	}
}

func NewFuncTicker(fn func(), duration time.Duration, opts ...Option) *FuncTicker {
	f := &FuncTicker{
		duration: duration,
		fn:       fn,
	}
	for _, opt := range opts {
		opt(f)
	}
	return f
}

func (f *FuncTicker) Exec(ctx context.Context) {
	if f.immediately {
		f.fn()
	}

	ticker := time.NewTicker(f.duration)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			f.fn()
		}
	}
}
