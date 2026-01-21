package models

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type WebhookStatus string

const (
	StatusQueued   WebhookStatus = "QUEUED"
	StatusSent     WebhookStatus = "SENT"
	StatusSuccess  WebhookStatus = "SUCCESS"
	StatusFailed   WebhookStatus = "FAILED"
	StatusRetrying WebhookStatus = "RETRYING"
)

type Webhook struct {
	ID               uint           `gorm:"primaryKey" json:"id"`
	ClientID         uint           `gorm:"not null;index" json:"client_id"`
	TargetURL        string         `gorm:"type:text;not null" json:"target_url"`
	Headers          string         `gorm:"type:text" json:"headers"` // disimpan sebagai JSON string
	Payload          string         `gorm:"type:longtext" json:"payload"`
	Status           WebhookStatus  `gorm:"type:enum('QUEUED','SENT','SUCCESS','FAILED','RETRYING');default:'QUEUED'" json:"status"`
	RetryCount       int            `gorm:"default:0" json:"retry_count"`
	LastResponseCode int            `json:"last_response_code"`
	LastResponseBody string         `gorm:"type:longtext" json:"last_response_body"`
	ErrorMessage     string         `gorm:"type:text" json:"error_message"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`
}

type WebhookTest struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	ClientId  uint           `gorm:"index" json:"client_id"`
	Headers   datatypes.JSON `json:"headers"`
	Payload   datatypes.JSON `json:"payload"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}
