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

/*
`code`: 1 -> 成功， 0: -> 错误，-1 -> 失败
*/

const randomStringLen = 10

// 更新 cookie 有效时间
func updateCookie(c *fiber.Ctx, cookieValue string) {
	cookie := fiber.Cookie{
		Name:    "userFlag",
		Value:   cookieValue,
		Expires: time.Now().Add(24 * time.Hour),
		Path:    "/",
	}
	//set cookie
	c.Cookie(&cookie)
}

// 验证 cookie，当验证成功时更新cookie有效时间
func verifyCookie(c *fiber.Ctx, cookie string) (any, error) {
	session, err := store.Get(c)
	if err != nil {
		return nil, fmt.Errorf("%s", "web server error!")
	}
	id := session.Get(cookie)
	if id == nil {
		return nil, fmt.Errorf("%s", "Authentication failed!")
	}

	return id, nil
}

// 用户注册 handler ["/user/api/register", "POST"] ["tested"]
func userRegisterHandler(c *fiber.Ctx) error {
	user := models.User{}
	if err := c.BodyParser(&user); err != nil {
		utils.PrintfColorStr(utils.Red, "post parse error!")

		c.SendStatus(500)
		return c.JSON(fiber.Map{
			"msg":  err.Error(),
			"code": 0,
		})
	}

	if err := user.Register(); err != nil {
		log.Println(err.Error())
		return c.JSON(fiber.Map{
			"msg":  err.Error(),
			"code": -1,
		})
	}

	return c.JSON(fiber.Map{
		"msg":  "register success",
		"code": 1,
	})
}

// 用户登录 handler ["POST", "/api/user/login"] ["tested"]
func userLoginHandler(c *fiber.Ctx) error {
	user := models.User{}
	if err := c.BodyParser(&user); err != nil {
		utils.PrintfColorStr(utils.Red, "post parse error!")
		return c.SendStatus(500)
	}

	if err := user.Login(); err != nil {
		log.Println(err.Error())
		return c.JSON(fiber.Map{
			"msg":  err.Error(),
			"code": -1,
		})
	}

	session, err := store.Get(c)
	if err != nil {
		c.SendStatus(500)
		return c.JSON(fiber.Map{
			"msg":  "web server error!",
			"code": 0,
		})
	}

	// 随机字符作为cookie
	rs := randstring.RandString(randomStringLen)
	for session.Get(rs) != nil {
		rs = randstring.RandString(randomStringLen)
	}
	session.Set(rs, user.Id)
	if err := session.Save(); err != nil {
		log.Println(err.Error())
		c.SendStatus(500)
		return c.JSON(fiber.Map{
			"msg":  "web server error!",
			"code": 0,
		})
	}

	// 更新并设置 cookie
	updateCookie(c, rs)

	return c.JSON(fiber.Map{
		"msg":  "login success",
		"code": 1,
		"id":   user.Id,
		"name": user.Name,
	})
}

// Get方法登录 ["GET", "/user/api/login"] ["tested"]
func userLoginForGet(c *fiber.Ctx) error {
	userFlag := c.Cookies("userFlag", "")
	if userFlag == "" {
		return c.JSON(fiber.Map{
			"msg":  "Authentication failed!",
			"code": -1,
		})
	}
	id, err := verifyCookie(c, userFlag)
	if err != nil {
		return c.JSON(fiber.Map{
			"msg":  err.Error(),
			"code": -1,
		})
	}
	user := &models.User{}
	user.Id = id.(uint64)
	user.GetInfoByID()

	updateCookie(c, userFlag)

	return c.JSON(fiber.Map{
		"msg":  "login success",
		"code": 1,
		"user": fiber.Map{
			"id":   user.Id,
			"name": user.Name,
		},
	})
}

