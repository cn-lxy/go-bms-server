package models

import (
	"fmt"

	"github.com/cn-lxy/bms_go/utils/dbutil"
)

type Admin struct {
	Id       uint64
	Account  string `json:"account"`
	Password string `json:"password"`
}

type userManager struct{}

// 管理员登录
func (a *Admin) Login() error {
	if err := a.verifyAP(); err != nil {
		return err
	}
	return nil
}

// 验证管理员账号密码
func (a *Admin) verifyAP() error {
	sql := "select * from `admins` where account = ? and password = ?"
	sqlArgs := make([]any, 0)
	sqlArgs = append(sqlArgs, a.Account, a.Password)
	res := dbutil.Query(sql, sqlArgs...)
	if len(res) == 0 {
		return fmt.Errorf("account or password error!")
	}
	return nil
}

// GetAllUser 获取所有用户
// `limit`: 数据条数
// `offset`: 偏移量
func (a *Admin) GetAllUser(limit, offset int) ([]User, error) {
	um := userManager{}
	return um.getAllUser(limit, offset)
}

// GetUserBySex 根据性别获取用户
// `limit`: 数据条数
// `offset`: 偏移量
func (a *Admin) GetUserBySex(sex string, limit, offset int) ([]User, error) {
	um := userManager{}
	return um.getUserBySex(sex, limit, offset)
}

// ModifyUser 更新用户信息
// `user`: 更新后的用户结构体，`Id`字段必须有效
func (a *Admin) ModifyUser(user *User) error {
	um := userManager{}
	return um.modifyUserInfo(user)
}

// DeleteUser 根据 `userId` 删除用户
// `userId`: 用户id
func (a *Admin) DeleteUser(userId uint64) error {
	um := userManager{}
	return um.deleteUserById(userId)
}

// 获取多个图书
// @Params `by: int` ["By_BookDefault", "By_BookTypeId"]
// example:
// bm.GetBooks(limit, offset, By_BookDefault, value[0])
// bm.GetBooks(limit, offset, By_BookTypeId, value[0])
func (a *Admin) GetBooks(limit, offset, by int, value ...any) ([]Book, error) {
	bm := BookManger{}
	return bm.GetBooks(limit, offset, by, value[0])
}

// 搜索图书
// `by: int` => By_BookISBN | By_BookName
func (a *Admin) SearchBook(by int, value ...string) (*Book, error) {
	bm := BookManger{}
	return bm.GetBook(by, value[0])
}

// 新建图书
func (a *Admin) InsertBook(b *Book) error {
	bm := BookManger{}
	return bm.InsertBook(b)
}

// 删除书籍
func (a *Admin) DeleteBook(isbn string) error {
	bm := BookManger{}
	return bm.DeleteBook(isbn)
}

// 更新书籍信息
func (a *Admin) UpdateBook(b *Book) error {
	bm := BookManger{}
	return bm.UpdateBook(b)
}

// 获取借阅记录信息
func (a *Admin) GetBorrowInfo(limit, offset int) ([]BorrowInfo, error) {
	bm := BorrowManager{}
	return bm.GetBorrows(limit, offset)
}

// 获取用户借阅信息
func (a *Admin) GetUserBorrows(userId uint64, limit, offset int) ([]BorrowInfo, error) {
	bm := BorrowManager{}
	return bm.GetUserAllBorrow(userId, limit, offset)
}

// 获取用户未归还借阅
func (a *Admin) GetUserNotBackBorrow(userId uint64, limit, offset int) ([]BorrowInfo, error) {
	bm := BorrowManager{}
	return bm.GetUserNotBackBorrow(userId, limit, offset)
}

// 获取用户已归还借阅
func (a *Admin) GetUserBackedBorrow(userId uint64, limit, offset int) ([]BorrowInfo, error) {
	bm := BorrowManager{}
	return bm.GetUserBackedBorrow(userId, limit, offset)
}

// 获取被借书籍用户
func (a *Admin) GetBookBorrowedUser(isbn string) ([]BorrowInfo, error) {
	bm := BorrowManager{}
	return bm.GetBookAllBorrow(isbn)
}

// 获取所有用户 切片
func (um *userManager) getAllUser(limit, offset int) ([]User, error) {
	sql := "select * from `users` limit ? offset ?"
	args := []any{limit, offset}
	res := dbutil.Query(sql, args...)
	return mapsToUsers(res)
}

// 根据性别获取用户
func (um *userManager) getUserBySex(sex string, limit, offset int) ([]User, error) {
	sql := "select * from `user` where sex = ? limit ? offset ?"
	args := []any{sex, limit, offset}
	res := dbutil.Query(sql, args...)
	return mapsToUsers(res)
}

// 更改用户信息
func (um *userManager) modifyUserInfo(user *User) error {
	return user.Modify()
}

// 删除用户
func (um *userManager) deleteUserById(userId uint64) error {
	user := User{
		Id: userId,
	}
	return user.Delete()
}
