package models

import (
	"fmt"
	"log"
	"strconv"

	"github.com/cn-lxy/bms_go/utils/dbutil"
)

const (
	User_Exist = iota
	User_RegisterError
	User_RegisterSuccess
	User_LoginError
	User_LoginSUccess
)

type User struct {
	Id       uint   // 用户id
	Name     string `json:"name"`     // 用户名
	Account  string `json:"account"`  // 账号/学号
	Password string `json:"password"` // 密码
	Sex      string `json:"sex"`      // 性别
	College  string `json:"college"`  // 学院
	Birthday string `json:"birthday"` // 出生日期
}

// 用户注册
func UserRegister(user *User) uint {
	// 验证账号未被注册
	querySql := fmt.Sprintf("select * from `users` where account = '%s'", user.Account)
	queryRes := dbutil.Query(querySql)
	if len(queryRes) != 0 {
		return User_Exist
	}

	sql := fmt.Sprintf(
		"insert into `users` (name, account, password, sex, college, birthday, register) values ('%s', '%s', '%s', '%s', '%s', '%s', now())",
		user.Name,
		user.Account,
		user.Password,
		user.Sex,
		user.College,
		user.Birthday)
	if err := dbutil.Update(sql); err != nil {
		log.Fatal(err.Error())
		return User_RegisterError
	}
	return User_RegisterSuccess
}

// 用户登录
func UserLogin(account string, password string) *User {
	queryUserSql := fmt.Sprintf(
		"select * from `users` where account = '%s' and password = '%s'",
		account, password)
	results := dbutil.Query(queryUserSql)
	if len(results) == 0 {
		return nil
	}

	id, _ := strconv.Atoi(results[0]["id"].(string))
	return &User{
		Id:       uint(id),
		Name:     results[0]["name"].(string),
		Password: results[0]["password"].(string),
		Sex:      results[0]["sex"].(string),
		College:  results[0]["college"].(string),
		Birthday: results[0]["birthday"].(string),
	}
}

// 更改用户信息：通过id
// 删除用户：通过id
