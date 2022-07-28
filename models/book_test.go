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
	bm := BookManger{}
	if err := bm.InsertBook(&book); err != nil {
		t.Fatal(err)
	}
	utils.PrintfColorStr(utils.Green, "添加书籍成功")
}

func TestGetBook(t *testing.T) {
	byName := "test书籍"
	bm := BookManger{}
	book, err := bm.GetBook(By_BookName, byName)
	if err != nil {
		t.Fatal(err)
	}
	utils.PrintfColorStr(utils.Green, fmt.Sprintf("%v", book))
}

func TestGetBooks(t *testing.T) {
	bm := BookManger{}
	typeId := 3
	books, err := bm.GetBooks(10, 0, By_BookTypeId, typeId)
	if err != nil {
		t.Fatal(err)
	}
	for _, book := range books {
		fmt.Printf("%v\n", book)
	}
	utils.PrintfColorStr(utils.Green, "books of len: "+strconv.Itoa(len(books)))
}

func TestDeleteBook(t *testing.T) {
	bm := BookManger{}

	book, err := bm.GetBook(By_BookName, "test书籍")
	if err != nil {
		t.Fatal(err)
	}
	if err := bm.DeleteBook(book.Isbn); err != nil {
		t.Fatal(err)
	}
	utils.PrintfColorStr(utils.Green, "删除成功")
	utils.PrintfColorStr(utils.Green, book.Isbn)
}
