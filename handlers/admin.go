package handlers

import (
	"net/http"

	"github.com/apex/log"
	"github.com/jinzhu/gorm"
	"github.com/tryy3/webbforum/models"
	"github.com/volatiletech/authboss"
)

func serveAdminPage(db *gorm.DB, data authboss.HTMLData) (authboss.HTMLData, string) {
	var categories []models.Category
	err := db.Find(&categories).Error
	if err != nil {
		log.WithError(err).Error("internal error when retrieving categories")
		return data, "internal error"
	}

	data = data.MergeKV("categories", categories)
	return data, ""
}

type AdminHandler struct {
	Database *gorm.DB
}

func (a AdminHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	data := LayoutData(w, r)
	data, err := serveAdminPage(a.Database, data)
	if err != "" {
		data = data.MergeKV("error", err)
	}

	mustRender(w, r, "admin", data)
}
