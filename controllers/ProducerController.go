package controllers

import (
	"encoding/json"
	"regexp"

	"api-budgeting.smartcodex.cloud/helpers"
	"api-budgeting.smartcodex.cloud/models"
	"api-budgeting.smartcodex.cloud/services"
	"api-budgeting.smartcodex.cloud/validations"

	"github.com/gofiber/fiber/v2"
)

func PushToQueue(c *fiber.Ctx) error {

	var req validations.RequestQueueValidation

	body := c.Body()

	validateErr := helpers.ValidatePayload(body, &req)
	if validateErr != "" {
		return helpers.ErrorResponse(c, 422, validateErr)
	}

	if err := validate.Struct(&req); err != nil {
		return helpers.ErrorResponse(c, 400, err.Error())
	}

	matchedUrlFormat, _ := regexp.MatchString(`(http|https)://`, req.TargetUrl)
	if !matchedUrlFormat {
		return helpers.ErrorResponse(c, 400, "Invalid targetUrl format (must be contain prefix http or https).")
	}

	availableMethod := []string{
		"POST",
		"GET",
		"DELETE",
		"PUT",
		"PATCH",
	}

	if !helpers.Contains(availableMethod, req.Method) {
		return helpers.ErrorResponse(c, 400, "Invalid httpMethod")
	}

	clientId := c.Locals("clientId").(uint)

	dataHeader, _ := json.Marshal(req.Headers)

	dataBody, _ := json.Marshal(req.Body)

	queueData := models.QueuePush{
		Client:    clientId,
		TargetUrl: req.TargetUrl,
		Method:    req.Method,
		Headers:   dataHeader,
		Body:      dataBody,
	}

	publishQueue := services.CreateQueue(queueData)
	if publishQueue["error"] != "" {
		return helpers.ErrorResponse(c, 500, publishQueue["error"])
	}

	return helpers.SuccessResponse(c, map[string]string{
		"success": "ok",
	})

}
