package models

import (
	"fmt"
	"strconv"

	"github.com/cn-lxy/bms_go/utils/dbutil"
)

const borrowTableName = "borrow"

// 借阅信息
type BorrowInfo struct {
	Id         uint64 `json:"id"`
	UserId     uint64 `json:"userId"`
	UserName   string `json:"userName"`
	BookISBN   string `json:"bookIsbn"`
	BookName   string `json:"bookName"`
	BookAuthor string `json:"bookAuthor"`
	Days       uint64 `json:"days"`
	BorrowDate string `json:"borrowDate"`
	BackDate   string `json:"backDate"`
}

type BorrowManager struct{}

// BorrowBook 通过 book-isbn 和 user-id 借阅书籍
func (b *BorrowManager) BorrowBook(userId uint64, bookIsbn string, days uint64) error {
	if userId == 0 || bookIsbn == "" || days == 0 {
		return fmt.Errorf("parameters exist nil")
	}
	sql := "insert into `borrow` (uid, bid, days, borrow_date, back_date) values (?, ?, ?, now(), null)"
	sqlArgs := make([]any, 0)
	sqlArgs = append(sqlArgs, userId, bookIsbn, days)
	return dbutil.Update(sql, sqlArgs...)
}

// BackBook 归还图书
func (b *BorrowManager) BackBook(borrowId uint64) error {
	sql := fmt.Sprintf("update `%s` set back_date = now() where id = ?", borrowTableName)
	sqlArgs := make([]any, 0)
	sqlArgs = append(sqlArgs, borrowId)
	return dbutil.Update(sql, sqlArgs...)
}

// GetUserAllBorrow 获取用户的所有借阅信息
func (b *BorrowManager) GetUserAllBorrow(uid uint64) ([]BorrowInfo, error) {
	sql := fmt.Sprintf(
		`select borrow.id, u.id user_id, u.name user_name, b.isbn, b.name book_name, b.author, days, borrow_date, back_date
		from %s join users u on u.id = borrow.uid join books b on b.isbn = borrow.bid
		where uid = ?`,
		borrowTableName)
	sqlArgs := make([]any, 0)
	sqlArgs = append(sqlArgs, uid)
	res := dbutil.Query(sql, sqlArgs...)
	if len(res) == 0 {
		return nil, nil
	}
	return mapsToBorrowInfos(res)
}

// 获取图书所有被借阅信息
func (b *BorrowManager) GetBookAllBorrow(bid string) ([]BorrowInfo, error) {
	sql := fmt.Sprintf(
		`select borrow.id, u.id user_id, u.name user_name, b.isbn, b.name book_name, b.author, days, borrow_date, back_date
		from %s join users u on u.id = borrow.uid join books b on b.isbn = borrow.bid
		where isbn = ?`,
		borrowTableName)
	sqlArgs := make([]any, 0)
	sqlArgs = append(sqlArgs, bid)
	res := dbutil.Query(sql, sqlArgs...)
	if len(res) == 0 {
		return nil, nil
	}
	return mapsToBorrowInfos(res)
}

// 获取用户未还借阅
func (b *BorrowManager) GetUserNotBackBorrow(userId uint64) ([]BorrowInfo, error) {
	sql := fmt.Sprintf(
		`select borrow.id, u.id user_id, u.name user_name, b.isbn, b.name book_name, b.author, days, borrow_date, back_date
		from %s join users u on u.id = borrow.uid join books b on b.isbn = borrow.bid
		where u.id = ? and back_date is null`,
		borrowTableName)
	sqlArgs := make([]any, 0)
	sqlArgs = append(sqlArgs, userId)
	res := dbutil.Query(sql, sqlArgs...)
	if len(res) == 0 {
		return nil, nil
	}
	return mapsToBorrowInfos(res)
}

// 获取用户已还还借阅
func (b *BorrowManager) GetUserBackedBorrow(userId uint64) ([]BorrowInfo, error) {
	sql := fmt.Sprintf(
		`select borrow.id, u.id user_id, u.name user_name, b.isbn, b.name book_name, b.author, days, borrow_date, back_date
		from %s join users u on u.id = borrow.uid join books b on b.isbn = borrow.bid
		where u.id = ? and back_date is not null`,
		borrowTableName)
	sqlArgs := make([]any, 0)
	sqlArgs = append(sqlArgs, userId)
	res := dbutil.Query(sql, sqlArgs...)
	if len(res) == 0 {
		return nil, nil
	}
	return mapsToBorrowInfos(res)
}

func mapToBorrowInfo(m map[string]any) (*BorrowInfo, error) {
	id, err := strconv.ParseUint(m["id"].(string), 0, 64)
	if err != nil {
		return nil, fmt.Errorf("parse uint to `id` error")
	}
	userId, err := strconv.ParseUint(m["user_id"].(string), 0, 64)
	if err != nil {
		return nil, fmt.Errorf("parse uint to `user_id` error")
	}
	days, err := strconv.ParseUint(m["days"].(string), 0, 64)
	if err != nil {
		return nil, fmt.Errorf("parse uint to `days` error")
	}

	return &BorrowInfo{
		Id:         id,
		UserId:     userId,
		UserName:   m["user_name"].(string),
		BookISBN:   m["isbn"].(string),
		BookName:   m["book_name"].(string),
		BookAuthor: m["author"].(string),
		Days:       days,
		BorrowDate: m["borrow_date"].(string),
		BackDate:   m["back_date"].(string),
	}, nil
}

func mapsToBorrowInfos(ms []map[string]any) ([]BorrowInfo, error) {
	bis := make([]BorrowInfo, 0, len(ms))
	for _, m := range ms {
		bi, err := mapToBorrowInfo(m)
		if err != nil {
			return nil, fmt.Errorf("[mapsToBorrowInfos] change error")
		}
		bis = append(bis, *bi)
	}
	return bis, nil
}
