package database

import (
	"github.com/fiufit/users/contracts"
	"gorm.io/gorm"
)

func Paginate(value interface{}, pagination *contracts.Pagination, db *gorm.DB) func(*gorm.DB) *gorm.DB {
	db.Model(value).Count(&pagination.TotalRows)

	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(pagination.ToOffset()).Limit(pagination.ToLimit())
	}
}