// 用户获取用户信息 handler ["GET", "/api/user/info?uid=[]"]
func userInfoHandler(c *fiber.Ctx) error {
	userFlag := c.Cookies("userFlag", "")
	if userFlag == "" {
		return c.JSON(fiber.Map{
			"msg":  "Authentication failed!",
			"code": -1,
		})
	}

	if _, err := verifyCookie(c, userFlag); err != nil {
		return c.JSON(fiber.Map{
			"msg":  err.Error(),
			"code": -1,
		})
	}

	uid, err := strconv.ParseUint(c.Query("uid", "0"), 10, 64)
	if err != nil {
		c.JSON(fiber.Map{
			"msg": "from `string` to `uint64` error",
		})
		return c.SendStatus(500)
	}
	if uid == 0 {
		return c.JSON(fiber.Map{
			"msg": "not found `uid` at url queries",
		})
	}
	user := models.User{Id: uid}
	if err := user.GetInfoByID(); err != nil {
		return c.JSON(fiber.Map{
			"msg": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"msg":  "success",
		"code": 1,
		"user": user,
	})
}

// 用户更改个人信息 handler ["/user/api/modify", "POST"] ["tested"]
func userModifyInfoHandler(c *fiber.Ctx) error {
	userFlag := c.Cookies("userFlag", "")
	if userFlag == "" {
		return c.JSON(fiber.Map{
			"msg":  "Authentication failed!",
			"code": -1,
		})
	}

	if _, err := verifyCookie(c, userFlag); err != nil {
		return c.JSON(fiber.Map{
			"msg":  err.Error(),
			"code": -1,
		})
	}

	user := models.User{}
	if err := c.BodyParser(&user); err != nil {
		utils.PrintfColorStr(utils.Red, "post parse error!")
		return c.JSON(fiber.Map{
			"msg":  err.Error(),
			"code": 0,
		})
	}
	log.Printf("%v\n", user)
	if err := user.Modify(); err != nil {
		utils.PrintfColorStr(utils.Red, "user modify error!")
		return c.JSON(fiber.Map{
			"msg":  err.Error(),
			"code": -1,
		})
	}

	updateCookie(c, userFlag)

	return c.JSON(fiber.Map{
		"msg":  "update success",
		"code": 1,
	})
}

// 用户借阅 handler ["POST", "/api/user/borrow" POST] ["tested"]
func userBorrowBookHandler(c *fiber.Ctx) error {
	userFlag := c.Cookies("userFlag", "")
	if userFlag == "" {
		return c.JSON(fiber.Map{
			"msg":  "Authentication failed!",
			"code": -1,
		})
	}

	if _, err := verifyCookie(c, userFlag); err != nil {
		return c.JSON(fiber.Map{
			"msg":  err.Error(),
			"code": -1,
		})
	}

	bi := models.BorrowInfo{}
	bm := models.BorrowManager{}

	if err := c.BodyParser(&bi); err != nil {
		utils.PrintfColorStr(utils.Yellow, "[userBorrowBookHandler]Parser Body err")
		return c.JSON(fiber.Map{
			"msg":  err.Error(),
			"code": 0,
		})
	}

	if bi.UserId == 0 || bi.BookISBN == "" || bi.Days == 0 {
		return c.JSON(fiber.Map{
			"msg":  "borrow info exist void variable.",
			"code": -1,
		})
	}
	if err := bm.BorrowBook(bi.UserId, bi.BookISBN, bi.Days); err != nil {
		return c.JSON(fiber.Map{
			"msg":  "borrow book error.",
			"code": -1,
		})
	}

	return c.JSON(fiber.Map{
		"msg":  "borrow success!",
		"code": 1,
	})
}

// 用户还书 handler ["GET", "/api/user/back?id=[borrowId]"] ["tested"]
func userBackBookHandler(c *fiber.Ctx) error {
	userFlag := c.Cookies("userFlag", "")
	if userFlag == "" {
		return c.JSON(fiber.Map{
			"msg":  "Authentication failed!",
			"code": -1,
		})
	}

	if _, err := verifyCookie(c, userFlag); err != nil {
		return c.JSON(fiber.Map{
			"msg":  err.Error(),
			"code": -1,
		})
	}

	var id uint64
	bm := models.BorrowManager{}

	idStr := c.Query("id", "-1")
	if idStr == "-1" {
		utils.PrintfColorStr(utils.Yellow, "[userBorrowBookHandler]Parser Body err")
		return c.JSON(fiber.Map{
			"msg":  "query `id` error!",
			"code": 0,
		})
	}
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return c.JSON(fiber.Map{
			"msg":  "transform `id` to int error!",
			"code": 0,
		})
	}
	if id == 0 {
		return c.JSON(fiber.Map{
			"msg":  "borrow info exist void variable.",
			"code": 0,
		})
	}
	if err := bm.BackBook(id); err != nil {
		return c.JSON(fiber.Map{
			"msg":  "back book error.",
			"code": 0,
		})
	}

	return c.JSON(fiber.Map{
		"code": 1,
		"msg":  "back book success!",
	})
}

