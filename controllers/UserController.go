package controllers

import (
	"api-budgeting.smartcodex.cloud/helpers"
	"api-budgeting.smartcodex.cloud/models"
	"api-budgeting.smartcodex.cloud/services"
	"api-budgeting.smartcodex.cloud/validations"
	"github.com/gofiber/fiber/v2"
)

func Register(c *fiber.Ctx) error {

	var req validations.RegisterUserValidation

	body := c.Body()

	validateErr := helpers.ValidatePayload(body, &req)
	if validateErr != "" {
		return helpers.ErrorResponse(c, 422, validateErr)
	}

	register := models.ReqRegisterUser{
		Name:            req.Name,
		Email:           req.Email,
		Password:        req.Password,
		ConfirmPassword: req.ConfirmPassword,
		RoleId:          req.RoleId,
	}

	registerService := services.RegisterUser(register)
	if !registerService.Success {
		return helpers.ErrorResponse(c, registerService.HttpCode, registerService.Message)
	}

	return helpers.SuccessResponse(c, registerService.Data)

}
