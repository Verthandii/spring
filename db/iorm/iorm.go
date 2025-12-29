package iorm

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"

	"github.com/Verthandii/spring/db/icache"
	"github.com/Verthandii/spring/db/iredis"
	"github.com/Verthandii/spring/env"
	"github.com/Verthandii/spring/ilogger"
)

const (
	RetryMaxCount = 3                     // 最大重试次数
	RetryNextTime = time.Millisecond * 10 // 下一次重试时间间隔
)

type Option func(o *IORM)

func WithRedis(r *iredis.RedisCli) Option {
	return func(o *IORM) {
		o.redis = r
	}
}

func WithMCache(mcache *icache.MemCache) Option {
	return func(o *IORM) {
		o.mcache = mcache
	}
}

func WithRetryMaxCount(count int64) Option {
	return func(o *IORM) {
		o.retryMaxCount = count
	}
}

func WithRetryNextTime(t time.Duration) Option {
	return func(o *IORM) {
		o.retryNextTime = t
	}
}

type IORM struct {
	db     *gorm.DB
	redis  *iredis.RedisCli
	mcache *icache.MemCache

	retryMaxCount int64
	retryNextTime time.Duration
}

// New 新建mysql连接
func New(cfg Config, opts ...Option) (*IORM, error) {
	db, err := newGORM(cfg)
	if err != nil {
		return nil, err
	}

	o := &IORM{
		db:            db,
		retryMaxCount: RetryMaxCount,
		retryNextTime: RetryNextTime,
	}

	for _, opt := range opts {
		opt(o)
	}

	return o, err
}

func newGORM(cfg Config) (*gorm.DB, error) {
	gc := &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "",             // table name prefix, table for `User` would be `t_users`
			SingularTable: true,           // use singular table name, table for `User` would be `user` with this option enabled
			NoLowerCase:   !cfg.LowerCase, // skip the snake_casing of names
		},
		PrepareStmt: false, // 这里必须关闭,每一条SQL会像服务器Prepare一次、再通过参数execute一次，交互次数变多，然后对计算节点的cache压力也变大了
		// 如果开启之后,实际没有达到prepare一次 execute多次的效果 关闭此参数可以让数据库QPS降低一倍
	}
	if !env.IsStagingOrProd() {
		level := logger.LogLevel(cfg.LoggerLevel)
		gc.Logger = logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				SlowThreshold:             time.Second, // 慢 SQL 阈值
				LogLevel:                  level,       // 日志级别
				IgnoreRecordNotFoundError: true,        // 是否忽略 ErrRecordNotFound（记录未找到）错误
				Colorful:                  true,        // 是否启用彩色打印
			},
		)
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local&interpolateParams=true",
		cfg.UserName, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)

	ilogger.Infof("【MySQL】开始连接 %s", cfg.Host)

	conn, err := gorm.Open(mysql.Open(dsn), gc)
	if err != nil {
		return nil, err
	}

	db, err := conn.DB()
	if err != nil {
		return nil, err
	}

	db.SetConnMaxIdleTime(time.Second * time.Duration(cfg.IdleTime))
	db.SetConnMaxLifetime(time.Hour)
	db.SetMaxOpenConns(cfg.Open)
	db.SetMaxIdleConns(cfg.Idle)

	err = conn.Use(&errorPlugin{})
	if err != nil {
		return nil, err
	}
	// 默认开启metrics，可以配置关闭
	if !cfg.DisableMetrics {
		mp := newMetricsPlugin(fmt.Sprintf("%s:%d", cfg.Host, cfg.Port), cfg.DBName)
		err = conn.Use(mp)
		if err != nil {
			return nil, err
		}
		_statCollector.registerClient(&statClient{
			addr:   fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
			dbName: cfg.DBName,
			inner:  conn,
		})
	}

	ilogger.Infof("【MySQL】连接成功 %s", cfg.Host)

	return conn, nil
}