// 用户获取所有借阅信息 ["GET", "/api/user/borrowed?id=[userId]&limit&offset"] ["tested"]
func userGetAllBorrow(c *fiber.Ctx) error {
	userFlag := c.Cookies("userFlag", "")
	if userFlag == "" {
		return c.JSON(fiber.Map{
			"msg":  "Authentication failed!",
			"code": -1,
		})
	}

	if _, err := verifyCookie(c, userFlag); err != nil {
		return c.JSON(fiber.Map{
			"msg":  err.Error(),
			"code": -1,
		})
	}

	id, err := strconv.ParseUint(c.Query("id", "0"), 0, 64)
	if err != nil {
		log.Println(err.Error())
		return c.JSON(fiber.Map{
			"msg":  "string to uint64 error",
			"code": 0,
		})
	}
	if id == 0 {
		return c.JSON(fiber.Map{
			"msg":  "parameters exit void variable.",
			"code": 0,
		})
	}

	limit, err := strconv.Atoi(c.Query("limit", "0"))
	if err != nil {
		log.Println(err.Error())
		return c.JSON(fiber.Map{
			"msg":  "string to int error",
			"code": 0,
		})
	}
	offset, err := strconv.Atoi(c.Query("offset", "0"))
	if err != nil {
		log.Println(err.Error())
		return c.JSON(fiber.Map{
			"msg":  "string to int error",
			"code": 0,
		})
	}

	bm := models.BorrowManager{}
	// BUG: 获取用户借阅信息将要通过两种模式：id/name，导致冲突
	bs, err := bm.GetUserAllBorrowById(id, limit, offset)
	if err != nil {
		log.Println(err.Error())
		return c.JSON(fiber.Map{
			"msg":  "Get user all borrow error!",
			"code": -1,
		})
	}

	return c.JSON(fiber.Map{
		"code": 1,
		"msg":  "success",
		"len":  len(bs),
		"bs":   bs,
	})
}

