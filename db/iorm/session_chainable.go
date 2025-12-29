package iorm

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// DB 获取当前会话连接
func (s *Session) DB() *gorm.DB {
	if s.tx != nil {
		return s.tx
	}
	if s.tran != nil {
		return s.tran
	}
	return s.o.db
}

// beginTran 开始事务
func (s *Session) beginTran(tx *gorm.DB) *Session {
	s.tran = tx
	return s
}

// endTran 结束事务
func (s *Session) endTran() *Session {
	s.tran = nil
	return s
}

func (s *Session) Field(q string, dest interface{}) *Session {
	return s.Select(q).Scan(dest)
}

// region 重写GORM ChainableAPI方法

// Model specify the model you would like to run db operations
//
//	// update all users's name to `hello`
//	db.Model(&User{}).Update("name", "hello")
//	// if user's primary key is non-blank, will use it as condition, then will only update the user's name to `hello`
//	db.Model(&user).Update("name", "hello")
func (s *Session) Model(value interface{}) *Session {
	s.tx = s.DB().Model(value)
	return s
}

// Clauses Add clauses
func (s *Session) Clauses(conds ...clause.Expression) *Session {
	s.tx = s.DB().Clauses(conds...)
	return s
}

// Table specify the table you would like to run db operations
func (s *Session) Table(name string, args ...interface{}) *Session {
	s.tx = s.DB().Table(name, args...)
	return s
}

// Distinct specify distinct fields that you want querying
func (s *Session) Distinct(args ...interface{}) *Session {
	s.tx = s.DB().Distinct(args)
	return s
}

// Select specify fields that you want when querying, creating, updating
func (s *Session) Select(query interface{}, args ...interface{}) *Session {
	s.tx = s.DB().Select(query, args...)
	return s
}

// Omit specify fields that you want to ignore when creating, updating and querying
func (s *Session) Omit(columns ...string) *Session {
	s.tx = s.DB().Omit(columns...)
	return s
}

// Where add conditions
func (s *Session) Where(query interface{}, args ...interface{}) *Session {
	s.tx = s.DB().Where(query, args...)
	return s
}

// Not add NOT conditions
func (s *Session) Not(query interface{}, args ...interface{}) *Session {
	s.tx = s.DB().Not(query, args...)
	return s
}

// Or add OR conditions
func (s *Session) Or(query interface{}, args ...interface{}) *Session {
	s.tx = s.DB().Or(query, args...)
	return s
}

// Joins specify Joins conditions
//
//	db.Joins("Account").Find(&user)
//	db.Joins("JOIN emails ON emails.user_id = users.id AND emails.email = ?", "jinzhu@example.org").Find(&user)
func (s *Session) Joins(query string, args ...interface{}) *Session {
	s.tx = s.DB().Joins(query, args...)
	return s
}

// Group specify the group method on the find
func (s *Session) Group(name string) *Session {
	s.tx = s.DB().Group(name)
	return s
}

// Having specify HAVING conditions for GROUP BY
func (s *Session) Having(query interface{}, args ...interface{}) *Session {
	s.tx = s.DB().Having(query, args...)
	return s
}

// Order specify order when retrieve records from database
//
//	db.Order("name DESC")
//	db.Order(clause.OrderByColumn{Column: clause.Column{Name: "name"}, Desc: true})
func (s *Session) Order(value interface{}) *Session {
	s.tx = s.DB().Order(value)
	return s
}

// Limit specify the number of records to be retrieved
func (s *Session) Limit(limit int) *Session {
	s.tx = s.DB().Limit(limit)
	return s
}

// Offset specify the number of records to skip before starting to return the records
func (s *Session) Offset(offset int) *Session {
	s.tx = s.DB().Offset(offset)
	return s
}

// Scopes pass current database connection to arguments `func(DB) DB`, which could be used to add conditions dynamically
//
//	func AmountGreaterThan1000(db *gorm.DB) *gorm.DB {
//	    return db.Where("amount > ?", 1000)
//	}
//
//	func OrderStatus(status []string) func (db *gorm.DB) *gorm.DB {
//	    return func (db *gorm.DB) *gorm.DB {
//	        return db.Scopes(AmountGreaterThan1000).Where("status in (?)", status)
//	    }
//	}
//
//	db.Scopes(AmountGreaterThan1000, OrderStatus([]string{"paid", "shipped"})).Find(&orders)
func (s *Session) Scopes(funcs ...func(*gorm.DB) *gorm.DB) *Session {
	s.tx = s.DB().Scopes(funcs...)
	return s
}

// Preload preload associations with given conditions
//
//	db.Preload("Orders", "state NOT IN (?)", "cancelled").Find(&users)
func (s *Session) Preload(query string, args ...interface{}) *Session {
	s.tx = s.DB().Preload(query, args...)
	return s
}

func (s *Session) Attrs(attrs ...interface{}) *Session {
	s.tx = s.DB().Attrs(attrs...)
	return s
}

func (s *Session) Assign(attrs ...interface{}) *Session {
	s.tx = s.DB().Assign(attrs...)
	return s
}

func (s *Session) Unscoped() *Session {
	s.tx = s.DB().Unscoped()
	return s
}

func (s *Session) Raw(sql string, values ...interface{}) *Session {
	s.tx = s.DB().Raw(sql, values...)
	return s
}

// endregion
