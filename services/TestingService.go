package services

import (
	"api-budgeting.smartcodex.cloud/config"
	"api-budgeting.smartcodex.cloud/models"

	"github.com/gofiber/fiber/v2"
)

func SaveTesting(c *fiber.Ctx, testModel models.WebhookTest) string {

	// 1. INSERT DB OF REQUEST TESTING
	if err := config.DB.Create(&testModel).Error; err != nil {
		return "Failed to create testing"
	}

	return ""

}
