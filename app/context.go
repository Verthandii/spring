package app

import (
	"context"
	"time"

	"github.com/Verthandii/spring/db/iredis"
)

type ICacheAppContext interface {
	GetContext() context.Context
	GetCacheRedis() *iredis.RedisCli
	Load(rc *iredis.RedisCli, key string, v interface{}) (ok bool)
	Store(rc *iredis.RedisCli, key string, v interface{}, ex time.Duration)
	Lock(key string, ex time.Duration) bool
	UnLock(key string)
}
