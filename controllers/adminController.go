package controllers

import (
	"fmt"
	"log"
	"strconv"

	"github.com/cn-lxy/bms_go/models"
	"github.com/gofiber/fiber/v2"
)

//! code: { -1: 身份验证失败, -2: 客户端请求错误, -3: 服务器错误, 0: 未知错误, 1: 请求成功 }

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
			"msg":  "get session field!",
			"code": -3,
		})
	}

	admin := models.Admin{
		Id:       0,
		Account:  "",
		Password: "",
	}
	if err := c.BodyParser(&admin); err != nil {
		log.Println(err.Error())
		return c.JSON(fiber.Map{
			"msg":  "body parser error",
			"code": -2,
		})
	}
	if err := admin.Login(); err != nil {
		log.Println(err.Error())
		return c.JSON(fiber.Map{
			"msg":  "login field!",
			"code": -1,
		})
	}

	session.Set("login", "true")
	if err := session.Save(); err != nil {
		log.Println(err.Error())
		c.SendStatus(500)
		return c.JSON(fiber.Map{
			"msg":  "set session field!",
			"code": -3, // 服务器错误
		})
	}

	return c.JSON(fiber.Map{
		"msg":  "login success",
		"code": 1,
	})
}

// 管理员登录 ["GET", "/admin/api/login"] ["tested"]
func adminLoginHandler_GET(c *fiber.Ctx) error {
	if err := verifyLogin_middle(c); err != nil {
		log.Println(err.Error())
		return c.JSON(fiber.Map{
			"msg":  err.Error(),
			"code": -1,
		})
	}

	return c.JSON(fiber.Map{
		"msg":  "login success",
		"code": 1,
	})
}

