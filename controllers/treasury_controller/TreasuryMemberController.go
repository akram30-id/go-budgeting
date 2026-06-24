package treasurycontroller

import (
	"api-budgeting.smartcodex.cloud/helpers"
	"api-budgeting.smartcodex.cloud/models"
	"api-budgeting.smartcodex.cloud/services/treasury"
	treasuryvalidation "api-budgeting.smartcodex.cloud/validations/treasury_validation"
	"github.com/gofiber/fiber/v2"
)

func ListMembers(c *fiber.Ctx) error {
	var req treasuryvalidation.ListTreasuryMembersValidation

	if err := c.QueryParser(&req); err != nil {
		return helpers.ErrorResponse(c, 400, err.Error())
	}

	if req.TreasuryNo == "" {
		return helpers.ErrorResponse(c, 422, "Treasury tidak valid.")
	}

	if req.Page == 0 || req.Page < 0 {
		req.Page = 1
	}

	request := models.TreasuryMemberReq{
		TreasuryNo: req.TreasuryNo,
		Page:       req.Page,
	}

	listTreasuryMembersService := treasury.ListMembersService(&request, c)

	if !listTreasuryMembersService.Success {
		return helpers.ErrorResponse(c, listTreasuryMembersService.HttpCode, listTreasuryMembersService.Message)
	}

	return helpers.SuccessResponse(c, map[string]any{
		"treasury_no": req.TreasuryNo,
		"page":        req.Page,
		"data":        listTreasuryMembersService.Data,
	})
}

func FindUsers(c *fiber.Ctx) error {
	var req treasuryvalidation.FindMemberValidation

	if err := c.QueryParser(&req); err != nil {
		return helpers.ErrorResponse(c, 400, err.Error())
	}

	if req.Keywords == "" {
		return helpers.ErrorResponse(c, 422, "Treasury tidak minimal 3 karakter.")
	}

	request := req.Keywords

	findMemberService := treasury.FindMemberService(&request, c)

	if !findMemberService.Success {
		return helpers.ErrorResponse(c, findMemberService.HttpCode, findMemberService.Message)
	}

	return helpers.SuccessResponse(c, findMemberService.Data)
}

func InviteMember(c *fiber.Ctx) error {
	var req treasuryvalidation.InviteMemberValidation

	if err := c.BodyParser(&req); err != nil {
		return helpers.ErrorResponse(c, 400, err.Error())
	}

	body := c.Body()

	validateErr := helpers.ValidatePayload(body, &req)
	if validateErr != "" {
		return helpers.ErrorResponse(c, 400, validateErr)
	}

	request := models.InviteTreasuryMemberReq{
		TreasuryNo: req.TreasuryNo,
		Email:      req.Email,
	}

	inviteMemberService := treasury.InviteMemberService(&request, c)

	if !inviteMemberService.Success {
		return helpers.ErrorResponse(c, inviteMemberService.HttpCode, inviteMemberService.Message)
	}

	return helpers.SuccessResponse(c, inviteMemberService.Data)
}

func AcceptInvitation(c *fiber.Ctx) error {
	var req treasuryvalidation.AcceptInvitationValidation

	if err := c.BodyParser(&req); err != nil {
		return helpers.ErrorResponse(c, 400, err.Error())
	}

	body := c.Body()

	validateErr := helpers.ValidatePayload(body, &req)
	if validateErr != "" {
		return helpers.ErrorResponse(c, 400, validateErr)
	}

	request := models.AcceptInvitationReq{
		TreasuryNo:       req.TreasuryNo,
		NotificationCode: req.NotificationCode,
	}

	acceptInvitationService := treasury.AcceptInvitationService(&request, c)

	if !acceptInvitationService.Success {
		return helpers.ErrorResponse(c, acceptInvitationService.HttpCode, acceptInvitationService.Message)
	}

	return helpers.SuccessResponse(c, acceptInvitationService.Data)
}

func UpdateMemberAccess(c *fiber.Ctx) error {
	var req treasuryvalidation.MemberAccessValidation

	if err := c.BodyParser(&req); err != nil {
		return helpers.ErrorResponse(c, 400, err.Error())
	}

	body := c.Body()

	validateErr := helpers.ValidatePayload(body, &req)
	if validateErr != "" {
		return helpers.ErrorResponse(c, 422, validateErr)
	}

	request := models.UpdateMemberAccessReq{
		MemberId:   req.MemberId,
		TreasuryNo: req.TreasuryNo,
		CanEdit:    req.CanEdit,
	}

	updateMemberAccess := treasury.UpdateMemberAccess(&request, c)

	if !updateMemberAccess.Success {
		return helpers.ErrorResponse(c, updateMemberAccess.HttpCode, updateMemberAccess.Message)
	}

	return helpers.SuccessResponse(c, updateMemberAccess.Data)
}

func RemoveMember(c *fiber.Ctx) error {
	var req treasuryvalidation.RemoveMemberValidation

	if err := c.BodyParser(&req); err != nil {
		return helpers.ErrorResponse(c, 400, err.Error())
	}

	body := c.Body()

	validateErr := helpers.ValidatePayload(body, &req)
	if validateErr != "" {
		return helpers.ErrorResponse(c, 422, validateErr)
	}

	request := models.RemoveMemberReq{
		TreasuryNo: req.TreasuryNo,
		MemberId:   req.MemberId,
	}

	removeMemberTreasury := treasury.RemoveMemberTreasury(c, request)

	if !removeMemberTreasury.Success {
		return helpers.ErrorResponse(c, removeMemberTreasury.HttpCode, removeMemberTreasury.Message)
	}

	return helpers.SuccessResponse(c, removeMemberTreasury.Data)
}
