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
		return fmt.Errorf("account or password error")
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

// AddUser 添加新用户
func (a *Admin) AddUser(user *User) error {
	um := userManager{}
	return um.addUser(user)
}

// GetUserBySex 根据性别获取用户
// `limit`: 数据条数
// `offset`: 偏移量
func (a *Admin) SearchUser(by int64, key string, limit, offset int) ([]User, error) {
	um := userManager{}
	return um.searchUser(by, key, limit, offset)
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
	return bm.GetBooks(limit, offset, by, value...)
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
func (a *Admin) GetBorrowInfo(_type uint64, limit, offset int) ([]BorrowInfo, error) {
	bm := BorrowManager{}
	return bm.GetBorrows(_type, limit, offset)
}

// 获取用户借阅信息
func (a *Admin) GetUserBorrows(userName string, limit, offset int) ([]BorrowInfo, error) {
	bm := BorrowManager{}
	return bm.GetUserAllBorrowByName(userName, limit, offset)
}

// 获取用户未归还借阅
func (a *Admin) GetUserNotBackBorrow(userName string, limit, offset int) ([]BorrowInfo, error) {
	bm := BorrowManager{}
	return bm.GetUserNotBackBorrowByName(userName, limit, offset)
}

// 获取用户已归还借阅
func (a *Admin) GetUserBackedBorrow(userName string, limit, offset int) ([]BorrowInfo, error) {
	bm := BorrowManager{}
	return bm.GetUserBackedBorrowByName(userName, limit, offset)
}

// 获取书籍所有借阅
func (a *Admin) GetBookBorrowAll(bookName string, limit, offset int) ([]BorrowInfo, error) {
	bm := BorrowManager{}
	return bm.GetBookAllBorrow(bookName, limit, offset)
}

// 获取书籍未还借阅
func (a *Admin) GetBookNotBackBorrow(isbn string, limit, offset int) ([]BorrowInfo, error) {
	bm := BorrowManager{}
	return bm.GetBookNotBackBorrow(isbn, limit, offset)
}

// 获取书籍已还借阅
func (a *Admin) GetBookBackedBorrow(isbn string, limit, offset int) ([]BorrowInfo, error) {
	bm := BorrowManager{}
	return bm.GetBookBackedBorrow(isbn, limit, offset)
}

// 完成借阅
func (a *Admin) CompleteBorrow(bid uint64) error {
	bm := BorrowManager{}
	return bm.BackBook(bid)
}

// 获取所有用户 切片
func (um *userManager) getAllUser(limit, offset int) ([]User, error) {
	sql := "select * from `users` limit ? offset ?"
	args := []any{limit, offset}
	res := dbutil.Query(sql, args...)
	return mapsToUsers(res)
}

// 添加新用户
func (um *userManager) addUser(user *User) error {
	// 验证空字段
	if user.Name == "" || user.Account == "" || user.Password == "" || user.Sex == "" || user.College == "" || user.Birthday == "" {
		return fmt.Errorf("exist void field")
	}
	// 验证账号是否已经被注册
	sql_1 := fmt.Sprintf("select * from users where account = %s", user.Account)
	res := dbutil.Query(sql_1)
	if len(res) != 0 {
		return fmt.Errorf("this account of `%s` already existed", user.Account)
	}
	// 验证通过，插入新用户
	sql_2 := "insert into users(name, account, password, sex, college, birthday, register) values (?, ?, ?, ?, ?, ?, now())"
	args := []any{
		user.Name,
		user.Account,
		user.Password,
		user.Sex,
		user.College,
		user.Birthday,
	}
	return dbutil.Update(sql_2, args...)
}

// 搜索用户：by: { 0:姓名, 1:学号, 2: 性别, 3: uid}
func (um *userManager) searchUser(by int64, key string, limit, offset int) ([]User, error) {
	var sql string
	switch by {
	case 0:
		sql = "select * from `users` where name = ? limit ? offset ?"
	case 1:
		sql = "select * from `users` where account = ? limit ? offset ?"
	case 2:
		sql = "select * from `users` where sex = ? limit ? offset ?"
	case 3:
		sql = "select * from `users` where id = ? limit ? offset ?"
	default:
		return nil, fmt.Errorf("`by` params error")
	}

	args := []any{key, limit, offset}
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
