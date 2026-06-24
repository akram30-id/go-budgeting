package main

import (
	"time"

	"api-budgeting.smartcodex.cloud/config"
	"api-budgeting.smartcodex.cloud/controllers"
	cashcontroller "api-budgeting.smartcodex.cloud/controllers/cash_controller"
	treasurycontroller "api-budgeting.smartcodex.cloud/controllers/treasury_controller"
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
		&models.DeveloperInfo{},
		&models.Treasury{},
		&models.UserNotification{},
	)

	// throttling
	appLimiter := limiter.New(limiter.Config{
		Max:        1000,
		Expiration: 30 * time.Second,
	})

	// routes

	// TESTING ENDPOINT
	app.Post("/api-test", appLimiter, controllers.TestPush)

	// WEBOSCKET ENDPOINT HANDLER
	app.Get("/ws/:id", controllers.WSNotificationHandler)

	api := app.Group("/api")

	// api.Post("/register", appLimiter, controllers.Register)
	// api.Post("/login", appLimiter, controllers.Login)
	api.Post("/test-middleware", middleware.ApiAuth, controllers.TestPush)

	// debt := api.Group("/debt")
	// debt.Get("/vendors", middleware.ApiAuth, debtcontroller.ShowListVendors)

	treasury := api.Group("/treasury")
	treasury.Post("/duplicate", middleware.ApiAuth, treasurycontroller.DuplicateTreasury)
	treasury.Get("/members", middleware.ApiAuth, treasurycontroller.ListMembers)
	treasury.Get("/find-users", middleware.ApiAuth, treasurycontroller.FindUsers)
	treasury.Post("/invite-member", middleware.ApiAuth, treasurycontroller.InviteMember)
	treasury.Post("/member-access", middleware.ApiAuth, treasurycontroller.UpdateMemberAccess)
	treasury.Post("/remove-member", middleware.ApiAuth, treasurycontroller.RemoveMember)

	cash := api.Group("/cash")
	cash.Post("/sort-update", middleware.ApiAuth, cashcontroller.UpdateSortController)

	auth := api.Group("/auth")
	auth.Post("/change-password", middleware.ApiAuth, controllers.ChangePassword)

	// NOTIFICATION ENDPOINTS
	notification := api.Group("/notification")
	notification.Get("/current", middleware.ApiAuth, controllers.GetListNotification)
	notification.Post("/accept-invite", middleware.ApiAuth, treasurycontroller.AcceptInvitation)

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
