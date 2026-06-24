package notificationrepo

import (
	"api-budgeting.smartcodex.cloud/config"
)

func GetTotalUnreadNotification(userId uint) int {

	db := config.DB

	var notifId []uint

	errGetTotalUnread := db.Table("user_notifications").
		Select("id").
		Where("user_id = ?", userId).
		Where("state = ?", 1).
		Limit(100).Scan(&notifId).Error

	if errGetTotalUnread != nil {
		return 0
	}

	return len(notifId)
}
