package models

import (
	"fmt"
	"testing"
	"time"

	"github.com/cn-lxy/bms_go/utils"
)

func TestTimeNow(t *testing.T) {
	fmt.Println(time.Now().Format("2006-01"))
}

func TestBookInsert(t *testing.T) {
	db.AutoMigrate(&Book{})
	book := Book{
		Name:       "Go语言高级编程",
		Author:     "曹大",
		Public:     "清华大学出版社",
		PublicTime: time.Now().Format("2006-01"),
	}
	book.Insert(db)
}

func TestBookQuery(t *testing.T) {
	db.AutoMigrate(&Book{})
	t.Run("query_book_byName", func(t *testing.T) {
		book := Book{
			Name: "Go语言高级编程",
		}
		book.Query(db, BOOK_NAME, book.Name)
		log := fmt.Sprintf("{bookName: %s, bookAuthor: %s, bookPublic: %s, bookPublicTime: %s}\n",
			book.Name, book.Author, book.Public, book.PublicTime)
		utils.PrintfColorStr(utils.GREEN, log)

	})
	t.Run("query_book_byAuthor", func(t *testing.T) {
		book := Book{
			Author: "曹大",
		}
		book.Query(db, BOOK_AUTHOR, book.Author)
		log := fmt.Sprintf("{bookName: %s, bookAuthor: %s, bookPublic: %s, bookPublicTime: %s}\n",
			book.Name, book.Author, book.Public, book.PublicTime)
		utils.PrintfColorStr(utils.GREEN, log)
	})
}

func TestBookUpdate(t *testing.T) {
	db.AutoMigrate(&Book{})

	book := Book{}
	book.Query(db, BOOK_NAME, "Go语言高级编程")

	updateBefore := fmt.Sprintf("{bookName: %s, bookAuthor: %s, bookPublic: %s, bookPublicTime: %s}\n",
		book.Name, book.Author, book.Public, book.PublicTime)
	utils.PrintfColorStr(utils.GREEN, updateBefore)

	book.Author = "曹大3"
	book.Update(db, BOOK_NAME)

	updateAfter := fmt.Sprintf("{bookName: %s, bookAuthor: %s, bookPublic: %s, bookPublicTime: %s}\n",
		book.Name, book.Author, book.Public, book.PublicTime)
	utils.PrintfColorStr(utils.GREEN, updateAfter)
}

func TestInsertBook(t *testing.T) {
	book := Book{
		Name:       "Go语言高级编程",
		Author:     "曹大2",
		BookTypeId: 1,
		Public:     "清华大学出版社",
		PublicTime: time.Now().Format("2006-01"),
	}
	if err = InsertBook(&book); err != nil {
		println(err)
	}
}

func TestGetBook(t *testing.T) {
	book, err := GetBook(BOOK_NAME, "Go语言高级编程")
	if err != nil {
		t.Fatal(err)
	}
	utils.PrintfColorStr(utils.GREEN, fmt.Sprintf(
		"Book {name: %s, author: %s, public: %s, publicTime: %s}",
		book.Name, book.Author, book.Public, book.PublicTime))
}

func TestGetBookSlice(t *testing.T) {
	books, err := GetBookSlice(-1, 10, 0)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%#v\n", books)
}

func TestCreateBookTypes(t *testing.T) {
	if err = CreateBookType("小说"); err != nil {
		t.Fatal(err)
	}
}
