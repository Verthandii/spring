package iredis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// Del 执行redis del操作
func (r *RedisCli) Del(ctx context.Context, keys ...string) *redis.IntCmd {
	if r.Client != nil {
		return r.Client.Del(ctx, keys...)
	}
	return r.Cluster.Del(ctx, keys...)
}

func (r *RedisCli) Get(ctx context.Context, key string) *redis.StringCmd {
	if r.Client != nil {
		return r.Client.Get(ctx, key)
	}
	return r.Cluster.Get(ctx, key)
}

func (r *RedisCli) Pipeline() redis.Pipeliner {
	if r.Client != nil {
		return r.Client.Pipeline()
	}
	return r.Cluster.Pipeline()
}

func (r *RedisCli) Incr(ctx context.Context, key string) *redis.IntCmd {
	if r.Client != nil {
		return r.Client.Incr(ctx, key)
	}
	return r.Cluster.Incr(ctx, key)
}

func (r *RedisCli) IncrBy(ctx context.Context, key string, value int64) *redis.IntCmd {
	if r.Client != nil {
		return r.Client.IncrBy(ctx, key, value)
	}
	return r.Cluster.IncrBy(ctx, key, value)
}

func (r *RedisCli) Expire(ctx context.Context, key string, ex time.Duration) *redis.BoolCmd {
	if r.Client != nil {
		return r.Client.Expire(ctx, key, ex)
	}
	return r.Cluster.Expire(ctx, key, ex)
}

func (r *RedisCli) Set(ctx context.Context, key string, value interface{}, ex time.Duration) *redis.StatusCmd {
	if r.Client != nil {
		return r.Client.Set(ctx, key, value, ex)
	}
	return r.Cluster.Set(ctx, key, value, ex)
}

func (r *RedisCli) SetEx(ctx context.Context, key string, value interface{}, ex time.Duration) *redis.StatusCmd {
	if r.Client != nil {
		return r.Client.SetEx(ctx, key, value, ex)
	}
	return r.Cluster.SetEx(ctx, key, value, ex)
}

func (r *RedisCli) SetNX(ctx context.Context, key string, value interface{}, ex time.Duration) *redis.BoolCmd {
	if r.Client != nil {
		return r.Client.SetNX(ctx, key, value, ex)
	}
	return r.Cluster.SetNX(ctx, key, value, ex)
}

func (r *RedisCli) Scan(ctx context.Context, cursor uint64, match string, count int64) *redis.ScanCmd {
	if r.Client != nil {
		return r.Client.Scan(ctx, cursor, match, count)
	}
	return r.Cluster.Scan(ctx, cursor, match, count)
}

func (r *RedisCli) Keys(ctx context.Context, pattern string) *redis.StringSliceCmd {
	if r.Client != nil {
		return r.Client.Keys(ctx, pattern)
	}
	return r.Cluster.Keys(ctx, pattern)
}

func (r *RedisCli) Exists(ctx context.Context, keys ...string) *redis.IntCmd {
	if r.Client != nil {
		return r.Client.Exists(ctx, keys...)
	}
	return r.Cluster.Exists(ctx, keys...)
}

func (r *RedisCli) AddHook(hook redis.Hook) {
	if r.Client != nil {
		r.Client.AddHook(hook)
	} else {
		r.Cluster.AddHook(hook)
	}
}

func (r *RedisCli) ZRange(ctx context.Context, key string, start, stop int64) *redis.StringSliceCmd {
	if r.Cluster != nil {
		return r.Cluster.ZRange(ctx, key, start, stop)
	}
	return r.Client.ZRange(ctx, key, start, stop)
}

func (r *RedisCli) ZAdd(ctx context.Context, key string, members ...redis.Z) *redis.IntCmd {
	if r.Cluster != nil {
		return r.Cluster.ZAdd(ctx, key, members...)
	}
	return r.Client.ZAdd(ctx, key, members...)
}

func (r *RedisCli) ZRemRangeByScore(ctx context.Context, key, min, max string) *redis.IntCmd {
	if r.Cluster != nil {
		return r.Cluster.ZRemRangeByScore(ctx, key, min, max)
	}
	return r.Client.ZRemRangeByScore(ctx, key, min, max)
}

func (r *RedisCli) Eval(ctx context.Context, script string, keys []string, args ...interface{}) *redis.Cmd {
	if r.Cluster != nil {
		return r.Cluster.Eval(ctx, script, keys, args...)
	}
	return r.Client.Eval(ctx, script, keys, args...)
}
