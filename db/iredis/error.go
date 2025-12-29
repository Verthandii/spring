package iredis

import (
	"context"
	"errors"

	"github.com/redis/go-redis/v9"

	"github.com/Verthandii/spring/db"
)

type errorHook struct{}

func (e *errorHook) DialHook(next redis.DialHook) redis.DialHook {
	return next
}

func (e *errorHook) ProcessHook(next redis.ProcessHook) redis.ProcessHook {
	return func(ctx context.Context, cmd redis.Cmder) error {
		err := next(ctx, cmd)
		if errors.Is(cmd.Err(), redis.Nil) {
			cmd.SetErr(errors.Join(redis.Nil, db.NotFound))
		}
		if errors.Is(err, redis.Nil) {
			return errors.Join(redis.Nil, db.NotFound)
		}
		return nil
	}
}

func (e *errorHook) ProcessPipelineHook(next redis.ProcessPipelineHook) redis.ProcessPipelineHook {
	return next
}
