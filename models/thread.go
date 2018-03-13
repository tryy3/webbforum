package models

import "time"

type Thread struct {
	ID   uint `gorm:"primary_key"`
	Name string

	// CreatedBy
	CreatedByID uint
	CreatedBy   *User

	// Category
	CategoryID uint
	Category   *Category

	// Extra filled by gorm
	CreatedAt time.Time
	UpdatedAt time.Time
}
