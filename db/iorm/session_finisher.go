package iorm

import (
	"database/sql"
	"strings"
	"time"

	"gorm.io/gorm"

	"github.com/Verthandii/spring/ilogger"
	"github.com/Verthandii/spring/it"
)

var (
	nilSS = strings.Builder{}
	nilSV []any
)

// r 返回DB结果 Result
func (s *Session) r(tx *gorm.DB, dest any) *Session {
	s.RetryNum = 0
	r := &Session{
		ctx:          s.ctx,
		o:            s.o,
		tx:           nil,
		tran:         nil,
		cache:        s.cache, // 存or取
		paginate:     s.paginate,
		Error:        tx.Error,
		RetryNum:     0,
		RowsAffected: tx.RowsAffected,
	}
	if dest != nil {
		r.StoreCache(dest)
	}
	return r
}

// retryRaw 重试机制 只有使用原生SQL才能触发重试机制,请使用 Raw 方法
// isStatement 是否为statement模式，如果是: 则ss必填，sv是vars参数
func (s *Session) retryRaw(tx *gorm.DB, isStatement bool, ss strings.Builder, sv []interface{}) bool {
	if tx.Error != nil {
		if s.RetryNum < s.o.retryMaxCount {
			if isStatement {
				if ss.Len() <= 0 {
					return false
				}
				// s.conn.Table("") 的具体解释如下
				// 让db.getInstance()的clone设为0，默认是1，表示开始新的builder，设置statement才不会被重置
				if s.tx != nil {
					s.tx = s.o.db.Table("")
					s.tx.Statement.SQL = ss
					s.tx.Statement.Vars = sv
				}
				if s.tran != nil {
					s.tran = s.o.db.Table("")
					s.tran.Statement.SQL = ss
					s.tran.Statement.Vars = sv
				}
			}
			s.RetryNum++
			time.Sleep(s.o.retryNextTime * time.Duration(s.RetryNum))
			return true
		}
		ilogger.ErrorwCtx(s.ctx, tx.Error.Error(), "Operate", "Retry", "RetryNum", s.RetryNum, "SQL", ss.String(), "Vars", sv)
	}
	return false
}

// retryIf 重试机制 只有部分写入操作可以出发
func (s *Session) retryIf(tx *gorm.DB, value interface{}) bool {
	if tx.Error != nil {
		if s.RetryNum < s.o.retryMaxCount {
			s.RetryNum++
			time.Sleep(s.o.retryNextTime * time.Duration(s.RetryNum))
			return true
		}
		ilogger.ErrorwCtx(s.ctx, tx.Error.Error(), "Operate", "Retry", "RetryNum", s.RetryNum, "Data", value, "SQL", tx.Statement)
	}
	return false
}

// Insert 用原生SQL执行insert并加载主键id
func (s *Session) Insert(sql string, id *int64, values ...interface{}) *Session {
	if s.Error != nil {
		return s
	}
	s.setup()
	defer s.teardown()
	result, err := s.DB().ConnPool.ExecContext(s.ctx, sql, values...)
	if err != nil {
		s.Error = err
		return s
	}
	if id != nil {
		*id, err = result.LastInsertId()
		if err != nil {
			s.Error = err
			return s
		}
	}
	s.RowsAffected, _ = result.RowsAffected()
	return s
}

// region 重写GORM FinisherAPI方法

// Create insert the value into database
func (s *Session) Create(value interface{}) *Session {
	if s.Error != nil {
		return s
	}
	s.setup()
	defer s.teardown()
	tx := s.DB().Create(value)
	if s.retryIf(tx, value) {
		return s.Create(value)
	}
	return s.r(tx, nil)
}

// CreateInBatches insert the value in batches into database
func (s *Session) CreateInBatches(value interface{}, batchSize int) *Session {
	if s.Error != nil {
		return s
	}
	s.setup()
	defer s.teardown()
	return s.r(s.DB().CreateInBatches(value, batchSize), nil)
}

