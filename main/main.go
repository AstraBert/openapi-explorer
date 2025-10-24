package main

import (
	"log"
	"openapi-explorer/handlers"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/storage/sqlite3"
)

func main() {
	app := Setup()

	if err := app.Listen(":8000"); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}

func Setup() *fiber.App {
	app := fiber.New()
	corsConfig := cors.Config{
		AllowOrigins: "https://openapi.clelias-dockploy.my.id",
		AllowMethods: "POST",
	}

	storage := sqlite3.New()
	limiterConfig := limiter.Config{
		Max: 5,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP() // Track limit per IP address
		},
		Expiration: 1 * time.Minute,
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"message": "You reached the maximum number of requests per minute, please retry soon!",
			})
		},
		Storage: storage,
	}
	app.Get("/", handlers.HomeRoute)
	app.Post("/code/generate", cors.New(corsConfig), limiter.New(limiterConfig), handlers.HandleCodeGeneration)

	return app
}
