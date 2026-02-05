package main

import (
	"time"

	"api-budgeting.smartcodex.cloud/config"
	"api-budgeting.smartcodex.cloud/controllers"
	cashcontroller "api-budgeting.smartcodex.cloud/controllers/cash_controller"
	"api-budgeting.smartcodex.cloud/middleware"
	"api-budgeting.smartcodex.cloud/models"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

func main() {
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOriginsFunc: AllowOriginCors,
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Authorization,Content-Type,X-API-KEY",
	}))

	// connect DB
	config.ConnectDB()
	// config.DB.AutoMigrate(&models.Item{}, &models.Client{}, &models.Webhook{}, &models.QueueLog{}, &models.ConsumerLog{})
	config.DB.AutoMigrate(
		&models.DebtVendor{},
		&models.DebtAccount{},
		&models.DebtVirtualAccount{},
		&models.DebtOutstanding{},
		&models.DebtPayment{},
	)

	// throttling
	appLimiter := limiter.New(limiter.Config{
		Max:        5,
		Expiration: 30 * time.Second,
	})

	// routes

	// TESTING ENDPOINT
	app.Post("/api-test", appLimiter, controllers.TestPush)

	api := app.Group("/api")

	// api.Post("/register", appLimiter, controllers.Register)
	// api.Post("/login", appLimiter, controllers.Login)
	api.Post("/test-middleware", middleware.ApiAuth, controllers.TestPush)

	// debt := api.Group("/debt")
	// debt.Get("/vendors", middleware.ApiAuth, debtcontroller.ShowListVendors)

	cash := api.Group("/cash")
	cash.Post("/sort-update", middleware.ApiAuth, cashcontroller.UpdateSortController)

	// api.Post("/items", controllers.CreateItem).Use(middleware.ApiAuth)
	// api.Get("/items", controllers.GetItems).Use(middleware.ApiAuth)
	// api.Get("/items/:id", controllers.GetItem).Use(middleware.ApiAuth)
	// api.Put("/items/:id", controllers.UpdateItem).Use(middleware.ApiAuth)
	// api.Delete("/items/:id", controllers.DeleteItem).Use(middleware.ApiAuth)

	// webhook := api.Group("/webhook")
	// webhook.Post("/clients", controllers.RegisterClient).Use(middleware.ApiAuth)
	// webhook.Post("/publish", controllers.PushToQueue).Use(middleware.ApiAuth)

	app.Listen(":3000")
}

func AllowOriginCors(origin string) bool {
	allowed := []string{
		"http://budgeting.test",
		"http://admin.budgeting.test",
		"https://budgeting.smartcodex.cloud",
		"http://202.10.47.104",
	}
	for _, o := range allowed {
		if o == origin {
			return true
		}
	}
	return false
}
