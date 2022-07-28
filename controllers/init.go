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

func userControllerInit() {
	UserController = fiber.New()
	UserController.Post("/login", userLoginHandler)
	UserController.Get("/login", userLoginForGet)
	UserController.Post("/register", userRegisterHandler)
	UserController.Post("/modify", userModifyInfoHandler)
	UserController.Post("/borrow", userBorrowBookHandler)
	UserController.Post("/back", userBackBookHandler)
	UserController.Get("/borrowed", userGetAllBorrow)
	UserController.Get("/not_back", userGetNotBackBorrow)
	UserController.Get("/backed", userGetBackedBorrow)
}

func adminControllerInit() {
	AdminController = fiber.New()

	AdminController.Post("/login", adminLoginHandler_POST)
	AdminController.Get("/login", adminLoginHandler_GET)

	AdminController.Get("/manager/user/getall", adminGetAllUser)
	AdminController.Post("/manager/user/update", adminUpdateUserHandler)
	AdminController.Get("/manager/user/delete", adminDeleteUserHandler)
	AdminController.Get("/manager/book/getbooks", adminGetBooksHandler)
	AdminController.Get("/manager/book/search", adminSearchBookHandler)

	AdminController.Post("/manager/book", adminInsertBookHandler)
	AdminController.Get("/manager/book/delete", adminDeleteBookHandler)
	AdminController.Post("/manager/book/update", adminUpdateBookHandler)
	AdminController.Get("/manager/borrow", adminGetBorrowInfoHandler)
	AdminController.Get("/manager/borrow/u/:userid", adminGetUserBorrowsHandler)
	AdminController.Get("/manager/borrow/b/:isbn", adminGetBookBorrowedUserHandler)
}
