package iorm

import (
	"errors"

	"gorm.io/gorm"

	"github.com/Verthandii/spring/db"
)

type errorPlugin struct{}

func (e *errorPlugin) Name() string {
	return "error_plugin"
}

func (e *errorPlugin) Initialize(tx *gorm.DB) error {
	_ = tx.Callback().Query().After("gorm:row").Register("error_plugin:row", e.transformErr)
	_ = tx.Callback().Query().After("gorm:raw").Register("error_plugin:raw", e.transformErr)
	_ = tx.Callback().Query().After("gorm:after_query").Register("error_plugin:after_query", e.transformErr)
	return nil
}

func (e *errorPlugin) transformErr(tx *gorm.DB) {
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		tx.Error = errors.Join(gorm.ErrRecordNotFound, db.NotFound)
	}
}
