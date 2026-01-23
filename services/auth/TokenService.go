package auth

import (
	"crypto/md5"
	"encoding/hex"
	"strconv"
	"time"

	"api-budgeting.smartcodex.cloud/config"
	"api-budgeting.smartcodex.cloud/helpers"
)

func CreateToken(userId int) helpers.ReturnService {

	db := config.DB

	var id int

	err := db.Raw("SELECT id FROM users WHERE id = ?", userId).Scan(&id).Error

	if err != nil {
		return helpers.ReturnService{
			Message:  "Failed to fetch user data (TOKEN_FETCH).",
			Code:     "500G00",
			Success:  false,
			Data:     nil,
			HttpCode: 500,
		}
	}

	plainText := strconv.Itoa(id) + strconv.Itoa(time.Now().Nanosecond())

	// generate token
	hash := md5.New()
	hash.Write([]byte(plainText))

	hashedToken := hex.EncodeToString(hash.Sum(nil))

	if err := db.Exec("INSERT INTO personal_access_tokens (tokenable_type, tokenable_id, name, token, abilities, created_at) VALUES (?, ?, ?, ?, ?, ?)", "App\\Models\\User", userId, "api-token", hashedToken, "['*']", time.Now()).Error; err != nil {
		return helpers.ReturnService{
			Message:  err.Error(),
			Code:     "500G01",
			Success:  false,
			Data:     nil,
			HttpCode: 500,
		}
	}

	data := map[string]any{
		"user_id": userId,
		"token":   hashedToken,
	}

	return helpers.ReturnService{
		Message:  "Success",
		Code:     "200F01",
		Success:  true,
		Data:     data,
		HttpCode: 200,
	}

}
