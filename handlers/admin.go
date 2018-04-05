package handlers

import (
	"net/http"

	"github.com/apex/log"
	"github.com/jinzhu/gorm"
	"github.com/tryy3/webbforum/models"
	"github.com/volatiletech/authboss"
)

// serveAdminPage takes care of retriving the general data that admin pages need
func serveAdminPage(db *gorm.DB, data authboss.HTMLData) (authboss.HTMLData, string) {
	var categories []models.Category
	err := db.Find(&categories).Error
	if err != nil {
		log.WithError(err).Error("internal error when retrieving categories")
		return data, "internal error"
	}

	var groups []*models.Group
	err = db.Find(&groups).Error
	if err != nil {
		log.WithError(err).Error("internal error when retrieving categories")
		return data, "internal error"
	}

	// loops through all groups to get the permissions
	for _, g := range groups {
		parsed, err := models.GetGroupPermission(db, g)
		if err != nil {
			log.WithError(err).Error("internal error when retrieving categories")
			return data, "internal error"
		}
		g.Permissions = parsed
	}

	data = data.MergeKV("categories", categories)
	data = data.MergeKV("groups", groups)
	return data, ""
}

// AdminHandler shows admin pages
type AdminHandler struct {
	Database *gorm.DB
}

// ServeHTTP renders the admin page
func (a AdminHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	data := LayoutData(w, r)
	data, err := serveAdminPage(a.Database, data)
	if err != "" {
		data = data.MergeKV("error", err)
	}

	mustRender(w, r, "admin", data)
}
