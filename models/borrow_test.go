package models

import (
	"fmt"
	"testing"

	"github.com/cn-lxy/bms_go/utils"
)

func TestBorrowBook(t *testing.T) {
	bm := BorrowManager{}
	if err := bm.BorrowBook(1, "9787010009223", 10); err != nil {
		t.Fatal(err)
	}
	utils.PrintfColorStr(utils.Green, "借阅成功")
}

func TestBackBook(t *testing.T) {
	bm := BorrowManager{}
	if err := bm.BackBook(3); err != nil {
		t.Fatal(err)
	}
	utils.PrintfColorStr(utils.Green, "归还成功")
}

func TestGetUserAllBorrow(t *testing.T) {
	bm := BorrowManager{}
	borrows, err := bm.GetUserAllBorrow(1, 10, 0)
	if err != nil {
		t.Fatal(err)
	}
	for i, b := range borrows {
		fmt.Printf("%d: %v\n", i, b)
	}
}

func TestGetBookAllBorrow(t *testing.T) {
	bm := BorrowManager{}
	borrows, err := bm.GetBookAllBorrow("9787010009223")
	if err != nil {
		t.Fatal(err)
	}
	for i, b := range borrows {
		fmt.Printf("%d: %v\n", i, b)
	}
}

func TestGetSUserNotBackBorrow(t *testing.T) {
	var userId uint64 = 1
	bm := BorrowManager{}
	bs, err := bm.GetUserNotBackBorrow(userId, 10, 0)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(len(bs))
	for _, b := range bs {
		fmt.Printf("%v\n", b)
	}
}

func TestGetSUserBackedBorrow(t *testing.T) {
	var userId uint64 = 1
	bm := BorrowManager{}
	bs, err := bm.GetUserBackedBorrow(userId, 10, 0)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(len(bs))
	for _, b := range bs {
		fmt.Printf("%v\n", b)
	}
}
