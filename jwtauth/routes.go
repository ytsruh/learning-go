package jwtauth

import (
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func SetupRoutes(app *fiber.App) {

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	app.Post("/register", func(c *fiber.Ctx) error {
		var data map[string]string
		if err := c.BodyParser(&data); err != nil {
			return err
		}
		password, err := bcrypt.GenerateFromPassword([]byte(data["password"]), 14)
		if err != nil {
			panic("Issue creating password")
		}
		user := User{
			Name:     data["name"],
			Email:    data["email"],
			Password: password,
		}

		DB.Create(&user)

		return c.JSON(user)
	})

	app.Post("/login", func(c *fiber.Ctx) error {
		var data map[string]string
		if err := c.BodyParser(&data); err != nil {
			return err
		}

		var user User

		DB.Where("email = ?", data["email"]).First(&user)

		if user.Id == 0 {
			c.Status(fiber.StatusNotFound)
			return c.JSON(fiber.Map{
				"message": "user not found",
			})
		}

		if err := bcrypt.CompareHashAndPassword(user.Password, []byte(data["password"])); err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{
				"message": "incorrect password",
			})
		}

		return c.JSON(user)
	})

}
