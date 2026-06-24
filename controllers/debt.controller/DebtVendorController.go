package debtcontroller

import (
	"api-budgeting.smartcodex.cloud/helpers"
	"api-budgeting.smartcodex.cloud/services/debt"
	"github.com/gofiber/fiber/v2"
)

func ShowListVendors(c *fiber.Ctx) error {

	showVendorService := debt.ListVendors(c)

	if !showVendorService.Success {
		return helpers.ErrorResponse(c, showVendorService.HttpCode, showVendorService.Message)
	}

	return helpers.SuccessResponse(c, showVendorService.Data)

}
