package middleware

import (
	"fmt"
	"regexp"

	"api-budgeting.smartcodex.cloud/config"
	"api-budgeting.smartcodex.cloud/helpers"

	"github.com/gofiber/fiber/v2"
)

func ApiAuth(c *fiber.Ctx) error {

	db := config.DB

	apiKey := c.Get("Authorization")

	if apiKey == "" {
		return helpers.ErrorResponse(c, 401, "Authorization is required.")
	}

	isBearer, err := regexp.MatchString(`(?i)^Bearer`, apiKey)

	if err != nil {
		return helpers.ErrorResponse(c, 500, "Failed to authorized this request.")
	}

	if !isBearer {
		return helpers.ErrorResponse(c, 401, "Invalid authorization format.")
	}

	cleanToken := RemoveBearer(apiKey)

	type TokenData struct {
		Token  string
		UserId int
	}

	var tokenData TokenData

	errGetToken := db.Raw(`
		SELECT a.token, a.tokenable_id
		FROM personal_access_tokens AS a
		WHERE a.token = ?
		AND a.created_at >= NOW() - INTERVAL 2 HOUR
	`, cleanToken).Scan(&tokenData).Error

	if errGetToken != nil {
		return helpers.ErrorResponse(c, 500, "Something's wrong when fetching access token.")
	}

	if tokenData.Token == "" {
		return helpers.ErrorResponse(c, 401, "Token is invalid.")
	}

	fmt.Printf("User id: %d", tokenData.UserId)

	isUserExist := GetUserById(tokenData.UserId)

	if isUserExist["success"] == false {
		if isUserExist["code"] == 404 {
			return helpers.ErrorResponse(c, 404, "Invalid user service.")
		} else {
			return helpers.ErrorResponse(c, 500, "Internal server error (authentication).")
		}
	}

	c.Set("token", tokenData.Token)
	c.Locals("userId", tokenData.UserId)

	return c.Next()

}

func RemoveBearer(token string) string {

	re := regexp.MustCompile(`(?i)^Bearer\s+`)

	return re.ReplaceAllString(token, "")

}

type AuthUser struct {
	Email  string
	Name   string
	RoleId int
}

func GetUserById(userId int) map[string]any {

	db := config.DB

	var authUser AuthUser

	getUser := db.Raw(`SELECT email, name, role_id FROM users WHERE id = ?`, userId).Scan(&authUser).Error

	if getUser != nil {
		return map[string]any{
			"success": false,
			"code":    500,
			"message": getUser.Error(),
		}
	}

	if authUser.Email == "" {
		fmt.Printf("SELECT email, name, role_id FROM users WHERE id = %d", userId)

		return map[string]any{
			"success": false,
			"code":    404,
			"message": "No user found.",
		}
	}

	return map[string]any{
		"success": true,
		"code":    200,
		"data":    authUser,
	}

}
