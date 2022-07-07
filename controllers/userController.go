package controllers

import (
	"fmt"
	"strconv"

	"github.com/cn-lxy/bms_go/models"
	"github.com/cn-lxy/bms_go/utils"
	"github.com/gofiber/fiber/v2"
)

var UserController *fiber.App

func init() {
	UserController = fiber.New()
	UserController.Post("/login", userLoginHandler)
	UserController.Post("/register", userRegisterHandler)
	UserController.Post("/modify", userModifyInfoHandler)
	UserController.Post("/borrow", userBorrowBookHandler)
	UserController.Post("/back", userBackBookHandler)
	UserController.Get("/borrowed", userGetAllBorrow)
	UserController.Get("/not_back", userGetNotBackBorrow)
	UserController.Get("/backed", userGetBackedBorrow)
}

// 用户注册 handler ["/user/api/register", "POST"] ["tested"]
func userRegisterHandler(c *fiber.Ctx) error {
	// user := &models.User{}
	user := models.User{}
	if err := c.BodyParser(&user); err != nil {
		utils.PrintfColorStr(utils.Red, "post parse error!")
	}

	var msg string
	if err := user.Register(); err != nil {
		msg = err.Error()
	} else {
		msg = "注册成功"
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
	utils.PrintfColorStr(utils.Green, fmt.Sprintf("`msg`: %s, account: %s", msg, user.Account))
	return c.SendString(userJson)
}

// 用户登录 handler ["POST", "user/api/login"] ["tested"]
func userLoginHandler(c *fiber.Ctx) error {
	user := models.User{}
	if err := c.BodyParser(&user); err != nil {
		utils.PrintfColorStr(utils.Red, "post parse error!")
	}
	utils.PrintfColorStr(utils.Green,
		fmt.Sprintf(`{"account": "%s", "password": "%s"}`, user.Account, user.Password))
	if err := user.Login(); err != nil {
		utils.PrintfColorStr(utils.Red, "登录失败")
		return c.SendString(`{"msg": "登录失败"}`)
	}
	msg := fmt.Sprintf(
		`{"msg": "登录成功!", "username": "%s"}`, user.Name)
	utils.PrintfColorStr(utils.Green, msg)
	return c.SendString(msg)
}

// 用户更改个人信息 handler ["/user/api/modify", "POST"] ["tested"]
func userModifyInfoHandler(c *fiber.Ctx) error {
	user := models.User{}
	if err := c.BodyParser(&user); err != nil {
		utils.PrintfColorStr(utils.Red, "post parse error!")
		return err
	}
	if err := user.Modify(); err != nil {
		utils.PrintfColorStr(utils.Red, "user modify error!")
		return err
	}
	return c.SendString("更新成功")
}

// 用户借阅 handler ["POST", "user/api/borrow" POST] ["tested"]
func userBorrowBookHandler(c *fiber.Ctx) error {
	bi := models.BorrowInfo{}
	bm := models.BorrowManager{}

	if err := c.BodyParser(&bi); err != nil {
		utils.PrintfColorStr(utils.Yellow, "[userBorrowBookHandler]Parser Body err")
	}

	if bi.UserId == 0 || bi.BookISBN == "" || bi.Days == 0 {
		return c.SendString(`{"msg":"borrow info exist void variable."}`)
	}
	if err := bm.BorrowBook(bi.UserId, bi.BookISBN, bi.Days); err != nil {
		return c.SendString(`{"msg":"borrow book error."}`)
	}
	return c.SendString(`{"msg":"borrow success!"}`)
}

// 用户还书 handler ["POST", "user/api/back"] ["tested"]
func userBackBookHandler(c *fiber.Ctx) error {
	bi := models.BorrowInfo{}
	bm := models.BorrowManager{}

	if err := c.BodyParser(&bi); err != nil {
		// log.Println(err.Error())
		utils.PrintfColorStr(utils.Yellow, "[userBorrowBookHandler]Parser Body err")
	}

	if bi.Id == 0 {
		return c.SendString(`{"msg":"borrow info exist void variable."}`)
	}
	if err := bm.BackBook(bi.Id); err != nil {
		return c.SendString(`{"msg":"back book error."}`)
	}
	return c.SendString(`{"msg":"back book success!"}`)
}

// 用户获取所有借阅信息 ["GET", "/user/api/borrowed?id=[userId]"] ["tested"]
func userGetAllBorrow(c *fiber.Ctx) error {
	id, _ := strconv.ParseUint(c.Query("id", ""), 0, 64)
	if id == 0 {
		return c.SendString(`{"msg":"parameters exit void variable."}`)
	}
	bm := models.BorrowManager{}
	bs, err := bm.GetUserAllBorrow(id)
	if err != nil {
		return c.SendString(`{"msg":"Get user all borrow error!"}`)
	}
	return c.JSON(bs)
}

// 获取用户未还借用 ["GET", "/user/api/not_back?id=[userId]"]
func userGetNotBackBorrow(c *fiber.Ctx) error {
	id, _ := strconv.ParseUint(c.Query("id", ""), 0, 64)
	if id == 0 {
		return c.SendString(`{"msg":"parameters exit void variable."}`)
	}
	bm := models.BorrowManager{}
	bs, err := bm.GetUserNotBackBorrow(id)
	if err != nil {
		return c.SendString(`{"msg":"Get user all borrow error!"}`)
	}
	return c.JSON(bs)
}

// 获取用户已还借阅 ["GET", "/user/api/backed?id=[userId]"]
func userGetBackedBorrow(c *fiber.Ctx) error {
	id, _ := strconv.ParseUint(c.Query("id", ""), 0, 64)
	if id == 0 {
		return c.SendString(`{"msg":"parameters exit void variable."}`)
	}
	bm := models.BorrowManager{}
	bs, err := bm.GetUserBackedBorrow(id)
	if err != nil {
		return c.SendString(`{"msg":"Get user all borrow error!"}`)
	}
	return c.JSON(bs)
}
