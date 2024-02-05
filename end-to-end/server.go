package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/gorm"
	"homethings.ytsruh.com/controllers"
	"homethings.ytsruh.com/lib"
)

func main() {
	// Load environment variables
	envErr := godotenv.Load()
	if envErr != nil {
		log.Println("Error loading .env file")
	}
	// Set port
	port := os.Getenv("PORT")
	if port == "" {
		port = ":1323"
	} else {
		port = ":" + port
	}

	// Initialize Echo, set routes & database
	e := echo.New()
	e.Use(echo.MiddlewareFunc(middleware.CORS()))
	db := lib.InitDB()
	setRoutes(e, db)
	// Start server
	go func() {
		if err := e.Start(port); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()
	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	// Use a buffered channel to avoid missing signals as recommended for signal.Notify
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// Disconnect from database
	d, err := db.DB()
	if err != nil {
		log.Println(err)
		log.Println("Error disconnecting from database")
	}
	if err := d.Close(); err != nil {
		log.Println(err)
		log.Println("Error disconnecting from database")
	}
	log.Println("Database successfully disconnected")
	// Shutdown server
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
	log.Println("Server successfully shut down")
}

func setRoutes(e *echo.Echo, db *gorm.DB) {
	api := controllers.API{
		DB: db,
	}
	group := e.Group("/v1")
	// Public Routes
	group.POST("/login", api.Login)

	// Configure JWT middleware for Authentication
	group.Use(api.SetJWTAuth())

	// Profile routes
	group.GET("/profile", api.GetProfile)
	group.PATCH("/profile", api.PatchProfile)

	// Feedback route
	group.POST("/feedback", api.CreateFeedback)
}