// 获取用户未还借用 ["GET", "/api/user/not_back?id=[userId]&limit&offset"]
func userGetNotBackBorrow(c *fiber.Ctx) error {
	userFlag := c.Cookies("userFlag", "")
	if userFlag == "" {
		return c.JSON(fiber.Map{
			"msg":  "Authentication failed!",
			"code": -1,
		})
	}

	if _, err := verifyCookie(c, userFlag); err != nil {
		return c.JSON(fiber.Map{
			"msg":  err.Error(),
			"code": -1,
		})
	}

	id, _ := strconv.ParseUint(c.Query("id", ""), 0, 64)
	limit, _ := strconv.Atoi(c.Query("limit", "0"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	if id == 0 {
		return c.JSON(fiber.Map{
			"msg":  "parameters exit void variable.",
			"code": 0,
		})
	}
	bm := models.BorrowManager{}
	bs, err := bm.GetUserNotBackBorrowById(id, limit, offset)
	if err != nil {
		return c.JSON(fiber.Map{
			"msg":  "Get user all borrow error!",
			"code": -1,
		})
	}

	return c.JSON(fiber.Map{
		"code": 1,
		"msg":  "success",
		"len":  len(bs),
		"bs":   bs,
	})
}

// 获取用户已还借阅 ["GET", "/api/user/backed?id=[userId]&limit&offset"]
func userGetBackedBorrow(c *fiber.Ctx) error {
	userFlag := c.Cookies("userFlag", "")
	if userFlag == "" {
		return c.JSON(fiber.Map{
			"msg":  "No permission! Please log in.",
			"code": -1,
		})
	}

	if _, err := verifyCookie(c, userFlag); err != nil {
		log.Println(err.Error())
		return c.JSON(fiber.Map{
			"msg":  err.Error(),
			"code": -1,
		})
	}

	id, err := strconv.ParseUint(c.Query("id", "0"), 0, 64)
	if err != nil {
		return c.JSON(fiber.Map{
			"msg":  "id: string to uint64 error",
			"code": 0,
		})
	}
	if id == 0 {
		return c.JSON(fiber.Map{
			"msg":  "parameters exit void variable.",
			"code": 0,
		})
	}
	limit, err := strconv.Atoi(c.Query("limit", "0"))
	if err != nil {
		return c.JSON(fiber.Map{
			"msg":  "limit: string to int error",
			"code": 0,
		})
	}
	offset, err := strconv.Atoi(c.Query("offset", "0"))
	if err != nil {
		return c.JSON(fiber.Map{
			"msg":  "limit: string to int error",
			"code": 0,
		})
	}

	bm := models.BorrowManager{}
	bs, err := bm.GetUserBackedBorrowById(id, limit, offset)
	if err != nil {
		return c.JSON(fiber.Map{
			"msg":  "Get user all borrow error!",
			"code": 0,
		})
	}

	return c.JSON(fiber.Map{
		"code": 1,
		"msg":  "success",
		"len":  len(bs),
		"bs":   bs,
	})
}

// 用户获取图书 ["GET", "/api/user/books?limit=xx&offset=xx&by={0, 1}&typeId=xx"]
// 0: default, 1: typeId
func userGetBooksHandler(c *fiber.Ctx) error {
	userFlag := c.Cookies("userFlag", "")
	if userFlag == "" {
		return c.JSON(fiber.Map{
			"msg":  "Authentication failed!",
			"code": -1,
		})
	}

	if _, err := verifyCookie(c, userFlag); err != nil {
		return c.JSON(fiber.Map{
			"msg":  err.Error(),
			"code": -1,
		})
	}

	limit, err := strconv.Atoi(c.Query("limit", "0"))
	if err != nil {
		log.Println(err.Error())
		c.SendStatus(500)
		return c.JSON(fiber.Map{
			"msg": err.Error(),
		})
	}
	offset, err := strconv.Atoi(c.Query("offset", "0"))
	if err != nil {
		log.Println(err.Error())
		c.SendStatus(500)
		return c.JSON(fiber.Map{
			"msg": err.Error(),
		})
	}

	// 0: default, 1: typeId
	var by int
	if c.Query("by", "0") == "1" {
		by = models.By_BookTypeId
	} else {
		by = models.By_BookDefault
	}

	admin := models.Admin{}
	var books []models.Book
	if by == models.By_BookTypeId {
		typeId, err := strconv.Atoi(c.Query("typeId", "-1"))
		if err != nil {
			log.Println(err.Error())
			c.JSON(fiber.Map{
				"msg": err.Error(),
			})
			return c.SendStatus(400)
		}
		if typeId == -1 {
			books, err = admin.GetBooks(limit, offset, models.By_BookDefault)
		} else {
			books, err = admin.GetBooks(limit, offset, by, typeId)
		}
		if err != nil {
			return c.JSON(fiber.Map{
				"msg":  err.Error(),
				"code": -1,
			})
		}
	} else {
		books, err = admin.GetBooks(limit, offset, by)
	}
	if err != nil {
		log.Println(err.Error())
		c.SendStatus(500)
		return c.JSON(fiber.Map{
			"msg": err.Error(),
		})
	}

	updateCookie(c, userFlag)

	return c.JSON(fiber.Map{
		"msg":   "get book success",
		"len":   len(books),
		"books": books,
	})
}

// 用户搜索图书 ["GET", "/api/user/book/search?by=xx&value=xx"]
// `by: int` => By_BookISBN -> 0 | By_BookName -> 1
func UserSearchBookHandler(c *fiber.Ctx) error {
	userFlag := c.Cookies("userFlag", "")
	if userFlag == "" {
		return c.JSON(fiber.Map{
			"msg":  "Authentication failed!",
			"code": -1,
		})
	}

	if _, err := verifyCookie(c, userFlag); err != nil {
		return c.JSON(fiber.Map{
			"msg":  err.Error(),
			"code": -1,
		})
	}

	var by int
	switch c.Query("by", "-1") {
	case "0":
		by = models.By_BookISBN
		break
	case "1":
		by = models.By_BookName
		break
	default:
		c.SendStatus(400)
		return c.JSON(fiber.Map{
			"msg": "get query param of `by` error",
		})
	}

	value := c.Query("value", "")
	if value == "" {
		c.SendStatus(400)
		return c.JSON(fiber.Map{
			"msg": "get query param of `value` error",
		})
	}

	admin := models.Admin{}
	book, err := admin.SearchBook(by, value)
	if err != nil {
		log.Println(err.Error())
		c.SendStatus(500)
		return c.JSON(fiber.Map{
			"msg": "search book error! " + value,
		})
	}

	updateCookie(c, userFlag)
	bs := []models.Book{}
	bs = append(bs, *book)
	return c.JSON(fiber.Map{
		"msg":   "search success.",
		"books": bs,
	})
}
