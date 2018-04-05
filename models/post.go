package models

import (
	"html/template"
	"time"
)

// Post is the database model for posts
type Post struct {
	// General information
	ID      uint `gorm:"primary_key"`
	Comment string

	// User
	UserID uint
	User   *User

	// Thread
	ThreadID uint
	Thread   *Thread

	// Time related information
	CreatedAt time.Time
	UpdatedAt time.Time

	// Extra columns used when rendering to a template
	DisplayComment template.HTML `gorm:"-"`
}
