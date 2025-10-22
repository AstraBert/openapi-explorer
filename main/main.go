package main

import (
	"log"
	"openapi-explorer/handlers"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := Setup()

	if err := app.Listen(":8000"); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}

func Setup() *fiber.App {
	app := fiber.New()
	app.Get("/", handlers.HomeRoute)
	app.Post("/code", handlers.HandleCodeGeneration)

	return app
}
