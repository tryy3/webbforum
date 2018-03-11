package models

import "time"

type File struct {
	ID            uint `primary_key`
	ContentType   string
	FileSizeBytes int64
	UploadName    string
	Base64Hash    string
	UserID        uint

	// Extra filled by gorm
	CreatedAt time.Time
	UpdatedAt time.Time
}
