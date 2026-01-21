package controllers

import (
	"github.com/gofiber/fiber/v2"
)

func TestPush(c *fiber.Ctx) error {

	body := c.Body()

	header := c.GetReqHeaders()

	// var bodyMap map[string]any

	// bodyReq := bytes.NewBuffer(body)

	// if err := json.Unmarshal(body, &bodyMap); err != nil {
	// 	return c.Status(500).JSON(fiber.Map{
	// 		"status":  false,
	// 		"message": "Invalid JSON body",
	// 	})
	// }

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  true,
		"message": "success",
		"data": map[string]any{
			"header": header,
			"body":   string(body),
		},
	})
}
