package models

type Treasury struct {
	InSorting int8 `gorm:"default=0" json:"in_sorting"`
}
