package models

import (
	"gorm.io/gorm"
)

type DBUtil interface {
	Insert(db *gorm.DB) error
	Query(db *gorm.DB, key string, value any) error
	Update(db *gorm.DB, mode string) error
	Delete(db *gorm.DB) error
}
