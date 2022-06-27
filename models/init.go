package models

import (
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB
var err error

func init() {
	// 数据库初始化
	db, err = gorm.Open(sqlite.Open("bms.sqlite"), &gorm.Config{})
	if err != nil {
		panic("failed connect database!")
	}
	// db.AutoMigrate(&Admin{})
	fmt.Println("init completed!")
}