// 获取所有用户 ["GET", "/api/admin/manager/users?limit=xx&offset=xx"] ["tested"]
func adminGetAllUser(c *fiber.Ctx) error {
	if err := verifyLogin_middle(c); err != nil {
		log.Println(err.Error())
		return c.JSON(fiber.Map{
			"msg":  err.Error(),
			"code": -1,
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

	if limit == 0 {
		return c.JSON(fiber.Map{
			"msg":  "query param error",
			"code": -2,
		})
	}

	users, err := admin.GetAllUser(limit, offset)
	if err != nil {
		log.Println(err.Error())
		return c.JSON(fiber.Map{
			"msg":  err.Error(),
			"code": -3,
		})
	}

	return c.JSON(fiber.Map{
		"len":   len(users),
		"users": users,
		"code":  1,
	})
}

// 搜索用户 ["GET", "/api/admin/manager/user/search?by=xx&key=xx"]
// by: { 0:姓名, 1:学号, 2: 性别, 3: uid }
func adminSearchUserHandler(c *fiber.Ctx) error {
	if err := verifyLogin_middle(c); err != nil {
		log.Println(err.Error())
		return c.JSON(fiber.Map{
			"msg":  err.Error(),
			"code": -1,
		})
	}

	admin := models.Admin{}
	limit, err := strconv.Atoi(c.Query("limit", "0"))
	if err != nil {
		log.Println(err.Error())
		c.SendStatus(500)
		return c.JSON(fiber.Map{
			"msg":  err.Error(),
			"code": -3,
		})
	}
	offset, err := strconv.Atoi(c.Query("offset", "0"))
	if err != nil {
		log.Println(err.Error())
		c.SendStatus(500)
		return c.JSON(fiber.Map{
			"msg":  err.Error(),
			"code": -3,
		})
	}
	by, err := strconv.ParseInt(c.Query("by", "-1"), 10, 64)
	if err != nil {
		log.Println(err.Error())
		c.SendStatus(500)
		return c.JSON(fiber.Map{
			"msg":  err.Error(),
			"code": -3,
		})
	}
	key := c.Query("key", "")

	if by == -1 || key == "" || limit == 0 {
		return c.JSON(fiber.Map{
			"msg":  "param is illegal",
			"code": -2,
		})
	}
	users, err := admin.SearchUser(by, key, limit, offset)
	if err != nil {
		return c.JSON(fiber.Map{
			"msg":  err.Error(),
			"code": -3,
		})
	}
	return c.JSON(fiber.Map{
		"msg":   "search success",
		"users": users,
		"code":  1,
	})
}

// 添加用户 ["POST", "/api/admin/manager/user"]
func adminAddUserHandler(c *fiber.Ctx) error {
	if err := verifyLogin_middle(c); err != nil {
		log.Println(err.Error())
		return c.JSON(fiber.Map{
			"msg":  err.Error(),
			"code": -1,
		})
	}

	newUser := models.User{}
	if err := c.BodyParser(&newUser); err != nil {
		log.Println(err.Error())
		return c.JSON(fiber.Map{
			"msg":  "parse body to `user` error!",
			"code": -2,
		})
	}
	admin := models.Admin{}
	if err := admin.AddUser(&newUser); err != nil {
		log.Println(err.Error())
		return c.JSON(fiber.Map{
			"msg":  "add user failed!",
			"code": -3,
		})
	}

	return c.JSON(fiber.Map{
		"msg":  "add user success.",
		"code": 1,
	})
}

// 更新用户信息 ["POST", "/api/admin/manager/user/update"] ["not test"]
func adminUpdateUserHandler(c *fiber.Ctx) error {
	if err := verifyLogin_middle(c); err != nil {
		log.Println(err.Error())
		return c.JSON(fiber.Map{
			"msg":  err.Error(),
			"code": -1,
		})
	}

	updateUser := models.User{}
	if err := c.BodyParser(&updateUser); err != nil {
		log.Println(err.Error())
		return c.JSON(fiber.Map{
			"msg":  "parse body to `user` error!",
			"code": -2,
		})
	}

	admin := models.Admin{}
	if err := admin.ModifyUser(&updateUser); err != nil {
		log.Println(err.Error())
		return c.JSON(fiber.Map{
			"msg":  "update user info field!",
			"code": -3,
		})
	}

	return c.JSON(fiber.Map{
		"msg":  "update user info success.",
		"code": 1,
	})
}

// 根据user_id删除用户 ["GET", "/api/admin/manager/user/delete?id=xx"] ["tested"]
func adminDeleteUserHandler(c *fiber.Ctx) error {
	if err := verifyLogin_middle(c); err != nil {
		log.Println(err.Error())
		return c.JSON(fiber.Map{
			"msg":  err.Error(),
			"code": -1,
		})
	}

	userId, err := strconv.ParseUint(c.Query("id", "0"), 0, 64)
	if err != nil {
		log.Println("parse uint:", err.Error())
		c.SendStatus(500)
		return c.JSON(fiber.Map{
			"msg":  err.Error(),
			"code": -3,
		})
	}
	if userId == 0 {
		return c.JSON(fiber.Map{
			"msg":  "user id is unlawful!",
			"code": -2,
		})
	}

	admin := models.Admin{}
	if err := admin.DeleteUser(userId); err != nil {
		log.Println("delete user:", err.Error())
		return c.JSON(fiber.Map{
			"msg":  "delete user filed!",
			"code": -3,
		})
	}

	// 删除成功
	return c.JSON(fiber.Map{
		"msg":  "delete user success",
		"code": 1,
	})
}

// 获取多本图书 ["GET", "/api/admin/manager/books?limit=xx&offset=xx&by={0, 1}&typeid=xx"]
// 0: default, 1: typeId
func adminGetBooksHandler(c *fiber.Ctx) error {
	if err := verifyLogin_middle(c); err != nil {
		log.Println(err.Error())
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
			"msg":  err.Error(),
			"code": -3,
		})
	}
	offset, err := strconv.Atoi(c.Query("offset", "0"))
	if err != nil {
		log.Println(err.Error())
		c.SendStatus(500)
		return c.JSON(fiber.Map{
			"msg":  err.Error(),
			"code": -3,
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
		typeId, err := strconv.Atoi(c.Query("typeid", "1"))
		if err != nil {
			log.Println(err.Error())
			return c.JSON(fiber.Map{
				"msg":  err.Error(),
				"code": -3,
			})
		}
		books, err = admin.GetBooks(limit, offset, by, typeId)
	} else {
		books, err = admin.GetBooks(limit, offset, by)
	}
	if err != nil {
		log.Println(err.Error())
		c.SendStatus(500)
		return c.JSON(fiber.Map{
			"msg":  err.Error(),
			"code": -3,
		})
	}

	return c.JSON(fiber.Map{
		"msg":   "get book success",
		"len":   len(books),
		"books": books,
		"code":  1,
	})
}

// 搜索图书 ["GET", "/api/admin/manager/book/search?by=xx&value=xx"]
// `by: int` => By_BookISBN -> 0 | By_BookName -> 1
func adminSearchBookHandler(c *fiber.Ctx) error {
	if err := verifyLogin_middle(c); err != nil {
		log.Println(err.Error())
		return c.JSON(fiber.Map{
			"msg":  err.Error(),
			"code": -1,
		})
	}

	var by int
	switch c.Query("by", "-1") {
	case "0":
		by = models.By_BookISBN
	case "1":
		by = models.By_BookName
	default:
		return c.JSON(fiber.Map{
			"msg":  "get query param of `by` error",
			"code": -2,
		})
	}

	value := c.Query("value", "")
	if value == "" {
		return c.JSON(fiber.Map{
			"msg":  "get query param of `value` error",
			"code": -2,
		})
	}

	admin := models.Admin{}
	book, err := admin.SearchBook(by, value)
	if err != nil {
		log.Println(err.Error())
		c.SendStatus(500)
		return c.JSON(fiber.Map{
			"msg":  "search book error! " + value,
			"code": -3,
		})
	}

	return c.JSON(fiber.Map{
		"msg":   "search success.",
		"books": []models.Book{*book},
		"code":  1,
	})
}

// 新建图书 handler ["POST", "/api/admin/manager/book"]
func adminInsertBookHandler(c *fiber.Ctx) error {
	if err := verifyLogin_middle(c); err != nil {
		log.Println(err.Error())
		return c.JSON(fiber.Map{
			"msg":  err.Error(),
			"code": -1,
		})
	}

	book := models.Book{}
	if err := c.BodyParser(&book); err != nil {
		log.Println(err.Error())
		return c.JSON(fiber.Map{
			"msg":  err.Error(),
			"code": -2,
		})
	}

	admin := models.Admin{}
	if err := admin.InsertBook(&book); err != nil {
		log.Println(err.Error())
		c.SendStatus(500)
		return c.JSON(fiber.Map{
			"msg":  err.Error(),
			"code": -3,
		})
	}

	return c.JSON(fiber.Map{
		"msg":  "insert success",
		"code": 1,
	})
}

// 删除书籍 ["GET", "/api/admin/manager/book/delete?isbn=xx"]
func adminDeleteBookHandler(c *fiber.Ctx) error {
	if err := verifyLogin_middle(c); err != nil {
		log.Println(err.Error())
		return c.JSON(fiber.Map{
			"msg":  err.Error(),
			"code": -1,
		})
	}

	isbn := c.Query("isbn", "")
	if isbn == "" {
		return c.JSON(fiber.Map{
			"msg":  "query param `isbn` filed!",
			"code": -2,
		})
	}

	admin := models.Admin{}
	if err := admin.DeleteBook(isbn); err != nil {
		log.Println(err.Error())
		return c.JSON(fiber.Map{
			"msg":  "delete book field!",
			"isbn": isbn,
			"code": -3,
		})
	}

	return c.JSON(fiber.Map{
		"msg":  "delete book success!",
		"isbn": isbn,
		"code": 1,
	})
}

// 更新书籍 ["POST", "/api/admin/manager/book/update"]
func adminUpdateBookHandler(c *fiber.Ctx) error {
	if err := verifyLogin_middle(c); err != nil {
		log.Println(err.Error())
		return c.JSON(fiber.Map{
			"msg":  err.Error(),
			"code": -1,
		})
	}

	book := models.Book{}
	if err := c.BodyParser(&book); err != nil {
		log.Println(err.Error())
		c.JSON(fiber.Map{
			"msg":  err.Error(),
			"code": -2,
		})
		// HTTP-Code: `400` => `Bad Request` 客户端请求的语法错误，服务器无法理解
		return c.SendStatus(400)
	}

	admin := models.Admin{}
	if err := admin.UpdateBook(&book); err != nil {
		log.Println(err.Error())
		return c.JSON(fiber.Map{
			"msg":  "update book field!",
			"name": book.Name,
			"code": -3,
		})
	}

	return c.JSON(fiber.Map{
		"msg":  "update book success!",
		"name": book.Name,
		"code": 1,
	})
}

// 获取借阅记录信息 ["GET", "/api/admin/manager/borrow?type=xx&limit=xx&&offset=xx"] ["not test"]
// query param: `type` => {1: 全部，2: 未归还，3: 已归还}
func adminGetBorrowInfoHandler(c *fiber.Ctx) error {
	if err := verifyLogin_middle(c); err != nil {
		log.Println(err.Error())
		return c.JSON(fiber.Map{
			"msg":  err.Error(),
			"code": -1,
		})
	}
	_type, err := strconv.ParseUint(c.Query("type", "0"), 10, 64)
	if err != nil {
		return err
	}
	if _type == 0 {
		return c.JSON(fiber.Map{
			"msg":  "query parma `type` is void",
			"code": -2,
		})
	}
	limit, err := strconv.Atoi(c.Query("limit", "0"))
	if err != nil {
		log.Println(err.Error())
		c.JSON(fiber.Map{
			"msg":  err.Error(),
			"code": -3,
		})
		return c.SendStatus(500)
	}
	offset, err := strconv.Atoi(c.Query("offset", "0"))
	if err != nil {
		log.Println(err.Error())
		c.JSON(fiber.Map{
			"msg":  err.Error(),
			"code": -3,
		})
		return c.SendStatus(500)
	}

	admin := models.Admin{}
	// borrowInfo slice => bis
	bis, err := admin.GetBorrowInfo(_type, limit, offset)
	if err != nil {
		log.Println(err.Error())
		c.JSON(fiber.Map{
			"msg":  err.Error(),
			"code": -3,
		})
		return c.SendStatus(500)
	}

	return c.JSON(fiber.Map{
		"msg":  "success",
		"len":  len(bis),
		"bis":  bis,
		"code": 1,
	})
}

// 获取用户借阅信息 ["GET", "/api/admin/manager/borrow/u/?name=xx&type=xx&limit=xx&offset=xx"]
// type = {"1":"全部借阅", "2":"未还借阅", "3":"已还借阅"}
func adminGetBorrowOfUserHandler(c *fiber.Ctx) error {
	if err := verifyLogin_middle(c); err != nil {
		log.Println(err.Error())
		return c.JSON(fiber.Map{
			"msg":  err.Error(),
			"code": -1,
		})
	}

	// 参数1：数字字符串，参数2：数字进制，参数3：bit大小 uint8,uint32..
	userName := c.Query("name", "")
	if userName == "" {
		return c.JSON(fiber.Map{
			"msg":  "Not found `user name` query params",
			"code": -2,
		})
	}
	_type := c.Query("type", "-1")
	limit, err := strconv.Atoi(c.Query("limit", "0"))
	if err != nil {
		log.Println(err.Error())
		return c.JSON(fiber.Map{
			"msg":  err.Error(),
			"code": -3,
		})
	}
	offset, err := strconv.Atoi(c.Query("offset", "0"))
	if err != nil {
		log.Println(err.Error())
		return c.JSON(fiber.Map{
			"msg":  err.Error(),
			"code": -3,
		})
	}

	admin := models.Admin{}
	var bis []models.BorrowInfo
	switch _type {
	case "1":
		bis, err = admin.GetUserBorrows(userName, limit, offset)
	case "2":
		bis, err = admin.GetUserNotBackBorrow(userName, limit, offset)
	case "3":
		bis, err = admin.GetUserBackedBorrow(userName, limit, offset)
	default:
		return c.JSON(fiber.Map{
			"msg":  "Get query params error",
			"code": -2,
		})
	}
	if err != nil {
		log.Println(err.Error())
		return c.JSON(fiber.Map{
			"msg":  err.Error(),
			"code": -3,
		})
	}

	return c.JSON(fiber.Map{
		"msg":  "success",
		"len":  len(bis),
		"bis":  bis,
		"code": 1,
	})
}

// 获取被借书籍用户 ["GET", "/api/admin/manager/borrow/b?name=xx&type=xx&limit=xx&offset=xx"]
// type = {"1":"全部借阅", "2":"未还借阅", "3":"已还借阅"}
func adminGetBorrowOfBookHandler(c *fiber.Ctx) error {
	if err := verifyLogin_middle(c); err != nil {
		log.Println(err.Error())
		return c.JSON(fiber.Map{
			"msg":  err.Error(),
			"code": -1,
		})
	}

	bookName := c.Query("name", "")
	_type := c.Query("type", "-1")
	if bookName == "" {
		return c.JSON(fiber.Map{
			"msg":  "Not found `borrow id`",
			"code": -2,
		})
	}
	limit, err := strconv.Atoi(c.Query("limit", "0"))
	if err != nil {
		log.Println(err.Error())
		return c.JSON(fiber.Map{
			"msg":  err.Error(),
			"code": -3,
		})
	}
	offset, err := strconv.Atoi(c.Query("offset", "0"))
	if err != nil {
		log.Println(err.Error())
		return c.JSON(fiber.Map{
			"msg":  err.Error(),
			"code": -3,
		})
	}

	admin := models.Admin{}
	bis := []models.BorrowInfo{}
	switch _type {
	case "1":
		bis, err = admin.GetBookBorrowAll(bookName, limit, offset)
	case "2":
		bis, err = admin.GetBookNotBackBorrow(bookName, limit, offset)
	case "3":
		bis, err = admin.GetBookBackedBorrow(bookName, limit, offset)
	default:
		return c.JSON(fiber.Map{
			"msg":  "Get query params error",
			"code": -2,
		})
	}
	if err != nil {
		log.Println(err.Error())
		return c.JSON(fiber.Map{
			"msg":  err.Error(),
			"code": -3,
		})
	}
	return c.JSON(fiber.Map{
		"msg":  "success",
		"len":  len(bis),
		"bis":  bis,
		"code": 1,
	})
}

// 管理员完成借阅 handler ["GET", "/api/admin/manager/borrow/back?bid=xxx"]
// query param: `bid` => borrow id
func adminCompleteBorrowHandler(c *fiber.Ctx) error {
	// 身份验证
	if err := verifyLogin_middle(c); err != nil {
		log.Println(err.Error())
		return c.JSON(fiber.Map{
			"msg":  err.Error(),
			"code": -1,
		})
	}

	bid, err := strconv.ParseUint(c.Query("bid", "0"), 10, 64)
	if err != nil {
		return err
	}
	if bid == 0 {
		return c.JSON(fiber.Map{
			"msg":  "query param `bid` is void",
			"code": -2,
		})
	}

	admin := models.Admin{}
	if err != admin.CompleteBorrow(bid) {
		return c.JSON(fiber.Map{
			"msg":  err.Error(),
			"code": -3,
		})
	}

	return c.JSON(fiber.Map{
		"msg":  "success",
		"code": 1,
	})
}
