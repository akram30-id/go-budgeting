package main

import (
	"api-budgeting.smartcodex.cloud/config"
	"api-budgeting.smartcodex.cloud/controllers"
	"api-budgeting.smartcodex.cloud/middleware"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	// connect DB
	config.ConnectDB()
	// config.DB.AutoMigrate(&models.Item{}, &models.Client{}, &models.Webhook{}, &models.QueueLog{}, &models.ConsumerLog{})
	config.DB.AutoMigrate()

	// routes

	// TESTING ENDPOINT
	app.Post("/api-test", controllers.TestPush)

	api := app.Group("/api")
	api.Post("/items", controllers.CreateItem).Use(middleware.ApiAuth)
	api.Get("/items", controllers.GetItems).Use(middleware.ApiAuth)
	api.Get("/items/:id", controllers.GetItem).Use(middleware.ApiAuth)
	api.Put("/items/:id", controllers.UpdateItem).Use(middleware.ApiAuth)
	api.Delete("/items/:id", controllers.DeleteItem).Use(middleware.ApiAuth)

	webhook := api.Group("/webhook")
	webhook.Post("/clients", controllers.RegisterClient).Use(middleware.ApiAuth)
	webhook.Post("/publish", controllers.PushToQueue).Use(middleware.ApiAuth)

	app.Listen(":3001")
}
