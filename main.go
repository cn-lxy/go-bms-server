package main

import (
	"fmt"

	"github.com/cn-lxy/bms_go/controllers"
	"github.com/cn-lxy/bms_go/utils"
	"github.com/cn-lxy/bms_go/utils/config"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	app := fiber.New()
	app.Static("/public", "./static")
	app.Use(logger.New())

	app.Mount("/api/user", controllers.UserController)
	app.Mount("/api/admin", controllers.AdminController)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello World!\nWelcome to GO-BMS!")
	})

	fmt.Printf("App Name: %v\n", config.Cfg.App.Name)
	fmt.Printf("App Author: %v\n", config.Cfg.App.Author)
	fmt.Printf("App Email: %v\n", config.Cfg.App.Email)
	fmt.Printf("Server Port: %v\n", config.Cfg.Server.Port)
	fmt.Printf("Database Host: %v\n", config.Cfg.Database.Host)
	fmt.Printf("Database Port: %v\n", config.Cfg.Database.Port)
	fmt.Printf("Database Name: %v\n", config.Cfg.Database.Name)
	fmt.Printf("Database Username: %v\n", config.Cfg.Database.UserName)
	fmt.Printf("Database Password: %v\n", config.Cfg.Database.Password)

	utils.PrintfColorStr(utils.Green, "App: "+config.Cfg.App.Name+" Ready!")

	app.Listen(fmt.Sprintf(":%d", config.Cfg.Server.Port))
}
