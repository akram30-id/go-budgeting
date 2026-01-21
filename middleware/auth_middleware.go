package middleware

import (
	"api-budgeting.smartcodex.cloud/config"
	"api-budgeting.smartcodex.cloud/helpers"
	"api-budgeting.smartcodex.cloud/models"

	"github.com/gofiber/fiber/v2"
)

func ApiAuth(c *fiber.Ctx) error {

	apiKey := c.Get("X-API-KEY")

	if apiKey == "" {
		return helpers.ErrorResponse(c, 401, "X-API-KEY is required.")
	}

	var clientModel models.Client
	result := config.DB.Table("clients").Where(map[string]string{
		"api_key": apiKey,
	}).First(&clientModel)

	if result.Error != nil {
		return helpers.ErrorResponse(c, 401, "API Key is invalid.")
	}

	c.Set("apiKey", apiKey)
	c.Locals("clientId", clientModel.ID)

	return c.Next()

}
