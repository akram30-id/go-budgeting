package controllers

import (
	"api-budgeting.smartcodex.cloud/config/socket"
	"api-budgeting.smartcodex.cloud/helpers"
	"api-budgeting.smartcodex.cloud/models"
	"api-budgeting.smartcodex.cloud/services/notification"
	treasuryvalidation "api-budgeting.smartcodex.cloud/validations/treasury_validation"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

// controllers/notification_controller.go
func WSNotificationHandler(c *fiber.Ctx) error {
	userId := c.Params("id") // Misal dari URL /ws/:id

	return websocket.New(func(conn *websocket.Conn) {
		// 1. Daftar ke Hub saat konek
		socket.GlobalHub.Register(userId, conn)

		// 2. Hapus dari Hub saat diskonek
		defer func() {
			socket.GlobalHub.Unregister(userId)
			conn.Close()
		}()

		// 3. Keep-alive loop
		for {
			if _, _, err := conn.ReadMessage(); err != nil {
				break
			}
		}
	})(c)
}

func GetListNotification(c *fiber.Ctx) error {

	var req treasuryvalidation.ListNotificationUser

	if err := c.QueryParser(&req); err != nil {
		return helpers.ErrorResponse(c, 400, err.Error())
	}

	request := models.ReqListNotification{
		Limit: req.Limit,
		Page:  req.Page,
	}

	getUserNotifService := notification.ListNotification(c, request)

	if !getUserNotifService.Success {
		return helpers.ErrorResponse(c, getUserNotifService.HttpCode, getUserNotifService.Message)
	}

	return helpers.SuccessResponse(c, getUserNotifService.Data)
}
