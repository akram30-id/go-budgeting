package controllers

import (
	"strconv"

	"api-budgeting.smartcodex.cloud/helpers"
	"api-budgeting.smartcodex.cloud/models"
	"api-budgeting.smartcodex.cloud/services"
	"api-budgeting.smartcodex.cloud/validations"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var validate = validator.New()

func CreateItem(c *fiber.Ctx) error {

	var req validations.CreateItemRequest
	body := c.Body()

	validateErr := helpers.ValidatePayload(body, &req)
	if validateErr != "" {
		return helpers.ErrorResponse(c, fiber.StatusBadRequest, validateErr)
	}

	if err := validate.Struct(&req); err != nil {
		return helpers.ErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	item := models.Item{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
	}

	created, err := services.CreateItem(item)
	if err != nil {
		return helpers.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return helpers.SuccessResponse(c, created)
}

func GetItems(c *fiber.Ctx) error {
	items, err := services.GetAllItems()
	if err != nil {
		return helpers.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return helpers.SuccessResponse(c, items)
}

func GetItem(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	item, err := services.GetItemByID(uint(id))
	if err != nil {
		return helpers.ErrorResponse(c, fiber.StatusNotFound, "Item not found")
	}
	return helpers.SuccessResponse(c, item)
}

func UpdateItem(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	var req validations.UpdateItemRequest
	if err := c.BodyParser(&req); err != nil {
		return helpers.ErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}
	if err := validate.Struct(&req); err != nil {
		return helpers.ErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	item, err := services.GetItemByID(uint(id))
	if err != nil {
		return helpers.ErrorResponse(c, fiber.StatusNotFound, "Item not found")
	}

	if req.Name != "" {
		item.Name = req.Name
	}
	if req.Description != "" {
		item.Description = req.Description
	}
	if req.Price > 0 {
		item.Price = req.Price
	}

	updated, err := services.UpdateItem(item)
	if err != nil {
		return helpers.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}
	return helpers.SuccessResponse(c, updated)
}

func DeleteItem(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	if err := services.DeleteItem(uint(id)); err != nil {
		return helpers.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}
	return helpers.SuccessResponse(c, fiber.Map{"deleted": true})
}
