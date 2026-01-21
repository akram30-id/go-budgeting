package models

import (
	"time"

	"gorm.io/datatypes"
)

type QueueLog struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Client    uint           `gorm:"index;default:null" json:"client_id"`
	TargetUrl string         `gorm:"index;size:128;default:null" json:"target_url"`
	Method    string         `gorm:"index;size:16;default:null" json:"http_method"`
	Headers   datatypes.JSON `gorm:"default:null" json:"header"`
	Body      datatypes.JSON `gorm:"default:null" json:"body"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

type QueuePush struct {
	Client    uint           `gorm:"index;default:null" json:"client_id"`
	TargetUrl string         `gorm:"index;size:128;default:null" json:"target_url"`
	Method    string         `gorm:"index;size:16;default:null" json:"http_method"`
	Headers   datatypes.JSON `gorm:"default:null" json:"header"`
	Body      datatypes.JSON `gorm:"default:null" json:"body"`
}

type ConsumerLog struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	Client         uint           `gorm:"index;default:null" json:"client_id"`
	TargetUrl      string         `gorm:"index;size:128;default:null" json:"target_url"`
	Method         string         `gorm:"index;size:16;default:null" json:"http_method"`
	RequestHeaders datatypes.JSON `gorm:"default:null" json:"request_header"`
	RequestBody    datatypes.JSON `gorm:"default:null" json:"request_body"`
	ResponseCode   uint           `gorm:"index,default:null" json:"response_code"`
	ResponseBody   datatypes.JSON `gorm:"index,default:null" json:"response_body"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
}
