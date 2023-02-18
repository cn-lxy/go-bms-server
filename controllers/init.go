package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

var store *session.Store
var UserController *fiber.App
var AdminController *fiber.App

func init() {
	store = session.New()
	userControllerInit()
	adminControllerInit()
}

// TAG 需要优化一下 api，并且 `user` 侧漏掉了图书查询功能
// 也可将 `user` 和 `admin` 关于图书查询的 api 解耦
func userControllerInit() {
	UserController = fiber.New()
	UserController.Post("/login", userLoginHandler)
	UserController.Get("/login", userLoginForGet)
	UserController.Post("/register", userRegisterHandler)
	UserController.Get("/info", userInfoHandler)
	UserController.Post("/modify", userModifyInfoHandler)
	UserController.Post("/borrow", userBorrowBookHandler)
	UserController.Get("/back", userBackBookHandler)
	UserController.Get("/borrowed", userGetAllBorrow)
	UserController.Get("/not_back", userGetNotBackBorrow)
	UserController.Get("/backed", userGetBackedBorrow)

	UserController.Get("/books", userGetBooksHandler)
	UserController.Get("/book/search", UserSearchBookHandler)
}

func adminControllerInit() {
	AdminController = fiber.New()
	AdminController.Post("/login", adminLoginHandlerPOST)
	AdminController.Get("/login", adminLoginHandlerGET)
	AdminController.Get("/manager/users", adminGetAllUser)
	AdminController.Get("/manager/user/search", adminSearchUserHandler)
	AdminController.Post("/manager/user", adminAddUserHandler)
	AdminController.Post("/manager/user/update", adminUpdateUserHandler)
	AdminController.Get("/manager/user/delete", adminDeleteUserHandler)
	AdminController.Get("/manager/books", adminGetBooksHandler)
	AdminController.Get("/manager/book/search", adminSearchBookHandler)
	AdminController.Post("/manager/book", adminInsertBookHandler)
	AdminController.Get("/manager/book/delete", adminDeleteBookHandler)
	AdminController.Post("/manager/book/update", adminUpdateBookHandler)
	AdminController.Get("/manager/borrow", adminGetBorrowInfoHandler)
	AdminController.Get("/manager/borrow/user", adminGetBorrowOfUserHandler)
	AdminController.Get("/manager/borrow/book", adminGetBorrowOfBookHandler)
	AdminController.Get("/manager/borrow/back", adminCompleteBorrowHandler)
}
