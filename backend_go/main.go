package main

import (
	"log"
	"os"

	"backend_go/db"
	"backend_go/handlers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	// initialize DB (reads env vars, provides sensible defaults)
	db.Init()

	app := fiber.New()
	app.Use(cors.New())

	api := app.Group("/api")
	api.Get("/posts", handlers.ListPosts)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Go API server started on :%s\n", port)
	log.Fatal(app.Listen(":" + port))
}
