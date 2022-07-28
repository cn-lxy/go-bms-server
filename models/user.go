package models

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/cn-lxy/bms_go/utils/dbutil"
)

const (
	User_Exist = iota
	User_RegisterError
	User_RegisterSuccess
	User_LoginError
	User_LoginSUccess
	By_UserName
	By_UserAccount
	By_UserSex
	By_UserDefault
)

type User struct {
	Id       uint64 `json:"id"`       // 用户id
	Name     string `json:"name"`     // 用户名
	Account  string `json:"account"`  // 账号/学号
	Password string `json:"password"` // 密码
	Sex      string `json:"sex"`      // 性别
	College  string `json:"college"`  // 学院
	Birthday string `json:"birthday"` // 出生日期
}

// 用户注册
func (u *User) Register() error {
	// 验证 `User` 结构体是否为空
	if u.Name == "" || u.Account == "" || u.Password == "" || u.Sex == "" || u.College == "" || u.Birthday == "" {
		return fmt.Errorf("user struct exist nil.")
	}

	// 验证账号未被注册
	if verifyAccountExisted(u.Id, u.Account) {
		return fmt.Errorf("this account : `%s` has been registered.", u.Account)
	}

	// 插入新数据
	insertSql := "insert into `users` (name, account, password, sex, college, birthday, register) values (?, ?, ?, ?, ?, ?, now())"
	insertArgs := make([]any, 0)
	insertArgs = append(insertArgs, u.Name, u.Account, u.Password, u.Sex, u.College, u.Birthday)
	if err := dbutil.Update(insertSql, insertArgs...); err != nil {
		log.Fatal(err.Error())
		return err
	}
	return nil
}

// 用户登录
func (u *User) Login() error {
	// 检查账号密码不为空
	if u.Account == "" || u.Password == "" {
		return fmt.Errorf("user of account or password is nil")
	}

	// 验证账号密码
	user, err := verifyAP(u.Account, u.Password)
	if err != nil {
		return err
	}

	// FIX: 赋值导致 `Id` 字段数据丢失
	// before fix: `u(*User) = user(*User)`
	// after fix: `*u = *user`
	*u = *user
	return nil
}

// 更改用户信息：通过id
func (u *User) Modify() error {
	if u.Id == 0 {
		return fmt.Errorf("err: %s", "id is nil")
	}

	// if _, err := verifyAP(u.Account, u.Password); err != nil {
	// 	return err
	// }

	fieldSlice := make([]string, 0, 10)
	if u.Account != "" {
		if verifyAccountExisted(u.Id, u.Account) {
			return fmt.Errorf("after modify `account`: `%s` is existed.", u.Account)
		}
		fieldSlice = append(fieldSlice, fmt.Sprintf("account = '%s'", u.Account))
	}
	if u.Name != "" {
		fieldSlice = append(fieldSlice, fmt.Sprintf("name = '%s'", u.Name))
	}

	if u.Password != "" {
		fieldSlice = append(fieldSlice, fmt.Sprintf("password = '%s'", u.Password))
	}
	if u.Sex != "" {
		fieldSlice = append(fieldSlice, fmt.Sprintf("sex = '%s'", u.Sex))
	}
	if u.College != "" {
		fieldSlice = append(fieldSlice, fmt.Sprintf("college = '%s'", u.College))
	}
	if u.Birthday != "" {
		fieldSlice = append(fieldSlice, fmt.Sprintf("birthday = '%s'", u.Birthday))
	}
	setField := strings.Join(fieldSlice, ", ")
	sql := fmt.Sprintf("update `users` set %s where id = %d", setField, u.Id)
	return dbutil.Update(sql)
}

// 删除用户：通过id
func (u *User) Delete() error {
	if u.Id == 0 {
		return fmt.Errorf("user's id is nil")
	}

	// TAG 验证设置不太合理 => 删除需要提供账号密码不方便管理员模块调用，但是不设置验证会导致用户模块非法调用
	// FIX: 尝试 在用户登录时设置 `session` 或者 `cookie`
	// if _, err := verifyAP(u.Account, u.Password); err != nil {
	// 	return err
	// }

	sql := "delete from `users` where id = ?"
	sqlArgs := make([]any, 0)
	sqlArgs = append(sqlArgs, u.Id)
	return dbutil.Update(sql, sqlArgs...)
}

// TAG 验证账号密码
func verifyAP(account, password string) (*User, error) {
	sql := "select * from `users` where account = ? and password = ?"
	sqlArgs := make([]any, 0)
	sqlArgs = append(sqlArgs, account, password)
	results := dbutil.Query(sql, sqlArgs...)
	if len(results) == 0 {
		return nil, fmt.Errorf("this verify of  account: `%s` is not pass! ", account)
	}
	return mapToUser(results[0])
}

// 验证用户是否存在
func verifyAccountExisted(uid uint64, account string) bool {
	sql := "select account from `users` where account = ? and id != ?"
	sqlArgs := make([]any, 0)
	sqlArgs = append(sqlArgs, account, uid)
	results := dbutil.Query(sql, sqlArgs...)
	if len(results) == 0 {
		return false // 不存在
	}
	return true // 存在
}

// 将 map 转换成 User 结构体
func mapToUser(m map[string]any) (*User, error) {
	id, err := strconv.ParseUint(m["id"].(string), 0, 64)
	if err != nil {
		return nil, fmt.Errorf("parse uint err: %s", err)
	}
	return &User{
		Id:       id,
		Name:     m["name"].(string),
		Account:  m["account"].(string),
		Password: m["password"].(string),
		Sex:      m["sex"].(string),
		College:  m["college"].(string),
		Birthday: m["birthday"].(string),
	}, nil
}

func mapsToUsers(ms []map[string]any) ([]User, error) {
	users := make([]User, 0, len(ms))
	for _, m := range ms {
		user, err := mapToUser(m)
		if err != nil {
			return nil, err
		}
		users = append(users, *user)
	}
	return users, nil
}

// 获取一个用户：通过 `name` | `account``
func GetUser(byType int, value string) (*User, error) {
	var sql string
	var typeField string
	switch byType {
	case By_UserName:
		typeField = "name"
		break
	case By_UserAccount:
		typeField = "account"
		break
	}
	sql = fmt.Sprintf("select * from `users` where %s = '%s'", typeField, value)
	result := dbutil.Query(sql)
	if len(result) == 0 {
		return nil, fmt.Errorf("query by `%s` err", typeField)
	}
	return mapToUser(result[0])
}

// 获取用户集，通过 `sex` | `default`（所有）
func GetUsers(limit, offset, by int, value ...string) ([]User, error) {
	var sql string
	switch by {
	case By_UserSex:
		sql = fmt.Sprintf("select * from `users` where sex = '%s' limit %d offset %d", value[0], limit, offset)
		break
	case By_UserDefault:
		sql = fmt.Sprintf("select * from `users` limit %d offset %d", limit, offset)
		break
	}
	results := dbutil.Query(sql)
	if len(results) == 0 {
		return nil, fmt.Errorf("get users error")
	}
	return mapsToUsers(results)
}
