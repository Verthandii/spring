package isnowflake

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/redis/go-redis/v9"

	"github.com/Verthandii/spring/db/iredis"
	"github.com/Verthandii/spring/ilogger"
	"github.com/Verthandii/spring/utils"
)

const (
	heartbeatSetInterval = time.Second * 10
)

var node *snowflake.Node

func getSnowFlakeKey() string {
	return "GLOBAL:snowflake_v2"
}

// ID 获取唯一ID
func ID() string {
	return node.Generate().String()
}

// Init 初始化雪花算法
func Init(ctx context.Context, ri *iredis.RedisCli) {
	// 获取容器名称
	key := getSnowFlakeKey()
	// 获取可使用ID
	nodeID := getNodeID(ctx, ri)
	// 使用并发锁控制NodeID
	lockKey := fmt.Sprintf("%s:%d", key, nodeID)
	if !ri.Lock(ctx, lockKey, time.Minute) {
		Init(ctx, ri)
		return
	}

	// 存储当前节点ID
	ri.ZAdd(ctx, key, redis.Z{Member: nodeID, Score: float64(utils.Now())})
	// 创建雪花实例
	var err error
	node, err = snowflake.NewNode(nodeID)
	if err != nil {
		ilogger.Fatalf("snowflake.NewNode err: %s", err.Error())
	}
	// 节点心跳
	go snowFlakeHeartbeat(ctx, ri, key, nodeID)
}

// getNodeID 获取已使用的节点ID
func getNodeID(ctx context.Context, ri *iredis.RedisCli) int64 {
	list, err := ri.ZRange(ctx, getSnowFlakeKey(), 0, -1).Result()
	if err != nil {
		ilogger.Fatal("snowflake.ZRange", err.Error())
	}
	var used []int
	for _, key := range list {
		id := ri.Get(ctx, key).Val()
		if id == "" {
			continue
		}
		u, err := strconv.Atoi(id)
		if err != nil {
			ilogger.Fatalf("snowflake.NewNode strconv.Atoi err: %s", err.Error())
		}
		used = append(used, u)
	}
	// 获取基础节点ID
	base := getBase(1024)
	// 计算可用节点ID
	available := utils.SliceDiff(base, used)
	if available == nil {
		ilogger.Fatal("snowflake.NewNode err: 无可用NodeID")
	}
	return int64(utils.SliceRand(available))
}

// snowFlakeHeartbeat 设置当前心跳
func snowFlakeHeartbeat(ctx context.Context, ri *iredis.RedisCli, key string, nodeID int64) {
	for {
		now := utils.Now()
		// 删除1分钟之前的key
		ri.ZRemRangeByScore(ctx, key, "-inf", fmt.Sprintf("%d", now-60))
		// 重新设置当前key
		ri.ZAdd(ctx, key, redis.Z{Member: nodeID, Score: float64(now)})
		time.Sleep(heartbeatSetInterval)
	}
}

// getBase 获取基础节点id
func getBase(max int) []int {
	var ids []int
	for i := 0; i < max; i++ {
		ids = append(ids, i)
	}
	return ids
}
