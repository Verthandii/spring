package iorm

import (
	"strings"
	"time"

	"gorm.io/gorm"
)

const (
	metricsInfoKey = "__metrics_info_key"
)

var _ gorm.Plugin = &metricsPlugin{}

// gorm metrics插件
type metricsPlugin struct {
	addr   string
	dbName string
}

type metricsInfo struct {
	startAt time.Time
	addr    string
	dbName  string
}

func newMetricsPlugin(addr, dbName string) *metricsPlugin {
	return &metricsPlugin{
		addr:   addr,
		dbName: dbName,
	}
}

func newMetricsInfo(addr, dbName string) *metricsInfo {
	return &metricsInfo{
		startAt: time.Now(),
		addr:    addr,
		dbName:  dbName,
	}
}

func (m *metricsPlugin) Name() string {
	return "metrics"
}

func (s *metricsPlugin) Initialize(db *gorm.DB) error {
	var err error
	for _, v := range beforeHookPoint(db) {
		if err != nil {
			continue
		}
		err = v("metrics:before", s.before)
	}
	for _, v := range afterHookPoint(db) {
		if err != nil {
			continue
		}
		err = v("metrics:after", s.after)
	}
	return err
}

// before 钩子, 用于记录请求耗时
func (s *metricsPlugin) before(db *gorm.DB) {
	info := newMetricsInfo(s.addr, s.dbName)
	db.Set(metricsInfoKey, info)
}

// after 钩子, 用于记录请求耗时
func (s *metricsPlugin) after(db *gorm.DB) {
	info, ok := db.Get(metricsInfoKey)
	if !ok {
		return
	}

	metricsInfo := info.(*metricsInfo)
	tableName := strings.Split(db.Statement.Table, " ")[0]
	elapsed := time.Since(metricsInfo.startAt).Milliseconds()
	metricReqDur.WithLabelValues(s.addr, s.dbName, tableName, getOperation(db)).Observe(float64(elapsed))
}

// https://github.com/go-gorm/gorm/blob/master/callbacks/callbacks.go 已定义的钩子
// 钩子函数, 用于注册metrics方法
func beforeHookPoint(db *gorm.DB) []func(name string, fn func(*gorm.DB)) error {
	return []func(name string, fn func(*gorm.DB)) error{
		db.Callback().Create().Before("gorm:before_create").Register,
		db.Callback().Query().Before("gorm:query").Register,
		db.Callback().Delete().Before("gorm:before_delete").Register,
		db.Callback().Update().Before("gorm:before_update").Register,
		db.Callback().Row().Before("gorm:row").Register,
		db.Callback().Raw().Before("gorm:raw").Register,
	}
}

func afterHookPoint(db *gorm.DB) []func(name string, fn func(*gorm.DB)) error {
	return []func(name string, fn func(*gorm.DB)) error{
		db.Callback().Create().After("gorm:after_create").Register,
		db.Callback().Query().After("gorm:after_query").Register,
		db.Callback().Delete().After("gorm:after_delete").Register,
		db.Callback().Update().After("gorm:after_update").Register,
		db.Callback().Row().After("gorm:row").Register,
		db.Callback().Raw().After("gorm:raw").Register,
	}
}

func getOperation(db *gorm.DB) string {
	sql := db.Statement.SQL.String()
	if sql != "" {
		return strings.ToUpper(strings.Split(sql, " ")[0])
	}

	cl := db.Statement.BuildClauses
	if len(cl) == 0 {
		return ""
	}

	return strings.ToUpper(cl[0])
}
