package services

import (
	"api-budgeting.smartcodex.cloud/config"
	"api-budgeting.smartcodex.cloud/models"
)

func CreateItem(item models.Item) (models.Item, error) {
	if err := config.DB.Create(&item).Error; err != nil {
		return item, err
	}

	return item, nil
}

func GetAllItems() ([]models.Item, error) {

	var items []models.Item

	if err := config.DB.Find(&items).Error; err != nil {
		return nil, err
	}

	return items, nil
}

func GetItemByID(id uint) (models.Item, error) {
	var item models.Item
	if err := config.DB.First(&item, id).Error; err != nil {
		return item, err
	}
	return item, nil
}

func UpdateItem(item models.Item) (models.Item, error) {
	if err := config.DB.Save(&item).Error; err != nil {
		return item, err
	}
	return item, nil
}

func DeleteItem(id uint) error {
	return config.DB.Delete(&models.Item{}, id).Error
}
