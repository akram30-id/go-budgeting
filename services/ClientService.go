package services

import (
	"crypto/sha256"
	"encoding/hex"

	"api-budgeting.smartcodex.cloud/config"
	"api-budgeting.smartcodex.cloud/models"
)

func RegisterClient(client models.Client) (models.Client, string) {

	if err := config.DB.Where("email = ?", client.Email).First(&client).Error; err == nil {
		return client, "Client is already registered"
	}

	hash := sha256.New()
	hash.Write([]byte(client.Email))

	client.ApiKey = hex.EncodeToString(hash.Sum(nil))

	if err := config.DB.Create(&client).Error; err != nil {

		return client, err.Error()

	}

	return client, ""

}
