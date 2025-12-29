package iredis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/Verthandii/spring/db"
	"github.com/Verthandii/spring/ilogger"
	"github.com/Verthandii/spring/utils"
)

type StoreType int64

const (
	StoreString = iota
	StoreHash
)

// RedisCli Redis对象
type RedisCli struct {
	*redis.Client                      // 集群代理或单节点Redis
	Cluster       *redis.ClusterClient // 原生集群版Redis
}

func New(conf Config) (*RedisCli, error) {
	addr := fmt.Sprintf("%s:%d", conf.Host, conf.Port)
	ilogger.Infof("【Redis】开始连接 %s", addr)
	opt := &redis.Options{
		Addr:            addr,
		Password:        conf.Password,
		DB:              conf.DB,
		DialTimeout:     time.Second * time.Duration(conf.DialTimeout),
		ReadTimeout:     time.Second * time.Duration(conf.ReadTimeout),
		WriteTimeout:    time.Second * time.Duration(conf.WriteTimeout),
		ConnMaxIdleTime: time.Second * time.Duration(conf.IdleTimeout),
		PoolTimeout:     time.Second * time.Duration(conf.PoolTimeout),
		PoolSize:        conf.PoolSize,
		MinIdleConns:    conf.MinIdleConns,
	}
	client := redis.NewClient(opt)
	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}
	rc := &RedisCli{Client: client}
	rc.addHookMetrics(conf, addr)
	ilogger.Infof("【Redis】连接成功 %s", addr)
	return rc, nil
}

func NewCluster(conf Config) (*RedisCli, error) {
	addr := fmt.Sprintf("%s:%d", conf.Host, conf.Port)
	ilogger.Infof("【RedisCluster】开始连接 %s", addr)
	opt := &redis.ClusterOptions{
		Addrs:           []string{addr},
		Password:        conf.Password,
		DialTimeout:     time.Second * time.Duration(conf.DialTimeout),
		ReadTimeout:     time.Second * time.Duration(conf.ReadTimeout),
		WriteTimeout:    time.Second * time.Duration(conf.WriteTimeout),
		ConnMaxIdleTime: time.Second * time.Duration(conf.IdleTimeout),
		PoolTimeout:     time.Second * time.Duration(conf.PoolTimeout),
		PoolSize:        conf.PoolSize,
		MinIdleConns:    conf.MinIdleConns,
	}
	client := redis.NewClusterClient(opt)
	// 添加连接重试逻辑
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 重试3次
	var lastErr error
	for i := 0; i < 3; i++ {
		if err := client.Ping(ctx).Err(); err != nil {
			lastErr = err
			time.Sleep(time.Second)
			continue
		}
		// 连接成功
		rc := &RedisCli{Cluster: client}
		rc.addHookMetrics(conf, addr)
		ilogger.Infof("【RedisCluster】连接成功 %s", addr)
		return rc, nil
	}

	return nil, fmt.Errorf("连接Redis集群失败: %v", lastErr)
}

// addHookMetrics 添加hook
func (r *RedisCli) addHookMetrics(c Config, addr string) {
	// error hook
	r.AddHook(&errorHook{})
	// 默认开启metrics，可以配置关闭
	if !c.DisableMetrics {
		r.AddHook(newDurationHook(addr, c.Host, c.Port, strconv.Itoa(c.DB)))
		// metric
		_statCollector.registerClient(&statClient{
			addr:     addr,
			dbName:   "default",
			poolSize: c.PoolSize,
			inner:    r,
		})
	}
}

// HLoad 加载数据并反序列化
// 返回值第一个参数表示是否加载成功
func (r *RedisCli) HLoad(ctx context.Context, key, field string, v interface{}) (bool, error) {
	if r == nil {
		return false, nil
	}
	var val []byte
	var err error
	if r.Client != nil {
		val, err = r.Client.HGet(ctx, key, field).Bytes()
	} else {
		val, err = r.Cluster.HGet(ctx, key, field).Bytes()
	}
	if errors.Is(err, db.NotFound) || errors.Is(err, redis.Nil) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	if val == nil {
		return false, nil
	}
	if err := json.Unmarshal(val, v); err != nil {
		return false, err
	}
	return true, nil
}

