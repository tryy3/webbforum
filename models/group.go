package models

import "time"

// Permission is the database model groups
type Group struct {
	// General information
	ID          uint `gorm:"primary_key"`
	Name        string
	Description string

	// Parent group
	ParentID uint
	Parent   *Group

	// Time related information
	CreatedAt time.Time
	UpdatedAt time.Time

	// Extra columns used when rendering to a template
	Permissions []ParsedPermission `gorm:"-"`
}
