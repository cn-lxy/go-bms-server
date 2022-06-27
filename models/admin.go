package models

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
)

// 使用Sqlite3

const (
	ID       = "id"
	ACCOUNT  = "account"
	PASSWORD = "password"
)

type Admin struct {
	gorm.Model
	Account  string `json:"account"`
	Password string `json:"password"`
}

func (admin *Admin) Insert(db *gorm.DB) error {
	db.Create(admin)
	return nil
}

// FIX
func (admin *Admin) Query(db *gorm.DB, key string, value any) error {
	db.First(&admin, fmt.Sprintf("%s = ?", key), value)
	return nil
}

// FIX
// Update by id or account update admin info.
func (admin *Admin) Update(db *gorm.DB, mode string) error {
	newAdmin := Admin{
		Model:    admin.Model,
		Account:  admin.Account,
		Password: admin.Password,
	}

	admin.Query(db, ACCOUNT, admin.Account)

	if admin.Model.ID == 0 {
		return errors.New("not find this model.")
	}
	db.Model(admin).Updates(newAdmin)
	return nil
}

func (admin *Admin) Delete(db *gorm.DB) error {
	db.Delete(admin)
	return nil
}