// Save update value in database, if the value doesn't have primary key, will insert it
func (s *Session) Save(value interface{}) *Session {
	if s.Error != nil {
		return s
	}
	s.setup()
	defer s.teardown()
	tx := s.DB().Save(value)
	if s.retryIf(tx, value) {
		return s.Save(value)
	}
	return s.r(tx, nil)
}

// First find first record that match given conditions, order by primary key
func (s *Session) First(dest interface{}, conds ...interface{}) *Session {
	if s.Error != nil {
		return s
	}
	s.setup()
	defer s.teardown()
	if s.LoadCache(dest) {
		return s
	}
	return s.r(s.DB().First(dest, conds...), dest)
}

// Take return a record that match given conditions, the order will depend on the database implementation
func (s *Session) Take(dest interface{}, conds ...interface{}) *Session {
	if s.Error != nil {
		return s
	}
	s.setup()
	defer s.teardown()
	if s.LoadCache(dest) {
		return s
	}
	return s.r(s.DB().Take(dest, conds...), dest)
}

// Last find last record that match given conditions, order by primary key
func (s *Session) Last(dest interface{}, conds ...interface{}) *Session {
	if s.Error != nil {
		return s
	}
	s.setup()
	defer s.teardown()
	if s.LoadCache(dest) {
		return s
	}
	return s.r(s.DB().Last(dest, conds...), dest)
}

func (s *Session) GetPaginate() *Paginate {
	return s.paginate
}

// Find find records that match given conditions
func (s *Session) Find(dest interface{}, conds ...interface{}) *Session {
	if s.Error != nil {
		return s
	}
	s.setup()
	defer s.teardown()
	if s.LoadCache(dest) {
		return s
	}
	if tx := s.loadPaginate(dest); tx.Error != nil {
		return tx
	}
	return s.r(s.DB().Find(dest, conds...), dest)
}

// FindInBatches find records in batches
func (s *Session) FindInBatches(dest interface{}, batchSize int, fc func(tx *gorm.DB, batch int) error) *Session {
	if s.Error != nil {
		return s
	}
	s.setup()
	defer s.teardown()
	return s.r(s.DB().FindInBatches(dest, batchSize, fc), nil)
}

func (s *Session) FirstOrInit(dest interface{}, conds ...interface{}) *Session {
	if s.Error != nil {
		return s
	}
	s.setup()
	defer s.teardown()
	if s.LoadCache(dest) {
		return s
	}
	return s.r(s.DB().FirstOrInit(dest, conds...), dest)
}

func (s *Session) FirstOrCreate(dest interface{}, conds ...interface{}) *Session {
	if s.Error != nil {
		return s
	}
	s.setup()
	defer s.teardown()
	if s.LoadCache(dest) {
		return s
	}
	return s.r(s.DB().FirstOrCreate(dest, conds...), dest)
}

// Update update attributes with callbacks, refer: https://gorm.io/docs/update.html#Update-Changed-Fields
func (s *Session) Update(column string, value interface{}) *Session {
	if s.Error != nil {
		return s
	}
	s.setup()
	defer s.teardown()
	return s.r(s.DB().Update(column, value), nil)
}

// Updates update attributes with callbacks, refer: https://gorm.io/docs/update.html#Update-Changed-Fields
func (s *Session) Updates(values interface{}) *Session {
	if s.Error != nil {
		return s
	}
	s.setup()
	defer s.teardown()

	switch v := values.(type) {
	case it.H:
		values = map[string]any(v)
	}

	return s.r(s.DB().Updates(values), nil)
}

func (s *Session) UpdateColumn(column string, value interface{}) *Session {
	if s.Error != nil {
		return s
	}
	s.setup()
	defer s.teardown()
	return s.r(s.DB().UpdateColumn(column, value), nil)
}

func (s *Session) UpdateColumns(values interface{}) *Session {
	if s.Error != nil {
		return s
	}
	s.setup()
	defer s.teardown()
	return s.r(s.DB().UpdateColumns(values), nil)
}

