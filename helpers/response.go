package helpers

import (
	"api-budgeting.smartcodex.cloud/models"
	"github.com/gofiber/fiber/v2"
)

func SuccessResponse(c *fiber.Ctx, data interface{}) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "success",
		"data":    data,
	})
}

func ErrorResponse(c *fiber.Ctx, code int, msg string) error {
	return c.Status(code).JSON(fiber.Map{
		"success": false,
		"message": msg,
	})
}

type ReturnService struct {
	Message  string
	Code     string
	Success  bool
	Data     any
	HttpCode int
}

func NewReturnService() ReturnService {
	return ReturnService{
		HttpCode: 200,
		Success:  true,
	}
}

type LoginResponse struct {
	Success bool                    `json:"success"`
	Message string                  `json:"message"`
	Token   string                  `json:"token"`
	User    models.UserLoginSuccess `json:"user"`
}
