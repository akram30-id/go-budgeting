package models

import "time"

type DeveloperInfo struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	ParentId  int       `gorm:"index;default=null" json:"parent_id"`
	Title     string    `gorm:"default=null;size=100" json:"title"`
	Content   string    `gorm:"default=null;size=250" json:"content"`
	StartDate time.Time `gorm:"index;default=null" json:"start_date"`
	EndDate   time.Time `gorm:"index;default=null" json:"end_date"`
}
