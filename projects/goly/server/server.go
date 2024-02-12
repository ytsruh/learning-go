package server

import (
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"ytsruh.com/goly/model"
	"ytsruh.com/goly/utils"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func redirect(c *fiber.Ctx) error {
	golyUrl := c.Params("redirect")
	goly, err := model.FindByGolyUrl(golyUrl)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "could not find goly in DB " + err.Error(),
		})
	}

	// Update the stats
	goly.Clicked += 1
	err = model.UpdateGoly(goly)
	if err != nil {
		log.Printf("error updating: %v\n", err)
	}
	// Check if string contains https or http otheriwse redirect doesn't work
	if !strings.Contains(goly.Redirect, "http") || !strings.Contains(goly.Redirect, "https") {
		goly.Redirect = "http://" + goly.Redirect
	}

	return c.Redirect(goly.Redirect, fiber.StatusTemporaryRedirect)
}

func getAllGolies(c *fiber.Ctx) error {
	golies, err := model.GetAllGolies()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "error getting all goly links " + err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(golies)
}

func getGoly(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "error could not parse id " + err.Error(),
		})
	}

	goly, err := model.GetGoly(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "error could not retreive goly from db " + err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(goly)
}

func createGoly(c *fiber.Ctx) error {
	c.Accepts("application/json")

	var goly model.Goly
	err := c.BodyParser(&goly)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "error parsing JSON " + err.Error(),
		})
	}

	if goly.Random {
		goly.Goly = utils.RandomURL(8)
	}

	err = model.CreateGoly(goly)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "could not create goly in db " + err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(goly)
}

func updateGoly(c *fiber.Ctx) error {
	c.Accepts("application/json")

	var goly model.Goly

	err := c.BodyParser(&goly)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "could not parse json " + err.Error(),
		})
	}

	err = model.UpdateGoly(goly)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "could not update goly link in DB " + err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(goly)
}

func deleteGoly(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "could not parse id from url " + err.Error(),
		})
	}

	err = model.DeleteGoly(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "could not delete from db " + err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "goly deleted.",
	})
}

func GetMostPopular(c *fiber.Ctx) error {
	golies, err := model.GetMostPopular()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "error getting most popular " + err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(golies)
}

func SetupAndListen() {
	port := os.Getenv("PORT")
	if port == "" {
		port = ":3000"
	} else {
		port = ":" + port
	}

	router := fiber.New(fiber.Config{
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
		IdleTimeout: 5 * time.Second,
	})

	// Set up middleware
	router.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))
	router.Use(RateLimiter)
	router.Use(Recover)
	router.Use(Compression)
	router.Use(Timer)
	router.Static("/", "./public")

	router.Use("/r", RedirectMiddleware)

	// Set up routes
	router.Get("/metrics", Monitor)
	router.Get("/r/:redirect", redirect)
	router.Get("/api/metrics", ApiMonitor)
	router.Get("/api/goly", getAllGolies)
	router.Get("/api/goly/:id", getGoly)
	router.Post("/api/goly", createGoly)
	router.Patch("/api/goly", updateGoly)
	router.Delete("/api/goly/:id", deleteGoly)
	router.Get("/api/popular", GetMostPopular)

	// Start server
	// Listen from a different goroutine
	go func() {
		if err := router.Listen(port); err != nil {
			log.Panic(err)
		}
	}()

	c := make(chan os.Signal, 1)                    // Create channel to signify a signal being sent
	signal.Notify(c, os.Interrupt, syscall.SIGTERM) // When an interrupt or termination signal is sent, notify the channel

	<-c // This blocks the main thread until an interrupt is received
	log.Println("Gracefully shutting down...")
	_ = router.Shutdown()

	log.Println("Running cleanup tasks...")

	// Your cleanup tasks go here
	log.Println("Fiber was successful shutdown.")

}
