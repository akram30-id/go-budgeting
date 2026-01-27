package middleware

import (
	"crypto/sha256"
	"encoding/hex"
	"regexp"
	"strings"

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
		return helpers.ErrorResponse(c, 500, "Failed to authorized this request. (BEARER_N_A)")
	}

	if !isBearer {
		return helpers.ErrorResponse(c, 401, "Invalid authorization format.")
	}

	cleanToken := RemoveBearer(apiKey)

	tokenParts := strings.Split(cleanToken, "|")
	if len(tokenParts) != 2 {
		return helpers.ErrorResponse(c, 401, "Invalid token authorization format.")
	}

	tokenId := tokenParts[0]
	plainToken := tokenParts[1]

	hash := sha256.Sum256([]byte(plainToken))
	hashedToken := hex.EncodeToString(hash[:])

	type TokenData struct {
		Token       string
		TokenableId int
	}

	var tokenData TokenData

	errGetToken := db.Raw(`
		SELECT a.token, a.tokenable_id
		FROM personal_access_tokens AS a
		WHERE a.token = ?
		AND a.id = ?
		-- AND a.created_at >= NOW() - INTERVAL 2 HOUR
		LIMIT 1
	`, hashedToken, tokenId).Scan(&tokenData).Error

	if errGetToken != nil {
		return helpers.ErrorResponse(c, 500, "Something's wrong when fetching access token.")
	}

	if tokenData.Token == "" {
		return helpers.ErrorResponse(c, 401, "Token is invalid.")
	}

	isUserExist := GetUserById(tokenData.TokenableId)

	if isUserExist.Success == false {
		if isUserExist.Code == 404 {
			return helpers.ErrorResponse(c, 404, "Invalid user service.")
		} else {
			return helpers.ErrorResponse(c, 500, "Internal server error (authentication).")
		}
	}

	c.Locals("userId", tokenData.TokenableId)
	c.Locals("name", isUserExist.Data.Name)
	c.Locals("email", isUserExist.Data.Email)
	c.Locals("roleId", isUserExist.Data.RoleId)

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

type GetUserReturn struct {
	Success bool
	Code    int
	Message string
	Data    AuthUser
}

func GetUserById(userId int) GetUserReturn {

	db := config.DB

	var authUser AuthUser

	getUser := db.Raw(`SELECT email, name, role_id FROM users WHERE id = ?`, userId).First(&authUser).Error

	if getUser != nil {
		return GetUserReturn{
			Success: false,
			Code:    500,
			Message: getUser.Error(),
		}
	}

	if authUser.Email == "" {
		return GetUserReturn{
			Success: false,
			Code:    404,
			Message: "No user found.",
		}
	}

	return GetUserReturn{
		Success: true,
		Code:    200,
		Message: "Success",
		Data:    authUser,
	}

}
