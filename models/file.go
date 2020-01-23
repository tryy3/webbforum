package models

import "time"

// File is the database model for file content
type File struct {
	// General information
	ID            uint `primary_key`
	ContentType   string
	FileSizeBytes int64
	UploadName    string
	Base64Hash    string
	UserID        uint

	// Time related information
	CreatedAt time.Time
	UpdatedAt time.Time
}
