package models

import "time"

// Token is the database model of a Token and used for the "remember me" functionality
type Token struct {
	ID    uint `gorm:"primary_key"`
	User  string
	Token string

	// Extra
	CreatedAt time.Time
	UpdatedAt time.Time
}
