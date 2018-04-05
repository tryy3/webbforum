package models

import "time"

type Category struct {
	// General information
	ID          uint `gorm:"primary_key"`
	Name        string
	Description string

	// Time related information
	CreatedAt time.Time
	UpdatedAt time.Time

	// Extra columns used when rendering to a template
	DisplayName  string `gorm:"-"`
	CountPost    int64  `gorm:"-"`
	CountThreads int64  `gorm:"-"`
	LatestUpdate *Post  `gorm:"-"`
}