// Delete delete value match given conditions, if the value has primary key, then will including the primary key as condition
func (s *Session) Delete(value interface{}, conds ...interface{}) *Session {
	if s.Error != nil {
		return s
	}
	s.setup()
	defer s.teardown()
	return s.r(s.DB().Delete(value, conds...), nil)
}

func (s *Session) Count(count *int64) *Session {
	if s.Error != nil {
		return s
	}
	s.setup()
	defer s.teardown()
	if s.LoadCache(count) {
		return s
	}
	return s.r(s.DB().Count(count), count)
}

func (s *Session) Row() *sql.Row {
	s.setup()
	defer s.teardown()
	return s.DB().Row()
}

func (s *Session) Rows() (*sql.Rows, error) {
	s.setup()
	defer s.teardown()
	return s.DB().Rows()
}

// Scan scan value to a struct
func (s *Session) Scan(dest interface{}) *Session {
	if s.Error != nil {
		return s
	}
	s.setup()
	defer s.teardown()
	if s.LoadCache(dest) {
		return s
	}
	ss, sv := s.DB().Statement.SQL, s.DB().Statement.Vars
	tx := s.DB().Scan(dest)
	if s.retryRaw(tx, true, ss, sv) {
		return s.Scan(dest)
	}
	return s.r(tx, dest)
}

// Pluck used to query single column from a model as a map
//
//	var ages []int64
//	db.Model(&users).Pluck("age", &ages)
func (s *Session) Pluck(column string, dest interface{}) *Session {
	if s.Error != nil {
		return s
	}
	s.setup()
	defer s.teardown()
	if s.LoadCache(dest) {
		return s
	}
	return s.r(s.DB().Pluck(column, dest), dest)
}

func (s *Session) ScanRows(rows *sql.Rows, dest interface{}) *Session {
	if s.Error != nil {
		return s
	}
	s.setup()
	defer s.teardown()
	if s.LoadCache(dest) {
		return s
	}
	s.Error = s.DB().ScanRows(rows, dest)
	s.StoreCache(dest)
	return s
}

// Transaction start a transaction as a block, return error will rollback, otherwise to commit.
func (s *Session) Transaction(fc func() error, opts ...*sql.TxOptions) (err error) {
	s.setup()
	defer s.teardown()
	return s.DB().Transaction(func(tx *gorm.DB) error {
		s.beginTran(tx)
		defer s.endTran()
		return fc()
	}, opts...)
}

// Begin begins a transaction
func (s *Session) Begin(opts ...*sql.TxOptions) *Session {
	if s.Error != nil {
		return s
	}
	s.setup()
	defer s.teardown()
	s.tx = s.DB().Begin(opts...)
	return s
}

// Commit commit a transaction
func (s *Session) Commit() *Session {
	if s.Error != nil {
		return s
	}
	s.setup()
	defer s.teardown()
	s.tx = s.DB().Commit()
	return s
}

// Rollback rollback a transaction
func (s *Session) Rollback() *Session {
	if s.Error != nil {
		return s
	}
	s.setup()
	defer s.teardown()
	s.tx = s.DB().Rollback()
	return s
}

func (s *Session) SavePoint(name string) *Session {
	if s.Error != nil {
		return s
	}
	s.setup()
	defer s.teardown()
	s.tx = s.DB().SavePoint(name)
	return s
}

func (s *Session) RollbackTo(name string) *Session {
	if s.Error != nil {
		return s
	}
	s.setup()
	defer s.teardown()
	s.tx = s.DB().RollbackTo(name)
	return s
}

// Exec execute raw sql
func (s *Session) Exec(sql string, values ...interface{}) *Session {
	if s.Error != nil {
		return s
	}
	s.setup()
	defer s.teardown()
	tx := s.DB().Exec(sql, values...)
	if s.retryRaw(tx, false, nilSS, nilSV) {
		return s.Exec(sql, values...)
	}
	return s.r(tx, nil)
}

// Preview 预览SQL
func (s *Session) Preview(sql string, values ...interface{}) string {
	return s.DB().Dialector.Explain(sql, values...)
}

// endregion
