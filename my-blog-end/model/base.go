package model

import "gorm.io/gorm"

// Paginate 通用分页
func Paginate(page, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if pageSize == 0 {
			pageSize = 10
		}
		return db.Offset(page * pageSize).Limit(pageSize)
	}
}
