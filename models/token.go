package models

import "time"

// Token is the database model of a Token and used for the "remember me" functionality
type Token struct {
	// General information
	ID    uint `gorm:"primary_key"`
	User  string
	Token string

	// Time related information
	CreatedAt time.Time
	UpdatedAt time.Time
}
