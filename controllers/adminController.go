package controllers

import (
	"fmt"
	"log"

	"github.com/cn-lxy/bms_go/models"
	"github.com/cn-lxy/bms_go/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// admin controller

var AdminController *fiber.App
var db *gorm.DB
var err error
var store *session.Store

func init() {
	// 数据库初始化
	db, err = gorm.Open(sqlite.Open("./models/bms.sqlite"), &gorm.Config{})
	if err != nil {
		panic("failed connect database!")
	}
	db.AutoMigrate(&models.Admin{})
	utils.PrintfColorStr(utils.GREEN, "admin db init completed!")

	// web服务初始化
	store = session.New()
	AdminController = fiber.New()
	AdminController.Post("/login", loginHandler)
	AdminController.Get("/login", loginGet)
	AdminController.Get("/", indexHandler)
}

func loginHandler(c *fiber.Ctx) error {
	admin := &models.Admin{}
	if err := c.BodyParser(admin); err != nil {
		fmt.Println("parse error!")
	}
	log.Printf("post admin: {account: %s, password: %s}", admin.Account, admin.Password)

	dbAdmin := models.Admin{}
	dbAdmin.Query(db, models.ACCOUNT, admin.Account)

	if dbAdmin.Password != admin.Password {
		return c.SendString("login error!")
	}

	// 登录成功
	log.Println("login success!")
	session, err := store.Get(c)
	if err != nil {
		panic(err)
	}
	session.Set("name", "dave")

	// Destroy session
	// if err := session.Destroy(); err != nil {
	// 	panic(err)
	// }

	// Save session
	if err := session.Save(); err != nil {
		panic(err)
	}

	return c.SendString("login success.")
}

func loginGet(c *fiber.Ctx) error {
	return c.SendString("login get")
}

// admin index page
func indexHandler(c *fiber.Ctx) error {
	session, err := store.Get(c)
	if err != nil {
		panic(err)
	}

	// Get value
	name := session.Get("name")
	if name == nil {
		c.Redirect("/admin/login")
	}

	session.Set("name", "dave2")

	// Save session
	if err := session.Save(); err != nil {
		panic(err)
	}

	fmt.Printf("name: %v\n", name)
	return c.SendString("admin index page")
}
