package main

import (
	"github.com/cn-lxy/bms_go/controllers"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	app := fiber.New()

	app.Use(logger.New())

	app.Mount("/user/api", controllers.UserController)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello World!\nWelcome to GO-BMS!")
	})

	app.Listen(":8080")
}
