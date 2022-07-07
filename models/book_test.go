package models

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/cn-lxy/bms_go/utils"
)

func TestBookRegister(t *testing.T) {
	book := Book{
		Isbn:       "12345678",
		Name:       "test书籍",
		TypeId:     3,
		Author:     "test",
		Public:     "test出版社",
		PublicDate: "1991-01-01",
		Stock:      100,
	}
	status := BookRegister(&book)
	switch status {
	case Book_RegisterSuccess:
		utils.PrintfColorStr(utils.Green, "添加书籍成功")
		break
	case Book_RegisterError:
		utils.PrintfColorStr(utils.Red, "添加书籍失败")
		panic("添加书籍失败")
	case Book_RegisterExisted:
		utils.PrintfColorStr(utils.Yellow, "书籍已存在")
		panic("书籍已存在")
	}
}

func TestGetBook(t *testing.T) {
	byName := "test书籍"
	book, err := GetBook(By_BookName, byName)
	if err != nil {
		t.Fatal(err)
	}
	utils.PrintfColorStr(utils.Green, fmt.Sprintf("%v", book))
}

func TestGetBooks(t *testing.T) {
	typeId := 3
	books := GetBooks(10, 0, By_BookTypeId, typeId)
	for _, book := range books {
		fmt.Printf("%v\n", book)
	}
	utils.PrintfColorStr(utils.Green, "books of len: "+strconv.Itoa(len(books)))
}

func TestDeleteBook(t *testing.T) {
	book, err := GetBook(By_BookName, "test书籍")
	if err != nil {
		t.Fatal(err)
	}
	if err := DeleteBook(book.Isbn); err != nil {
		t.Fatal(err)
	}
	utils.PrintfColorStr(utils.Green, "删除成功")
	utils.PrintfColorStr(utils.Green, book.Isbn)
}
