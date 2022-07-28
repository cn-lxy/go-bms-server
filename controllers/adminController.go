package controllers

import (
	"fmt"
	"log"
	"strconv"

	"github.com/cn-lxy/bms_go/models"
	"github.com/gofiber/fiber/v2"
)

func verifyLogin_middle(c *fiber.Ctx) error {
	session, err := store.Get(c)
	if err != nil {
		return fmt.Errorf("%s", "get session field!")
	}
	login := session.Get("login")
	if login == nil {
		return fmt.Errorf("%s", "get session data of `login` error!")
	}

	if login == "true" {
		return nil
	}
	return fmt.Errorf("%s", "login field")
}

// 管理员登录 ["POST", "/admin/api/login"] ["tested"]
func adminLoginHandler_POST(c *fiber.Ctx) error {
	session, err := store.Get(c)
	if err != nil {
		c.SendStatus(500)
		return c.JSON(fiber.Map{
			"msg": "get session field!",
		})
	}

	admin := models.Admin{
		Id:       0,
		Account:  "",
		Password: "",
	}
	if err := c.BodyParser(&admin); err != nil {
		log.Println(err.Error())
		c.SendStatus(400)
		return c.SendString(`{"msg": "body parser error"}`)
	}
	if err := admin.Login(); err != nil {
		log.Println(err.Error())
		c.SendStatus(401)
		return c.SendString(`{"msg": "login field!"}`)
	}

	session.Set("login", "true")
	if err := session.Save(); err != nil {
		log.Println(err.Error())
		c.SendStatus(500)
		return c.JSON(fiber.Map{
			"msg": "set session field!",
		})
	}

	return c.JSON(fiber.Map{
		"msg": "login success",
	})
}

