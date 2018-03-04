package models

import "time"

type Post struct {
	ID      uint `gorm:"primary_key"`
	Name    string
	Comment string

	// User
	UserID uint
	User   User

	// Thread
	ThreadID uint
	Thread   Thread

	// Extra filled by gorm
	CreatedAt time.Time
	UpdatedAt time.Time
}
