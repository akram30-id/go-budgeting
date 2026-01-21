package controllers

import (
	"api-budgeting.smartcodex.cloud/helpers"
	"api-budgeting.smartcodex.cloud/models"
	"api-budgeting.smartcodex.cloud/services"
	"api-budgeting.smartcodex.cloud/validations"

	"github.com/gofiber/fiber/v2"
)

func RegisterClient(c *fiber.Ctx) error {

	var req validations.RegisterClientValidation
	body := c.Body()

	validateErr := helpers.ValidatePayload(body, &req)
	if validateErr != "" {
		return helpers.ErrorResponse(c, 400, validateErr)
	}

	if errr := validate.Struct(&req); errr != nil {
		return helpers.ErrorResponse(c, 400, errr.Error())
	}

	client := models.Client{
		Name:  req.Name,
		Email: req.Email,
	}

	created, errInsert := services.RegisterClient(client)

	if errInsert != "" {
		return helpers.ErrorResponse(c, 500, errInsert)
	}

	return helpers.SuccessResponse(c, created)

}
