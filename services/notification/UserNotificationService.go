package notification

import (
	"strconv"
	"time"

	"api-budgeting.smartcodex.cloud/config"
	"api-budgeting.smartcodex.cloud/helpers"
	"api-budgeting.smartcodex.cloud/models"
	notificationrepo "api-budgeting.smartcodex.cloud/repository/notification_repo"
	"github.com/gofiber/fiber/v2"
)

func CreateNotification(data models.CreateNotification) helpers.ReturnService {

	db := config.DB

	userNotification := models.UserNotification{
		UserId:       data.UserId,
		Title:        data.Title,
		Message:      data.Message,
		State:        1,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		UserSenderId: data.UserSenderId,
	}

	tx := db.Begin()

	errCreateNotification := tx.Table("user_notifications").Create(&userNotification).Error

	if errCreateNotification != nil {
		return helpers.ReturnService{
			Message:  "Gagal saat menyimpan notifikasi.",
			Code:     "500G00",
			Success:  false,
			HttpCode: 500,
		}
	}

	lastNotifID := userNotification.ID
	if lastNotifID == 0 {
		tx.Rollback()

		return helpers.ReturnService{
			Message:  "Gagal saat menyimpan notifikasi.",
			Code:     "500G01",
			Success:  false,
			HttpCode: 500,
		}
	}

	var generatedNotifCode string
	generatedNotifCode = "NTF" + helpers.StrPadLeft(strconv.Itoa(int(lastNotifID)), 10, "0")

	message := data.Message

	if message == "" {
		message = `<div class="col-12 list-notif rounded-2 mb-1" style="height: 128px; width: 380px;">
					<div class="row align-items-center">
						<div class="col-2 mt-3">
							<i class="bi bi-person-fill" style="font-size: 3rem;"></i>
						</div>
						<div class="col-10">
							` + data.OwnerEmail + ` mengundang kamu ke Treasury <b> ` + data.TreasuryNo + `</b>
						</div>
						<div class="col-2"></div>
						<div class="col-10">
							<div class="d-flex align-items-center">
								<button class="btn btn-sm btn-primary fs-6 btn-notif-accept" data='{"treasury_no":"` + data.TreasuryNo + `","notification_code":"` + generatedNotifCode + `"}' style="margin-right: 10px;">Accept</button>
								<button class="btn btn-sm btn-outline-danger fs-6">Delete</button>
							</div>
						</div>
					</div>
				</div>`
	}

	errUpdateNotifCode := tx.Table("user_notifications").
		Where("id = ?", userNotification.ID).
		Update("notification_code", generatedNotifCode).
		Update("message", message).Error

	if errUpdateNotifCode != nil {
		return helpers.ReturnService{
			Message:  "Gagal saat menyimpan notifikasi.",
			Code:     "500G02",
			Success:  false,
			HttpCode: 500,
		}
	}

	if errDBTransaction := tx.Commit().Error; errDBTransaction != nil {
		return helpers.ReturnService{
			Message:  "Something's wrong. Please try again later.",
			Code:     "502E00",
			Success:  false,
			HttpCode: 502,
		}
	}

	return helpers.ReturnService{
		Message:  "Sukses.",
		Code:     "200S00",
		Success:  true,
		HttpCode: 200,
		Data:     generatedNotifCode,
	}

}

func ListNotification(c *fiber.Ctx, limitPage models.ReqListNotification) helpers.ReturnService {

	db := config.DB

	userId := c.Locals("userId").(int)

	type UserNotified struct {
		NotificationCode string
		title            string
		Message          string
		SenderName       string
		CreatedAt        string
	}

	var userNotified []UserNotified

	if limitPage.Limit == 0 {
		limitPage.Limit = 1
	}

	if limitPage.Page <= 1 {
		limitPage.Page = 0
	} else {
		limitPage.Page = (limitPage.Page - 1) * limitPage.Limit
	}

	errGetNotif := db.Debug().Table("user_notifications AS a").
		Joins("JOIN users AS b ON a.user_sender_id=b.id").
		Select("a.notification_code, a.title, a.message, b.name AS sender_name, a.created_at").
		Where("a.user_id = ?", userId).
		Where("a.state = ?", 1).
		Order("a.id DESC").
		Limit(limitPage.Limit).Offset(limitPage.Page).Scan(&userNotified).Error

	if errGetNotif != nil {
		return helpers.ReturnService{
			Message:  "Gagal saat mengambil notifikasi.",
			Code:     "502G0",
			Success:  false,
			HttpCode: 502,
		}
	}

	type DataResponse struct {
		NotificationList []UserNotified
		Count            int
	}

	getUnreadCount := notificationrepo.GetTotalUnreadNotification(uint(userId))

	response := DataResponse{
		NotificationList: userNotified,
		Count:            getUnreadCount,
	}

	return helpers.ReturnService{
		Message:  "Success.",
		Code:     "200S01",
		Success:  true,
		HttpCode: 200,
		Data:     response,
	}

}
