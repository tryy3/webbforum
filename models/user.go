package models

import (
	"html/template"
	"time"

	"github.com/volatiletech/authboss"
)

// User is the database model of a user
type User struct {
	// General information
	ID             uint   `gorm:"primary_key"`
	Username       string `gorm:"unique"`
	FirstName      string
	LastName       string
	ProfileImageID uint
	ProfileImage   *File `gorm:"association_autoupdate:false"`
	Description    string
	Attachments    []*File

	// Permission
	Permission uint64
	GroupID    uint
	Group      *Group

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

	// Time related information
	CreatedAt time.Time
	UpdatedAt time.Time

	// Template related columns
	ProfileImageURL   string        `gorm:"-"`
	ParsedDescription template.HTML `gorm:"-"`
}

// UserStorer is the interface required for taking care of a user
type UserStorer interface {
	Create(key string, attr authboss.Attributes) error
	Put(key string, attr authboss.Attributes) error
	Get(key string) (interface{}, error)
}
