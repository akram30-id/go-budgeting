package helpers

import "github.com/gofiber/fiber/v2"

func SuccessResponse(c *fiber.Ctx, data interface{}) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  true,
		"message": "success",
		"data":    data,
	})
}

func ErrorResponse(c *fiber.Ctx, code int, msg string) error {
	return c.Status(code).JSON(fiber.Map{
		"status":  false,
		"message": msg,
	})
}

type ReturnService struct {
	Message  string
	Code     string
	Success  bool
	Data     map[string]any
	HttpCode int
}

func NewReturnService() ReturnService {
	return ReturnService{
		HttpCode: 200,
		Success:  true,
	}
}
