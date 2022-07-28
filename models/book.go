package models

import (
	"fmt"
	"strconv"

	"github.com/cn-lxy/bms_go/utils/dbutil"
)

const (
	Book_RegisterSuccess = iota
	Book_RegisterExisted
	Book_RegisterError

	By_BookName
	By_BookISBN
	By_BookDefault
	By_BookTypeId
)

type Book struct {
	Isbn       string `json:"isbn"`        // ISBN号
	Name       string `json:"name"`        // 书名
	TypeId     uint64 `json:"type_id"`     // 类别id
	Author     string `json:"author"`      // 作者
	Public     string `json:"public"`      // 出版社
	PublicDate string `json:"public_date"` // 出版时间
	Stock      uint64 `json:"stock"`       // 库存量
}

type BookManger struct{}

// 添加书籍
func (b BookManger) InsertBook(book *Book) error {
	querySql := fmt.Sprintf("select name from `books` where isbn = '%s'", book.Isbn)
	results := dbutil.Query(querySql)
	if len(results) != 0 {
		return fmt.Errorf("[insert book]: book is existed!\n")
	}

	insertSql := fmt.Sprintf(
		"insert into `books` (isbn, name, type_id, author, public, public_date, register, stock) "+
			"values ('%s', '%s', %d, '%s', '%s', '%s', now(), %d)",
		book.Isbn,
		book.Name,
		book.TypeId,
		book.Author,
		book.Public,
		book.PublicDate,
		book.Stock)
	if err := dbutil.Update(insertSql); err != nil {
		return err
	}
	return nil
}

// 获取一本书籍：isbn, name
func (b BookManger) GetBook(by int, value string) (*Book, error) {
	var whereName string
	switch by {
	case By_BookName:
		whereName = "name"
		break
	case By_BookISBN:
		whereName = "isbn"
		break
	}
	selectSql := fmt.Sprintf("select * from `books` where %s = '%s'", whereName, value)
	results := dbutil.Query(selectSql)
	if len(results) == 0 {
		return nil, fmt.Errorf("%s", "无查询结果")
	}

	return mapToBook(results[0])
}

// 获取图书集合: typeId, default
func (b BookManger) GetBooks(limit, offset int, by int, value ...any) ([]Book, error) {
	var sql string
	switch by {
	case By_BookTypeId:
		sql = fmt.Sprintf("select * from `books` where type_id = %d limit %d offset %d", value[0].(int), limit, offset)
		break
	case By_BookDefault:
		sql = fmt.Sprintf("select * from `books` limit %d offset %d", limit, offset)
		break
	}
	results := dbutil.Query(sql)
	return mapsToBooks(results)
}

// 删除书籍
func (b BookManger) DeleteBook(isbn string) error {
	sql := "delete from `books` where isbn = ?"
	args := []any{isbn}
	return dbutil.Update(sql, args...)
}

// 更新书籍
func (b BookManger) UpdateBook(book *Book) error {
	// sql := "update from `books` set "

	// if book.Isbn == "" {
	// 	return fmt.Errorf("[update book]: isbn is nil\n")
	// }

	// args := make([]any, 0)

	// if book.Name != "" {
	// 	sql += "name = ? ,"
	// 	args = append(args, book.Name)
	// }
	// if book.TypeId != 0 {
	// 	sql += "type_id = ? ,"
	// 	args = append(args, book.TypeId)
	// }
	// if book.Author != "" {
	// 	sql += "author = ? ,"
	// 	args = append(args, book.Author)
	// }
	// if book.Public != "" {
	// 	sql += "public = ? ,"
	// 	args = append(args, book.Public)
	// }
	// if book.PublicDate != "" {
	// 	sql += "public_date = ? ,"
	// 	args = append(args, book.PublicDate)
	// }
	// if book.Stock >= 0 {
	// 	sql += "stock = ? "
	// 	args = append(args, book.Stock)
	// }
	// sql += " where isbn = ? "

	if book.Isbn == "" {
		return fmt.Errorf("[update book]: isbn is nil\n")
	}
	sql := "update books set name = ?, author = ?, type_id = ?, public = ?, public_date = ?, stock = ? where isbn = ?"
	args := []any{book.Name, book.Author, book.TypeId, book.Public, book.PublicDate, book.Stock, book.Isbn}

	return dbutil.Update(sql, args...)
}

// 将map slice 转换成 结构体 slice
func mapsToBooks(maps []map[string]any) ([]Book, error) {
	books := make([]Book, 0, len(maps)) // make slice 得注意设置 len 和 cap
	for _, m := range maps {
		b, err := mapToBook(m)
		if err != nil {
			return nil, err
		}
		books = append(books, *b)
	}
	return books, nil
}

// 将 map 转换成 Book 结构体
func mapToBook(m map[string]any) (*Book, error) {
	typeId, err := strconv.ParseUint(m["type_id"].(string), 0, 64)
	if err != nil {
		return nil, err
	}
	stock, err := strconv.ParseUint(m["stock"].(string), 0, 64)
	if err != nil {
		return nil, err
	}
	return &Book{
		Isbn:       m["isbn"].(string),
		Name:       m["name"].(string),
		TypeId:     typeId,
		Author:     m["author"].(string),
		Public:     m["public"].(string),
		PublicDate: m["public_date"].(string),
		Stock:      stock,
	}, nil
}
