package models

import "time"

type UserNotification struct {
	ID               uint      `gorm:"primaryKey" json:"id"`
	NotificationCode string    `gorm:"index;default=null;size=50" json:"vendor_code"`
	UserId           int       `gorm:"index;default=null" json:"user_id"`
	UserSenderId     int       `gorm:"index;default=null" json:"user_sender_id"`
	Title            string    `gorm:"default=null;size=150" json:"vendor_name"`
	Message          string    `gorm:"default=null;size=200" json:"description"`
	State            int       `gorm:"index;default=1" json:"state"`
	CreatedAt        time.Time `gorm:"index;default=null" json:"created_at"`
	UpdatedAt        time.Time `gorm:"default=null" json:"updated_at"`
}

type CreateNotification struct {
	UserId       int
	UserSenderId int
	OwnerEmail   string
	TreasuryNo   string
	Title        string
	Message      string
	CreatedAt    time.Time
	updatedAt    time.Time
}

type ReqListNotification struct {
	Limit int
	Page  int
}
