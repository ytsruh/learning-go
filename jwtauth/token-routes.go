package jwtauth

import (
	"time"

	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v2"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

// Protected protect routes
func Protected() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey:   []byte(SecretKey),
		ErrorHandler: jwtError,
	})
}

func jwtError(c *fiber.Ctx, err error) error {
	if err.Error() == "Missing or malformed JWT" {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"status": "error", "message": "Missing or malformed JWT", "data": nil})
	}
	return c.Status(fiber.StatusUnauthorized).
		JSON(fiber.Map{"status": "error", "message": "Invalid or expired JWT", "data": nil})
}

func SetupTokenRoutes(app *fiber.App) {

	tokenRoutes := app.Group("/token")

	tokenRoutes.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("These are the token routes.")
	})

	tokenRoutes.Post("/register", func(c *fiber.Ctx) error {
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

	tokenRoutes.Post("/login", func(c *fiber.Ctx) error {
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

		// Create the Claims
		claims := jwt.MapClaims{
			"email": user.Email,
			"id":    user.Id,
			"exp":   time.Now().Add(time.Hour * 24).Unix(),
		}

		// Create token
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		// Sign and get the complete encoded token as a string using the secret
		signedToken, err := token.SignedString([]byte(SecretKey))
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.JSON(fiber.Map{
				"message": "could not login",
			})
		}

		c.Status(fiber.StatusOK)
		return c.JSON(fiber.Map{
			"message": "logged in",
			"token":   signedToken,
		})
	})

	tokenRoutes.Get("/user", Protected(), func(c *fiber.Ctx) error {
		// c.Locals is set by the Protected() middleware & claims can be accessed from there
		user := c.Locals("user").(*jwt.Token) // Type assert to type of *jwt.Token
		claims := user.Claims.(jwt.MapClaims) // Type assert to type of claims
		email := claims["email"].(string)     // Type assert to string
		id := claims["id"].(float64)
		c.Status(fiber.StatusOK)
		return c.JSON(fiber.Map{
			"message": "authenticated",
			"name":    email,
			"id":      id,
		})
	})

}
