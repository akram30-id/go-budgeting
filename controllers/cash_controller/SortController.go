package cashcontroller

import (
	"api-budgeting.smartcodex.cloud/helpers"
	"api-budgeting.smartcodex.cloud/models"
	"api-budgeting.smartcodex.cloud/services/cash"
	cashvalidation "api-budgeting.smartcodex.cloud/validations/cash_validation"
	"github.com/gofiber/fiber/v2"
)

func UpdateSortController(c *fiber.Ctx) error {

	var req cashvalidation.UpdateSortValidation

	body := c.Body()

	validateErr := helpers.ValidatePayload(body, &req)
	if validateErr != "" {
		return helpers.ErrorResponse(c, 400, validateErr)
	}

	request := models.SortUpdate{
		TreasuryDetailNo: req.TreasuryDetailNo,
		Sorts:            req.Sorts,
	}

	updateSortService := cash.UpdateSortCash(request)

	if !updateSortService.Success {
		return helpers.ErrorResponse(c, updateSortService.HttpCode, updateSortService.Message)
	}

	return helpers.SuccessResponse(c, updateSortService.Message)

}
