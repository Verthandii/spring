package iorm

import (
	"time"

	"github.com/Verthandii/spring/db/icache"
	"github.com/Verthandii/spring/db/iredis"
	"github.com/Verthandii/spring/ilogger"
	"github.com/Verthandii/spring/it"
	"github.com/Verthandii/spring/utils"
)

type Cache struct {
	key       string
	field     string
	redis     *iredis.RedisCli
	mcache    *icache.MemCache
	StoreType iredis.StoreType
	expire    time.Duration
}

// HCache 使用 redis hash
func (s *Session) HCache(key, field string, mcache ...bool) *Session {
	if len(mcache) > 0 && mcache[0] {
		// 用 mcache
		s.cache = &Cache{key: key, field: field, redis: s.o.redis, mcache: s.o.mcache}
	} else {
		s.cache = &Cache{key: key, field: field, redis: s.o.redis}
	}
	s.cache.StoreType = iredis.StoreHash
	return s
}

// Cache 使用 redis string
func (s *Session) Cache(key string, mcache ...bool) *Session {
	if len(mcache) > 0 && mcache[0] {
		// 用 mcache
		s.cache = &Cache{key: key, redis: s.o.redis, mcache: s.o.mcache}
	} else {
		s.cache = &Cache{key: key, redis: s.o.redis}
	}
	s.cache.StoreType = iredis.StoreString
	return s
}

func (s *Session) CacheEx(ex time.Duration) *Session {
	s.cache.expire = ex
	return s
}

// LoadCache 加载缓存
func (s *Session) LoadCache(v interface{}) bool {
	if s.cache == nil {
		return false
	}
	cache := s.cache
	var ok bool
	if cache.StoreType == iredis.StoreString {
		ok = s.Load(cache.redis, cache.key, v)
	} else {
		ok = s.HLoad(cache.key, cache.field, v)
	}
	if ok {
		s.RowsAffected = 1
		return true
	}
	return false
}

// StoreCache 存储缓存
func (s *Session) StoreCache(v interface{}) {
	if s.Error != nil || s.cache == nil || s.RowsAffected == 0 {
		return
	}
	cache := s.cache
	if cache != nil {
		if cache.StoreType == iredis.StoreString {
			s.Store(cache.redis, cache.key, v, cache.expire)
		} else {
			s.HStore(cache.key, cache.field, v)
		}
	}
}

// Load 加载数据并反序列化
func (s *Session) Load(ri *iredis.RedisCli, key string, v interface{}) (ok bool) {
	if ri == nil {
		return ok
	}
	ok, err := ri.Load(s.ctx, key, v)
	if err != nil {
		ilogger.ErrorwCtx(s.ctx, err.Error(), "Operate", "Load", "Key", key)
		return false
	}
	return ok
}

// Store 序列化数据并储存数据
func (s *Session) Store(ri *iredis.RedisCli, key string, v interface{}, ex time.Duration) {
	if ri == nil {
		return
	}
	ri.Store(s.ctx, key, v, ex)
}

// HLoad 加载通用缓存 Hash
func (s *Session) HLoad(key, field string, v interface{}) bool {
	mcache := s.cache.mcache != nil
	if mcache {
		if m, ex := s.cache.mcache.Get(key + field); ex {
			if err := utils.ModifyIT(v, m); err != nil {
				ilogger.ErrorwCtx(s.ctx, err.Error(), "Operate", "ModifyIT", "Key", key+field)
				return false
			}
			return true
		}
	}
	if v == nil {
		v = &it.H{}
	}
	ok, err := s.cache.redis.HLoad(s.ctx, key, field, v)
	if err != nil {
		ilogger.ErrorwCtx(s.ctx, err.Error(), "Operate", "Load", "Key", key)
	}
	if ok && mcache {
		s.cache.mcache.Set(key+field, v)
		return ok
	}
	return ok
}

// HStore 存储通用缓存 Hash
func (s *Session) HStore(key, field string, v interface{}) {
	if m := s.cache.mcache; m != nil {
		m.Set(key+field, v)
	}
	if r := s.cache.redis; r != nil {
		r.HStore(s.ctx, key, field, v, 0)
	}
}

// CHLoad 自定义 Hash 加载缓存
func (s *Session) CHLoad(key, field string, v interface{}, mcache ...bool) bool {
	m := len(mcache) > 0 && mcache[0]
	s.HCache(key, field, m)
	defer s.teardown()
	return s.HLoad(key, field, v)
}

// CHStore 自定义 Hash 存储缓存
func (s *Session) CHStore(key, field string, v interface{}, mcache ...bool) {
	m := len(mcache) > 0 && mcache[0]
	s.HCache(key, field, m)
	defer s.teardown()
	s.HStore(key, field, v)
}
