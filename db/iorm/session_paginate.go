package iorm

import (
	"reflect"

	"gorm.io/gorm/clause"

	"github.com/Verthandii/spring/transport/ihttp"
)

// Paginate 分页组件
type Paginate struct {
	Page  int   // 页码
	Limit int   // 每页条数
	Total int64 // 总条数
}

// PageData 分页数据
type PageData struct {
	List     interface{}
	Paginate *Paginate
}

// GPaginate 使用Gin 加载分页数据
func (s *Session) GPaginate(c *ihttp.Context) *Session {
	var bind Paginate
	if err := c.ShouldBind(&bind); err != nil {
		s.Error = err
		return s
	}
	return s.Paginate(bind.Page, bind.Limit)
}

// Paginate 加载分页数据
func (s *Session) Paginate(page, limit int) *Session {
	s.setup()
	bind := Paginate{
		Page:  page,
		Limit: limit,
	}
	if bind.Page <= 0 {
		bind.Page = 1
	}
	if bind.Limit == 0 {
		bind.Limit = 20
	}
	if bind.Limit > 5000 {
		bind.Limit = 5000
	}
	s.paginate = &bind
	return s
}

// loadPaginate 加载分页
func (s *Session) loadPaginate(dest interface{}) *Session {
	page := s.paginate
	if page == nil {
		return s
	}

	// 获取Model
	model := s.DB().Statement.Model
	if model == nil {
		// 获取Dest
		destSlice := reflect.Indirect(reflect.ValueOf(dest))
		destType := destSlice.Type().Elem()
		model = reflect.New(destType).Interface()
	}

	// 保存当前的 Selects 和 OrderBy 条件
	selects := s.DB().Statement.Selects
	order := s.DB().Statement.Clauses["ORDER BY"]

	// 清除 OrderBy 条件以进行 Count 查询 使用空结构体{}来清除
	s.DB().Statement.Clauses["ORDER BY"] = clause.Clause{}

	if tx := s.DB().Select("COUNT(1)").Model(model).Count(&page.Total); tx.Error != nil {
		return s.r(tx, nil)
	}

	// 恢复原来的 Selects 和 OrderBy 条件
	s.DB().Statement.Selects = selects
	s.DB().Statement.Clauses["ORDER BY"] = order

	q := s.Limit(page.Limit)
	if page.Page > 1 {
		q.Offset((page.Page - 1) * page.Limit)
	}
	return q
}

// PD 获取分页数据 PageData
func (s *Session) PD(list interface{}) *PageData {
	return &PageData{
		List:     list,
		Paginate: s.paginate,
	}
}

// GetPDTotal 获取分页总条数
func (s *Session) GetPDTotal() int64 {
	return s.paginate.Total
}

// SetCount 设置总条数
func (s *Session) SetCount(total int64) {
	s.paginate.Total = total
	return
}
