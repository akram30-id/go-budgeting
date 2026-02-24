package controllers

import (
	"api-budgeting.smartcodex.cloud/helpers"
	"api-budgeting.smartcodex.cloud/models"
	"api-budgeting.smartcodex.cloud/services"
	"api-budgeting.smartcodex.cloud/services/auth"
	"api-budgeting.smartcodex.cloud/validations"
	uservalidation "api-budgeting.smartcodex.cloud/validations/user_validation"
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

func Login(c *fiber.Ctx) error {

	var reqValidate validations.LoginUserValidation

	body := c.Body()

	err := helpers.ValidatePayload(body, &reqValidate)
	if err != "" {
		return helpers.ErrorResponse(c, 422, err)
	}

	req := models.ReqUserLogin{
		Email:    reqValidate.Email,
		Password: reqValidate.Password,
	}

	login := services.Login(req)

	if !login.Success {
		return helpers.ErrorResponse(c, 401, login.Message)
	}

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": login.Message,
		"token":   login.Token,
		"user":    login.User,
	})

}

func ChangePassword(c *fiber.Ctx) error {
	var req uservalidation.ChangePasswordValidation

	if err := c.BodyParser(&req); err != nil {
		return helpers.ErrorResponse(c, 400, err.Error())
	}

	body := c.Body()

	validateErr := helpers.ValidatePayload(body, &req)
	if validateErr != "" {
		return helpers.ErrorResponse(c, 400, validateErr)
	}

	request := models.ChangePasswordRequest{
		OldPassword:     req.OldPassword,
		NewPassword:     req.NewPassword,
		ConfirmPassword: req.ConfirmPassword,
	}

	changePasswordService := auth.ChangePasswordService(request, c)

	if !changePasswordService.Success {
		return helpers.ErrorResponse(c, changePasswordService.HttpCode, changePasswordService.Message)
	}

	return helpers.SuccessResponse(c, changePasswordService.Data)
}
