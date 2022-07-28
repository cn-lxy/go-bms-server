package controllers

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/cn-lxy/bms_go/models"
	"github.com/cn-lxy/bms_go/utils"
	"github.com/cn-lxy/bms_go/utils/randstring"
	"github.com/gofiber/fiber/v2"
)

const randomStringLen = 10

// 更新 cookie 有效时间
func updateCookie(c *fiber.Ctx, cookieValue string) {
	cookie := fiber.Cookie{
		Name:    "userFlag",
		Value:   cookieValue,
		Expires: time.Now().Add(24 * time.Hour),
		Path:    "/user/api/",
	}
	//set cookie
	c.Cookie(&cookie)
}

// 验证 cookie，当验证成功时更新cookie有效时间
func verifyCookie(c *fiber.Ctx, cookie string) error {
	session, err := store.Get(c)
	if err != nil {
		return fmt.Errorf("%s", "web server error!")
	}
	if session.Get(cookie) == nil {
		return fmt.Errorf("%s", "Authentication field!")
	}

	return nil
}

// 用户注册 handler ["/user/api/register", "POST"] ["tested"]
func userRegisterHandler(c *fiber.Ctx) error {
	user := models.User{}
	if err := c.BodyParser(&user); err != nil {
		utils.PrintfColorStr(utils.Red, "post parse error!")
	}

	if err := user.Register(); err != nil {
		log.Println(err.Error())
		return c.JSON(fiber.Map{
			"msg": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"msg": "register success",
	})
}

// 用户登录 handler ["POST", "user/api/login"] ["tested"]
func userLoginHandler(c *fiber.Ctx) error {
	user := models.User{}
	if err := c.BodyParser(&user); err != nil {
		utils.PrintfColorStr(utils.Red, "post parse error!")
	}

	if err := user.Login(); err != nil {
		return c.JSON(fiber.Map{
			"msg": "login field!",
		})
	}

	session, err := store.Get(c)
	if err != nil {
		return c.JSON(fiber.Map{
			"msg": "web server error!",
		})
	}

	// 随机字符作为cookie
	rs := randstring.RandString(randomStringLen)
	for session.Get(rs) != nil {
		rs = randstring.RandString(randomStringLen)
	}
	session.Set(rs, "1")
	if err := session.Save(); err != nil {
		fmt.Println("error: ", err.Error())
		return c.JSON(fiber.Map{
			"msg": "web server error!",
		})
	}

	// 更新并设置 cookie
	updateCookie(c, rs)

	return c.JSON(fiber.Map{
		"msg": "login success",
	})
}

// Get方法登录 ["GET", "/user/api/login"] ["tested"]
func userLoginForGet(c *fiber.Ctx) error {
	userFlag := c.Cookies("userFlag", "")
	if userFlag == "" {
		return c.JSON(fiber.Map{
			"msg": "login field!",
		})
	}

	if err := verifyCookie(c, userFlag); err != nil {
		return c.JSON(fiber.Map{
			"msg": err.Error(),
		})
	}

	updateCookie(c, userFlag)

	return c.JSON(fiber.Map{
		"msg": "login success",
	})
}

// 用户更改个人信息 handler ["/user/api/modify", "POST"] ["tested"]
func userModifyInfoHandler(c *fiber.Ctx) error {
	userFlag := c.Cookies("userFlag", "")
	if userFlag == "" {
		return c.JSON(fiber.Map{
			"msg": "No permission! Please log in.",
		})
	}

	if err := verifyCookie(c, userFlag); err != nil {
		return c.JSON(fiber.Map{
			"msg": err.Error(),
		})
	}

	user := models.User{}
	if err := c.BodyParser(&user); err != nil {
		utils.PrintfColorStr(utils.Red, "post parse error!")
		return err
	}
	if err := user.Modify(); err != nil {
		utils.PrintfColorStr(utils.Red, "user modify error!")
		return c.JSON(fiber.Map{
			"msg": err.Error(),
		})
	}

	updateCookie(c, userFlag)

	return c.JSON(fiber.Map{
		"msg": "update success",
	})
}

// 用户借阅 handler ["POST", "user/api/borrow" POST] ["tested"]
func userBorrowBookHandler(c *fiber.Ctx) error {
	userFlag := c.Cookies("userFlag", "")
	if userFlag == "" {
		return c.JSON(fiber.Map{
			"msg": "No permission! Please log in.",
		})
	}

	if err := verifyCookie(c, userFlag); err != nil {
		return c.JSON(fiber.Map{
			"msg": err.Error(),
		})
	}

	bi := models.BorrowInfo{}
	bm := models.BorrowManager{}

	if err := c.BodyParser(&bi); err != nil {
		utils.PrintfColorStr(utils.Yellow, "[userBorrowBookHandler]Parser Body err")
		return c.JSON(fiber.Map{
			"msg": err.Error(),
		})
	}

	if bi.UserId == 0 || bi.BookISBN == "" || bi.Days == 0 {
		return c.JSON(fiber.Map{
			"msg": "borrow info exist void variable.",
		})
	}
	if err := bm.BorrowBook(bi.UserId, bi.BookISBN, bi.Days); err != nil {
		return c.JSON(fiber.Map{
			"msg": "borrow book error.",
		})
	}
	return c.JSON(fiber.Map{
		"msg": "borrow success!",
	})
}

// 用户还书 handler ["POST", "user/api/back"] ["tested"]
func userBackBookHandler(c *fiber.Ctx) error {
	userFlag := c.Cookies("userFlag", "")
	if userFlag == "" {
		return c.JSON(fiber.Map{
			"msg": "No permission! Please log in.",
		})
	}

	if err := verifyCookie(c, userFlag); err != nil {
		return c.JSON(fiber.Map{
			"msg": err.Error(),
		})
	}

	bi := models.BorrowInfo{}
	bm := models.BorrowManager{}

	if err := c.BodyParser(&bi); err != nil {
		utils.PrintfColorStr(utils.Yellow, "[userBorrowBookHandler]Parser Body err")
		return c.JSON(fiber.Map{
			"msg": err.Error(),
		})
	}

	if bi.Id == 0 {
		return c.SendString(`{"msg":"borrow info exist void variable."}`)
	}
	if err := bm.BackBook(bi.Id); err != nil {
		return c.SendString(`{"msg":"back book error."}`)
	}
	return c.SendString(`{"msg":"back book success!"}`)
}

// 用户获取所有借阅信息 ["GET", "/user/api/borrowed?id=[userId]&limit&offset"] ["tested"]
func userGetAllBorrow(c *fiber.Ctx) error {
	userFlag := c.Cookies("userFlag", "")
	if userFlag == "" {
		return c.JSON(fiber.Map{
			"msg": "No permission! Please log in.",
		})
	}

	if err := verifyCookie(c, userFlag); err != nil {
		return c.JSON(fiber.Map{
			"msg": err.Error(),
		})
	}

	id, err := strconv.ParseUint(c.Query("id", "0"), 0, 64)
	if err != nil {
		log.Println(err.Error())
		return c.SendString(`{"msg": "string to uint64 error"}`)
	}
	if id == 0 {
		return c.SendString(`{"msg":"parameters exit void variable."}`)
	}
	limit, err := strconv.Atoi(c.Query("limit", "0"))
	if err != nil {
		log.Println(err.Error())
		c.SendString(`{"msg": "string to int error"}`)
	}
	offset, err := strconv.Atoi(c.Query("offset", "0"))
	if err != nil {
		log.Println(err.Error())
		c.SendString(`{"msg": "string to int error"}`)
	}
	bm := models.BorrowManager{}
	bs, err := bm.GetUserAllBorrow(id, limit, offset)
	if err != nil {
		log.Println(err.Error())
		return c.SendString(`{"msg":"Get user all borrow error!"}`)
	}
	return c.JSON(bs)
}

// 获取用户未还借用 ["GET", "/user/api/not_back?id=[userId]&limit&offset"]
func userGetNotBackBorrow(c *fiber.Ctx) error {
	userFlag := c.Cookies("userFlag", "")
	if userFlag == "" {
		return c.JSON(fiber.Map{
			"msg": "No permission! Please log in.",
		})
	}

	if err := verifyCookie(c, userFlag); err != nil {
		return c.JSON(fiber.Map{
			"msg": err.Error(),
		})
	}

	id, _ := strconv.ParseUint(c.Query("id", ""), 0, 64)
	limit, _ := strconv.Atoi(c.Query("limit", "0"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	if id == 0 {
		return c.SendString(`{"msg":"parameters exit void variable."}`)
	}
	bm := models.BorrowManager{}
	bs, err := bm.GetUserNotBackBorrow(id, limit, offset)
	if err != nil {
		return c.SendString(`{"msg":"Get user all borrow error!"}`)
	}
	return c.JSON(bs)
}

// 获取用户已还借阅 ["GET", "/user/api/backed?id=[userId]&limit&offset"]
func userGetBackedBorrow(c *fiber.Ctx) error {
	userFlag := c.Cookies("userFlag", "")
	if userFlag == "" {
		return c.JSON(fiber.Map{
			"msg": "No permission! Please log in.",
		})
	}

	if err := verifyCookie(c, userFlag); err != nil {
		return c.JSON(fiber.Map{
			"msg": err.Error(),
		})
	}

	id, err := strconv.ParseUint(c.Query("id", ""), 0, 64)
	if err != nil {
		return c.SendString(`{"msg": "id: string to uint64 error"}`)
	}
	if id == 0 {
		return c.SendString(`{"msg":"parameters exit void variable."}`)
	}
	limit, err := strconv.Atoi(c.Query("limit", "0"))
	if err != nil {
		return c.SendString(`{"msg": "limit: string to int error"}`)
	}
	offset, err := strconv.Atoi(c.Query("offset", "0"))
	if err != nil {
		return c.SendString(`{"msg": "limit: string to int error"}`)
	}
	bm := models.BorrowManager{}
	bs, err := bm.GetUserBackedBorrow(id, limit, offset)
	if err != nil {
		return c.SendString(`{"msg":"Get user all borrow error!"}`)
	}
	return c.JSON(bs)
}
