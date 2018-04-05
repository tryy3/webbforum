package models

import "time"

type Thread struct {
	// General information
	ID   uint `gorm:"primary_key"`
	Name string

	// CreatedBy
	CreatedByID uint
	CreatedBy   *User

	// Category
	CategoryID uint
	Category   *Category

	// Time related information
	CreatedAt time.Time
	UpdatedAt time.Time

	// Template related columns
	DisplayName string    `gorm:"-"`
	LatestPost  time.Time `gorm:"-"`
	CountPost   int64     `gorm:"-"`
}
