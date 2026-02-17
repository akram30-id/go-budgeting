package treasurycontroller

import (
	"api-budgeting.smartcodex.cloud/helpers"
	"api-budgeting.smartcodex.cloud/models"
	"api-budgeting.smartcodex.cloud/services/treasury"
	treasuryvalidation "api-budgeting.smartcodex.cloud/validations/treasury_validation"
	"github.com/gofiber/fiber/v2"
)

func DuplicateTreasury(c *fiber.Ctx) error {

	var req treasuryvalidation.DuplicateTreasuryValidation

	if err := c.BodyParser(&req); err != nil {
		return helpers.ErrorResponse(c, 400, err.Error())
	}

	body := c.Body()

	validateErr := helpers.ValidatePayload(body, &req)
	if validateErr != "" {
		return helpers.ErrorResponse(c, 400, validateErr)
	}

	request := models.DuplicateTreasuryReq{
		TreasuryNo:       req.TreasuryNo,
		TreasuryDetailNo: req.TreasuryDetailNo,
		Month:            req.Month,
		Year:             req.Year,
	}

	if len(req.TreasuryDetailNo) == 0 {
		return helpers.ErrorResponse(c, 422, "Harus pilih salah satu cash.")
	}

	duplicateTreasuryService := treasury.DuplicateTreasuryService(request, c)

	if !duplicateTreasuryService.Success {
		return helpers.ErrorResponse(c, duplicateTreasuryService.HttpCode, duplicateTreasuryService.Message)
	}

	return helpers.SuccessResponse(c, map[string]any{
		"treasury_no": duplicateTreasuryService.Data,
	})

}
