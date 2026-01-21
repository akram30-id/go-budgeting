package models

import (
	"time"
)

type Client struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"size:100;null" json:"name"`
	Email     string    `gorm:"size:50;index:idx_client_email;null" json:"email"`
	ApiKey    string    `gorm:"size:255;uniqueIndex;null" json:"api_key"`
	IsActive  bool      `gorm:"default:true" json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
