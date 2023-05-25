package jwtauth

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func RunServer() {
	ConnectDb()

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
	}))

	SetupCookieRoutes(app)
	SetupTokenRoutes(app)

	app.Listen(":3000")
}
