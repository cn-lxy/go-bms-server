package models

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/cn-lxy/bms_go/utils"
)

var user User

func init() {
	user = User{
		Name:     "aa",
		Account:  "2019001",
		Password: "123456",
		Sex:      "男",
		College:  "信息与通信工程学院",
		Birthday: "2019-01-01",
	}
}

func TestUserRegister(t *testing.T) {
	if err := user.Register(); err != nil {
		t.Fatal(err)
	}
	utils.PrintfColorStr(utils.Green, "注册成功")
}

func TestUserLogin(t *testing.T) {
	if err := user.Login(); err != nil {
		t.Fatal(err)
	}
	utils.PrintfColorStr(utils.Green, "登录成功, userName: "+user.Name)
	// fmt.Printf("after login, the id is `%d`\n", user.Id)
	fmt.Printf("after login user is:\n%v\n", user)
}

func TestModifyUser(t *testing.T) {
	user.Name = "bb"
	if err := user.Modify(); err != nil {
		t.Fatal(err)
	}
	utils.PrintfColorStr(utils.Green, "修改成功")
}

func TestDeleteUser(t *testing.T) {
	if err := user.Delete(); err != nil {
		t.Fatal(err)
	}
	utils.PrintfColorStr(utils.Green, "删除成功")
}

func TestGetUser(t *testing.T) {
	// by user-name
	t.Run("GetUser - by name", func(t *testing.T) {
		user1, err := GetUser(By_UserName, "小明")
		if err != nil {
			t.Fatal(err)
		}
		fmt.Printf("%v\n", user1)
	})

	// by user-account
	t.Run("GetUser - by account", func(t *testing.T) {
		user2, err := GetUser(By_UserAccount, "2019136417")
		if err != nil {
			t.Fatal(err)
		}
		fmt.Printf("%v\n", user2)
	})
}

func TestGetUsers(t *testing.T) {
	t.Run("GetUsers - by sex", func(t *testing.T) {
		users, err := GetUsers(10, 0, By_UserSex, "男")
		if err != nil {
			t.Fatal(err)
		}
		for i, u := range users {
			fmt.Printf("%d: %v\n", i+1, u)
		}
	})
	t.Run("GetUsers - by default", func(t *testing.T) {
		users, err := GetUsers(10, 0, By_UserDefault)
		if err != nil {
			t.Fatal(err)
		}
		for i, u := range users {
			fmt.Printf("%d: %v\n", i+1, u)
		}
	})
}

func TestStringToUint(t *testing.T) {
	id := "24"
	uint_id, err := strconv.ParseUint(id, 0, 64)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("uint id: %v, type: %T\n", uint_id, uint_id)
}
