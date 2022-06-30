package models

import (
	"testing"

	"github.com/cn-lxy/bms_go/utils"
)

func TestUserRegister(t *testing.T) {
	user := &User{
		Name:     "aa",
		Account:  "2019",
		Password: "123456",
		Sex:      "男",
		College:  "信息与通信工程学院",
		Birthday: "2019-01-01",
	}
	status := UserRegister(user)
	if status == User_RegisterSuccess {
		utils.PrintfColorStr(utils.Green, "注册成功")
	} else if status == User_Exist {
		utils.PrintfColorStr(utils.Yellow, "用户已存在")
	} else if status == User_RegisterError {
		panic("注册失败")
	}
}

func TestUserLogin(t *testing.T) {
	user := UserLogin("2019", "123456")
	if user != nil {
		utils.PrintfColorStr(utils.Green, "登录成功, userName: "+user.Name)
	} else {
		panic("登录失败!")
	}
}
