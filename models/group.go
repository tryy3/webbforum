package models

import "time"

// Permission is the database model groups
type Group struct {
	ID         uint `gorm:"primary_key"`
	Name       string
	Permission uint64

	// Parent
	ParentID uint
	Parent   *Group

	// Extra filled by gorm
	CreatedAt time.Time
	UpdatedAt time.Time
}
