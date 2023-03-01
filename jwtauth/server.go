package jwtauth

import (
	"github.com/gofiber/fiber/v2"
)

func RunServer() {
	ConnectDb()

	app := fiber.New()

	SetupRoutes(app)

	app.Listen(":3000")
}
