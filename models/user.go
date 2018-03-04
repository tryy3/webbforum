package models

import "time"

// User is the database model of a user
type User struct {
	ID           uint   `gorm:"primary_key"`
	Username     string `gorm:"unique"`
	FirstName    string
	LastName     string
	ProfileImage string // Should be an url
	Description  string

	// Permission
	Permission uint64
	GroupID    uint
	Group      Group

	// Auth
	Email    string `gorm:"unique"`
	Password string

	// Confirm
	ConfirmToken string
	Confirmed    bool

	// Lock
	AttemptNumber int64
	AttemptTime   time.Time
	Locked        time.Time

	// Recover
	RecoverToken       string
	RecoverTokenExpiry time.Time

	// Extra filled by gorm
	CreatedAt time.Time
	UpdatedAt time.Time
}
