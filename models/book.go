package models

import (
	"fmt"
	"log"
	"strconv"

	"github.com/cn-lxy/bms_go/utils/dbutil"
)

type Book struct {
	Isbn       string // ISBN号
	Name       string // 书名
	TypeId     int    // 类别id
	Author     string // 作者
	Public     string // 出版社
	PublicDate string // 出版时间
	Stock      int    // 库存量
}

const (
	Book_RegisterSuccess = iota
	Book_RegisterExisted
	Book_RegisterError

	By_BookName
	By_BookISBN
	By_BookDefault
	By_BookTypeId
)

// 添加书籍
func BookRegister(book *Book) int {
	querySql := fmt.Sprintf("select name from `books` where isbn = '%s'", book.Isbn)
	results := dbutil.Query(querySql)
	if len(results) != 0 {
		return Book_RegisterExisted
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
		log.Fatal(err.Error())
		return Book_RegisterError
	}
	return Book_RegisterSuccess
}

// 获取一本书籍：isbn, name
func GetBook(by int, value string) (*Book, error) {
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

	return mapToBook(results[0]), nil

}

// 获取图书集合: typeId, default
func GetBooks(limit, offset int, by int, value ...any) []Book {
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

// 将map slice 转换成 结构体 slice
func mapsToBooks(maps []map[string]any) []Book {
	books := make([]Book, 0, len(maps)) // make slice 得注意设置 len 和 cap
	for _, m := range maps {
		books = append(books, *mapToBook(m))
	}
	return books
}

// 将 map 转换成 Book 结构体
func mapToBook(m map[string]any) *Book {
	var typeId int
	var stock int
	typeId, _ = strconv.Atoi(m["type_id"].(string)) // 断言
	stock, _ = strconv.Atoi(m["stock"].(string))
	return &Book{
		Isbn:       m["isbn"].(string),
		Name:       m["name"].(string),
		TypeId:     typeId,
		Author:     m["author"].(string),
		Public:     m["public"].(string),
		PublicDate: m["public_date"].(string),
		Stock:      stock,
	}
}