// HStore 序列化数据并储存数据
func (r *RedisCli) HStore(ctx context.Context, key, field string, v interface{}, ex time.Duration) {
	if r == nil {
		return
	}
	if ex == 0 {
		ex = time.Hour * 3
	}
	var pipe redis.Pipeliner
	if r.Client != nil {
		pipe = r.Client.Pipeline()
	} else {
		pipe = r.Cluster.Pipeline()
	}
	pipe.HSet(ctx, key, field, utils.JSONMarshal(v))
	pipe.Expire(ctx, key, ex)
	_, err := pipe.Exec(ctx)
	if err != nil {
		_, err2 := pipe.Exec(ctx)
		if err2 != nil {
			msg := fmt.Sprintf("HStore error for key:%s field:%s first error: %s, second error: %s", key, field, err, err2)
			ilogger.ErrorwCtx(ctx, err2.Error(), "HStore", msg)
		}
	}
}

// Load 加载数据并反序列化
// 返回值第一个参数表示是否加载成功
func (r *RedisCli) Load(ctx context.Context, key string, v interface{}) (bool, error) {
	if r == nil {
		return false, nil
	}
	var val []byte
	var err error
	if r.Client != nil {
		val, err = r.Client.Get(ctx, key).Bytes()
	} else {
		val, err = r.Cluster.Get(ctx, key).Bytes()
	}
	if errors.Is(err, db.NotFound) || errors.Is(err, redis.Nil) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	if val == nil {
		return false, nil
	}
	if err := json.Unmarshal(val, v); err != nil {
		return false, err
	}
	return true, nil
}

// Store 序列化数据并储存数据
func (r *RedisCli) Store(ctx context.Context, key string, v interface{}, ex time.Duration) {
	if r == nil {
		return
	}
	if ex == 0 {
		ex = time.Hour * 3
	}
	if r.Client != nil {
		r.Client.Set(ctx, key, utils.JSONMarshal(v), ex)
	} else {
		r.Cluster.Set(ctx, key, utils.JSONMarshal(v), ex)
	}
}

// IncrEx IncrEx
func (r *RedisCli) IncrEx(ctx context.Context, key string, ex time.Duration) *redis.IntCmd {
	if r == nil {
		return nil
	}
	var rst *redis.IntCmd
	if r.Client != nil {
		rst = r.Client.Incr(ctx, key)
		r.Client.Expire(ctx, key, ex)
	} else {
		rst = r.Cluster.Incr(ctx, key)
		r.Cluster.Expire(ctx, key, ex)
	}
	return rst
}

// GetKeys 通过scan获取keys*
func (r *RedisCli) GetKeys(ctx context.Context, pattern string) ([]string, error) {
	// 初始化游标变量和 key 列表
	var cursor uint64
	var keys []string
	// 匹配模式
	for {
		var results []string
		var nextCursor uint64
		var err error

		// 执行 SCAN 命令
		if r.Client != nil {
			results, nextCursor, err = r.Client.Scan(ctx, cursor, pattern, 100).Result()
		} else {
			results, nextCursor, err = r.Cluster.Scan(ctx, cursor, pattern, 100).Result()
		}
		if err != nil {
			return nil, err
		}
		// 将返回的 key 列表添加到切片中
		keys = append(keys, results...)
		// 更新游标变量，以便进行下一轮扫描
		cursor = nextCursor
		// 如果游标为 0，则表示已经遍历完整个键空间
		if cursor == 0 {
			break
		}
	}
	return keys, nil
}

// Lock 创建一个并发锁
func (r *RedisCli) Lock(ctx context.Context, key string, ex time.Duration) bool {
	if r == nil {
		return false
	}

	if r.Client != nil {
		return r.Client.SetNX(ctx, key, 1, ex).Val()
	}
	return r.Cluster.SetNX(ctx, key, 1, ex).Val()
}

// UnLock 删除一个并发锁
func (r *RedisCli) UnLock(ctx context.Context, key string) error {
	return r.RDel(ctx, key)
}

// RDel 执行redis del操ct
func (r *RedisCli) RDel(ctx context.Context, keys ...string) error {
	return r.doDel(ctx, 3, keys...)
}

// RDel 执行redis del操作
func (r *RedisCli) doDel(ctx context.Context, retry int, keys ...string) error {
	if r == nil {
		return nil
	}

	var err error
	if r.Client != nil {
		_, err = r.Client.Del(ctx, keys...).Result()
	} else {
		_, err = r.Cluster.Del(ctx, keys...).Result()
	}
	if err != nil {
		if retry > 0 {
			return r.doDel(ctx, retry-1, keys...)
		}
		ilogger.ErrorwCtx(ctx, err.Error(), "Operate", "RedisCli Del", "Retry", retry, "Key", utils.JSONMarshal(keys))
		return err
	}
	return nil
}
