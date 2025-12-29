package iorm

import (
	"context"

	"gorm.io/gorm"
)

// Session DB会话
type Session struct {
	ctx          context.Context
	o            *IORM     // DB连接
	tx           *gorm.DB  // 查询会话
	tran         *gorm.DB  // 事务
	cache        *Cache    // 缓存组件
	paginate     *Paginate // 分页组件
	Error        error     // 会话错误
	RetryNum     int64     // 重试次数
	RowsAffected int64     // 受影响的行数
}

func NewSession(ctx context.Context, o *IORM) *Session {
	return &Session{
		ctx: ctx,
		o:   o,
	}
}

// Clone 克隆会话
func (s *Session) Clone(ctx context.Context) *Session {
	return &Session{
		ctx: ctx,
		o:   s.o,
	}
}

// setup 前置
func (s *Session) setup() {
	s.Error = nil
	s.RowsAffected = 0
}

// teardown 后置
func (s *Session) teardown() {
	s.cache = nil
	s.paginate = nil
	s.tx = nil
}
