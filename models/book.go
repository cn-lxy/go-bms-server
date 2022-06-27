package models

import (
	"fmt"

	"gorm.io/gorm"
)

const (
	BOOK_ID         = "id"
	BOOK_NAME       = "name"
	BOOK_ISBN       = "isbn"
	BOOK_TYPE       = "type"
	BOOK_AUTHOR     = "author"
	BOOK_PUBLIC     = "public"
	BOOK_PUBLICTIME = "public_time"
)

type BookType struct {
	gorm.Model
	TypeName string
}

func CreateBookType(typeName string) error {
	db.AutoMigrate(&BookType{})
	return db.Create(&BookType{TypeName: typeName}).Error
}

type Book struct {
	gorm.Model
	Name       string `json:"name"`
	Author     string `json:"author"`
	BookTypeId uint   `json:"type" gorm:"foreignkey:BookTypeID"`
	Public     string `json:"public"`
	PublicTime string `json:"publicTime"`
	Stock      int    `json:"stock"`
}

// 查询图书
func (book *Book) Query(db *gorm.DB, key string, value any) error {
	queryParams := fmt.Sprintf("%s = ?", key)
	db.First(&book, queryParams, value)
	return nil
}

func (book *Book) Insert(db *gorm.DB) error {
	db.Create(book)
	return nil
}

// FIX
// Update by name
func (book *Book) Update(db *gorm.DB, mode string) error {
	newBook := Book{
		Model:      book.Model,
		Name:       book.Name,
		Author:     book.Author,
		Public:     book.Public,
		PublicTime: book.PublicTime,
	}

	book.Query(db, BOOK_NAME, book.Name)

	db.Model(book).Updates(&newBook)
	return nil
}

func (book *Book) Delete(db *gorm.DB) error {
	db.Delete(book)
	return nil
}

// rebuild

func InsertBook(book *Book) error {
	db.AutoMigrate(&Book{})
	result := db.Create(book)
	return result.Error
}

// UpdateBook 通过id更新图书信息
func UpdateBook(book *Book) error {
	db.AutoMigrate(&Book{})
	book.Query(db, BOOK_ID, book.ID)
	return db.Model(book).Updates(&book).Error
}

// DeleteBook 根据主键删除图书信息
func DeleteBook(id int) error {
	db.AutoMigrate(&Book{})
	// 根据主键删除
	return db.Delete(&Book{}, id).Error
}

// GetBook 获取一本图书，通过name或isbn
func GetBook(byType string, value string) (Book, error) {
	db.AutoMigrate(&Book{})
	queryParams := fmt.Sprintf("%s = ?", byType)
	book := Book{}
	result := db.First(&book, queryParams, value)
	return book, result.Error
}

// GetBooksSlice 获取多个图书,通过typeid，没有typeid(-1)默认返回
func GetBookSlice(typeid int, limit, offset int) ([]Book, error) {
	db.AutoMigrate(&Book{})
	books := []Book{}
	var result *gorm.DB
	if typeid != -1 {
		result = db.Joins("left join book_types on book_types.id = book.id").Where("type <> ?", typeid).Limit(limit).Offset(offset).Find(&books)
	} else {
		result = db.Limit(limit).Offset(offset).Find(&books)
	}
	return books, result.Error
}
