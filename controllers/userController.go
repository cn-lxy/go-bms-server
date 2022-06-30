package controllers

import (
	"fmt"

	"github.com/cn-lxy/bms_go/models"
	"github.com/cn-lxy/bms_go/utils"
	"github.com/gofiber/fiber/v2"
)

var UserController *fiber.App

func init() {
	UserController = fiber.New()
	UserController.Post("/login", userLoginHandler)
	UserController.Post("/register", userRegisterHandler)
}

// 用户注册 web handler
func userRegisterHandler(c *fiber.Ctx) error {
	user := &models.User{}
	if err := c.BodyParser(&user); err != nil {
		utils.PrintfColorStr(utils.Red, "post parse error!")
	}

	status := models.UserRegister(user)
	var msg string
	switch status {
	case models.User_RegisterSuccess:
		msg = "注册成功"
		utils.PrintfColorStr(utils.Green, msg)
		break
	case models.User_Exist:
		msg = "账号已存在"
		utils.PrintfColorStr(utils.Yellow, msg)
		break
	case models.User_LoginError:
		msg = "注册失败"
		utils.PrintfColorStr(utils.Red, msg)
		break
	}

	userJson := fmt.Sprintf(
		`{"msg": "%s", "user":{"name":"%s", "account":"%s", "password":"%s", "sex": "%s", "college":"%s", "birthday":"%s"}}`,
		msg,
		user.Name,
		user.Account,
		user.Password,
		user.Sex,
		user.College,
		user.Birthday)
	return c.SendString(userJson)
}

// 用户登录 web handler
func userLoginHandler(c *fiber.Ctx) error {
	var ap struct {
		Account  string `json:"account"`
		Password string `json:"password"`
	}
	if err := c.BodyParser(&ap); err != nil {
		utils.PrintfColorStr(utils.Red, "post parse error!")
	}
	utils.PrintfColorStr(utils.Green,
		fmt.Sprintf(`{"account": "%s", "password": "%s"}`, ap.Account, ap.Password))
	user := models.UserLogin(ap.Account, ap.Password)
	if user == nil {
		utils.PrintfColorStr(utils.Red, "登录失败")
		return c.SendString(`{"msg": "登录失败"}`)
	}
	msg := fmt.Sprintf(
		`{"msg": "登录成功!", "username": "%s"}`, user.Name)
	utils.PrintfColorStr(utils.Green, msg)
	return c.SendString(msg)
}

// TODO
