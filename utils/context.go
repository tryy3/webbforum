package utils

import "github.com/jinzhu/gorm"

type Context struct {
	Database *gorm.DB
}
