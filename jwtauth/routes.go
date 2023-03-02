package jwtauth

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

const SecretKey = "ewvfdjknl"

type CustomClaims struct {
	User                 string `json:"user"`
	jwt.RegisteredClaims `json:"claims"`
}

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

		claims := CustomClaims{
			user.Email,
			jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
				Issuer:    "jwt-test",
			},
		}

		// Create a new token object, specifying signing method and the claims
		// you would like it to contain.
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		// Sign and get the complete encoded token as a string using the secret
		tokenString, err := token.SignedString([]byte(SecretKey))
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.JSON(fiber.Map{
				"message": "could not login",
			})
		}

		cookie := fiber.Cookie{
			Name:     "jwt",
			Value:    tokenString,
			Expires:  time.Now().Add(time.Hour * 24),
			HTTPOnly: true, // Meant only for the server
		}

		c.Cookie(&cookie)

		return c.JSON(fiber.Map{
			"message": "success",
		})
	})

	app.Get("/user", func(c *fiber.Ctx) error {
		cookie := c.Cookies("jwt")

		token, err := jwt.ParseWithClaims(cookie, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(SecretKey), nil
		})

		if err != nil {
			log.Println(err)
			c.Status(fiber.StatusUnauthorized)
			return c.JSON(fiber.Map{
				"message": "unauthenticated",
			})
		}

		claims, ok := token.Claims.(*CustomClaims)
		if ok && token.Valid {
			c.Status(fiber.StatusOK)
			return c.JSON(fiber.Map{
				"message": "authenticated",
				"claims":  claims,
			})
		}

		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "unauthenticated",
		})

	})
}
