package models

import "time"

type Category struct {
	ID          uint `gorm:"primary_key"`
	Name        string
	Description string

	// Extra filled by gorm
	CreatedAt time.Time
	UpdatedAt time.Time
}