// 管理员登录 ["GET", "/admin/api/login"] ["tested"]
func adminLoginHandler_GET(c *fiber.Ctx) error {
	if err := verifyLogin_middle(c); err != nil {
		log.Println(err.Error())
		c.SendStatus(401)
		return c.JSON(fiber.Map{
			"msg": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"msg": "login success",
	})
}

// 获取所有用户 ["GET", "admin/api/manager/user/getall?limit=xx&offset=xx"] ["tested"]
func adminGetAllUser(c *fiber.Ctx) error {
	if err := verifyLogin_middle(c); err != nil {
		log.Println(err.Error())
		c.SendStatus(401)
		return c.JSON(fiber.Map{
			"msg": err.Error(),
		})
	}

	admin := models.Admin{}
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

	users, err := admin.GetAllUser(limit, offset)
	if err != nil {
		log.Println(err.Error())
		c.SendStatus(500)
		return c.JSON(fiber.Map{
			"msg": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"len":   len(users),
		"users": users,
	})
}

// 更新用户信息 ["POST", "/admin/api/manager/user/update"] ["not test"]
func adminUpdateUserHandler(c *fiber.Ctx) error {
	if err := verifyLogin_middle(c); err != nil {
		log.Println(err.Error())
		c.SendStatus(401)
		return c.JSON(fiber.Map{
			"msg": err.Error(),
		})
	}

	updateUser := models.User{}
	if err := c.BodyParser(&updateUser); err != nil {
		log.Println(err.Error())
		c.SendStatus(400)
		return c.JSON(fiber.Map{
			"msg": "parse body to `user` error!",
		})
	}

	admin := models.Admin{}
	if err := admin.ModifyUser(&updateUser); err != nil {
		log.Println(err.Error())
		c.SendStatus(500)
		return c.JSON(fiber.Map{
			"msg": "update user info field!",
		})
	}

	return c.JSON(fiber.Map{
		"msg": "update user info success.",
	})
}

// 根据user_id删除用户 ["GET", "/admin/api/manager/user/delete?id=xx"] ["tested"]
func adminDeleteUserHandler(c *fiber.Ctx) error {
	if err := verifyLogin_middle(c); err != nil {
		log.Println(err.Error())
		c.SendStatus(401)
		return c.JSON(fiber.Map{
			"msg": err.Error(),
		})
	}

	userId, err := strconv.ParseUint(c.Query("id", "0"), 0, 64)
	if err != nil {
		log.Println("parse uint:", err.Error())
		c.SendStatus(500)
		return c.JSON(fiber.Map{
			"msg": err.Error(),
		})
	}
	if userId == 0 {
		c.SendStatus(400)
		return c.JSON(fiber.Map{
			"msg": "user id is unlawful!",
		})
	}

	admin := models.Admin{}
	if err := admin.DeleteUser(userId); err != nil {
		log.Println("delete user:", err.Error())
		c.SendStatus(500)
		return c.JSON(fiber.Map{
			"msg": "delete user filed!",
		})
	}

	// 删除成功
	return c.JSON(fiber.Map{
		"msg": "delete user success",
	})
}

// 获取多本图书 ["GET", "/admin/api/manager/book/getbooks?limit=xx&offset=xx&by={0, 1}&typeid=xx"]
// 0: default, 1: typeId
// ["not test"]
func adminGetBooksHandler(c *fiber.Ctx) error {
	if err := verifyLogin_middle(c); err != nil {
		log.Println(err.Error())
		c.SendStatus(401)
		return c.JSON(fiber.Map{
			"msg": err.Error(),
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
		by = models.By_BookTypeId
	}

	admin := models.Admin{}
	var books []models.Book
	if by == models.By_BookTypeId {
		typeId, err := strconv.Atoi(c.Query("typeid", "1"))
		if err != nil {
			log.Println(err.Error())
			c.JSON(fiber.Map{
				"msg": err.Error(),
			})
			return c.SendStatus(400)
		}
		books, err = admin.GetBooks(limit, offset, by, typeId)
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

	return c.JSON(fiber.Map{
		"msg":   "get book success",
		"len":   len(books),
		"books": books,
	})
}

// 搜索图书 ["GET", "/admin/api/manager/book/search/by=xx&value=xx"] ["not test"]
// `by: int` => By_BookISBN -> 0 | By_BookName -> 1
func adminSearchBookHandler(c *fiber.Ctx) error {
	if err := verifyLogin_middle(c); err != nil {
		log.Println(err.Error())
		c.SendStatus(401)
		return c.JSON(fiber.Map{
			"msg": err.Error(),
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

	return c.JSON(fiber.Map{
		"msg":  "search success.",
		"book": book,
	})
}

// 新建图书 handler ["POST", "/admin/api/manager/book"] ["not tested"]
func adminInsertBookHandler(c *fiber.Ctx) error {
	if err := verifyLogin_middle(c); err != nil {
		log.Println(err.Error())
		c.SendStatus(401)
		return c.JSON(fiber.Map{
			"msg": err.Error(),
		})
	}

	book := models.Book{}
	if err := c.BodyParser(&book); err != nil {
		log.Println(err.Error())
		return c.JSON(fiber.Map{
			"msg": err.Error(),
		})
	}

	admin := models.Admin{}
	if err := admin.InsertBook(&book); err != nil {
		log.Println(err.Error())
		c.SendStatus(500)
		return c.JSON(fiber.Map{
			"msg": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"msg": "insert success",
	})
}

// 删除书籍 ["GET", "/admin/api/manager/book/delete?isbn=xx"]
func adminDeleteBookHandler(c *fiber.Ctx) error {
	if err := verifyLogin_middle(c); err != nil {
		log.Println(err.Error())
		c.SendStatus(401)
		return c.JSON(fiber.Map{
			"msg": err.Error(),
		})
	}

	isbn := c.Query("isbn", "")
	if isbn == "" {
		return c.JSON(fiber.Map{
			"msg": "query param `isbn` filed!",
		})
	}

	admin := models.Admin{}
	if err := admin.DeleteBook(isbn); err != nil {
		log.Println(err.Error())
		c.SendStatus(500)
		return c.JSON(fiber.Map{
			"msg":  "delete book field!",
			"isbn": isbn,
		})
	}

	return c.JSON(fiber.Map{
		"msg":  "delete book success!",
		"isbn": isbn,
	})
}

// 更新书籍 ["POST", "/admin/api/manager/book/update"]
func adminUpdateBookHandler(c *fiber.Ctx) error {
	if err := verifyLogin_middle(c); err != nil {
		log.Println(err.Error())
		c.SendStatus(401)
		return c.JSON(fiber.Map{
			"msg": err.Error(),
		})
	}

	book := models.Book{}
	if err := c.BodyParser(&book); err != nil {
		log.Println(err.Error())
		c.JSON(fiber.Map{
			"msg": err.Error(),
		})
		// HTTP-Code: `400` => `Bad Request` 客户端请求的语法错误，服务器无法理解
		return c.SendStatus(400)
	}

	admin := models.Admin{}
	if err := admin.UpdateBook(&book); err != nil {
		log.Println(err.Error())
		c.SendStatus(500)
		return c.JSON(fiber.Map{
			"msg":  "update book field!",
			"name": book.Name,
		})
	}

	return c.JSON(fiber.Map{
		"msg":  "update book success!",
		"name": book.Name,
	})
}

// 获取借阅记录信息 ["GET", "/admin/api/manager/borrow?limit=xx&&offset=xx"] ["not test"]
func adminGetBorrowInfoHandler(c *fiber.Ctx) error {
	if err := verifyLogin_middle(c); err != nil {
		log.Println(err.Error())
		c.SendStatus(401)
		return c.JSON(fiber.Map{
			"msg": err.Error(),
		})
	}

	limit, err := strconv.Atoi(c.Query("limit", "0"))
	if err != nil {
		log.Println(err.Error())
		c.JSON(fiber.Map{
			"msg": err.Error(),
		})
		return c.SendStatus(500)
	}
	offset, err := strconv.Atoi(c.Query("offset", "0"))
	if err != nil {
		log.Println(err.Error())
		c.JSON(fiber.Map{
			"msg": err.Error(),
		})
		return c.SendStatus(500)
	}

	admin := models.Admin{}
	// borrowInfo slice => bis
	bis, err := admin.GetBorrowInfo(limit, offset)
	if err != nil {
		log.Println(err.Error())
		c.JSON(fiber.Map{
			"msg": err.Error(),
		})
		return c.SendStatus(500)
	}

	return c.JSON(fiber.Map{
		"msg": "success",
		"len": len(bis),
		"bis": bis,
	})
}

// 获取用户借阅信息 ["GET", "/admin/api/manager/borrow/u/{userid}?type=xx&limit=xx&offset=xx"] ["not test"]
// type = {"1":"全部借阅", "2":"未还借阅", "3":"已还借阅"}
func adminGetUserBorrowsHandler(c *fiber.Ctx) error {
	if err := verifyLogin_middle(c); err != nil {
		log.Println(err.Error())
		c.SendStatus(401)
		return c.JSON(fiber.Map{
			"msg": err.Error(),
		})
	}

	// 参数1：数字字符串，参数2：数字进制，参数3：bit大小 uint8,uint32..
	userId, err := strconv.ParseUint(c.Params("userid", "0"), 10, 64)
	if err != nil {
		log.Println(err.Error())
		c.JSON(fiber.Map{
			"msg": err.Error(),
		})
		return c.SendStatus(400)
	}
	if userId == 0 {
		c.JSON(fiber.Map{
			"msg": "Not found `user id`",
		})
		return c.SendStatus(400)
	}

	limit, err := strconv.Atoi(c.Query("limit", "0"))
	if err != nil {
		log.Println(err.Error())
		c.JSON(fiber.Map{
			"msg": err.Error(),
		})
		return c.SendStatus(500)
	}
	offset, err := strconv.Atoi(c.Query("offset", "0"))
	if err != nil {
		log.Println(err.Error())
		c.JSON(fiber.Map{
			"msg": err.Error(),
		})
		return c.SendStatus(500)
	}

	admin := models.Admin{}
	_type := c.Query("type", "-1")
	var bis []models.BorrowInfo
	switch _type {
	case "1":
		bis, err = admin.GetUserBorrows(userId, limit, offset)
		break
	case "2":
		bis, err = admin.GetUserNotBackBorrow(userId, limit, offset)
		break
	case "3":
		bis, err = admin.GetUserBackedBorrow(userId, limit, offset)
		break
	default:
		c.JSON(fiber.Map{
			"msg": "Get query params error",
		})
		return c.SendStatus(400)
	}
	if err != nil {
		log.Println(err.Error())
		c.JSON(fiber.Map{
			"msg": err.Error(),
		})
		return c.SendStatus(500)
	}

	return c.JSON(fiber.Map{
		"msg": "success",
		"len": len(bis),
		"bis": bis,
	})
}

// 获取被借书籍用户 ["GET", "/admin/api/manager/borrow/b/{isbn}"]
func adminGetBookBorrowedUserHandler(c *fiber.Ctx) error {
	if err := verifyLogin_middle(c); err != nil {
		log.Println(err.Error())
		c.SendStatus(401)
		return c.JSON(fiber.Map{
			"msg": err.Error(),
		})
	}

	isbn := c.Params("isbn", "-1")
	if isbn == "-1" {
		c.JSON(fiber.Map{
			"msg": "Not found `book id: isbn`",
		})
		return c.SendStatus(400)
	}

	admin := models.Admin{}
	bis, err := admin.GetBookBorrowedUser(isbn)
	if err != nil {
		log.Println(err.Error())
		c.JSON(fiber.Map{
			"msg": err.Error(),
		})
		return c.SendStatus(500)
	}

	return c.JSON(fiber.Map{
		"msg": "success",
		"len": len(bis),
		"bis": bis,
	})
}
