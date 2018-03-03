package utils

import "github.com/jinzhu/gorm"

type Context struct {
	Config   *Config
	Database *gorm.DB
}
