package handlers

import (
	"net/http"
	"strconv"

	"github.com/apex/log"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/tryy3/webbforum/models"
	"github.com/volatiletech/authboss"
)

// CategoryCreateHandler creates a handler for creating categories
type CategoryCreateHandler struct {
	Database *gorm.DB
}

// ServeHTTP handle creating new categories
func (c CategoryCreateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	data := LayoutData(w, r)

	attr, err := authboss.AttributesFromRequest(r)
	if err != nil {
		log.WithError(err).Error("internal error when parsing request")
		data = data.MergeKV("error", "internal error")
		data, _ = serveAdminPage(c.Database, data)
		mustRender(w, r, "admin", data)
		return
	}

	var category models.Category
	if name, ok := attr.String("category_name"); ok {
		category.Name = name
	}
	if description, ok := attr.String("category_description"); ok {
		category.Description = description
	}
	data = data.MergeKV("category", category)

	if category.Name == "" {
		data = data.MergeKV("errs", map[string][]string{"category_name": {"Missing category name"}})
		data, errStr := serveAdminPage(c.Database, data)
		if errStr != "" {
			data = data.MergeKV("error", errStr)
		}
		mustRender(w, r, "admin", data)
		return
	}

	err = c.Database.Create(&category).Error
	if err != nil {
		log.WithError(err).Error("internal error when creating new category")
		data = data.MergeKV("error", "internal error")
		data, _ = serveAdminPage(c.Database, data)
		mustRender(w, r, "admin", data)
		return
	}

	http.Redirect(w, r, "/admin", http.StatusFound)
}

// CategoryEditHandler edits the category in the database
type CategoryEditHandler struct {
	Database *gorm.DB
}

// ServeHTTP handles editing category.
func (c CategoryEditHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	data := LayoutData(w, r)

	attr, err := authboss.AttributesFromRequest(r)
	if err != nil {
		log.WithError(err).Error("internal error when parsing request")
		data = data.MergeKV("error", "internal error")
		data, _ = serveAdminPage(c.Database, data)
		mustRender(w, r, "admin", data)
		return
	}

	var category models.Category
	id, errStr := getCategoryID(attr)
	if errStr != "" {
		data, serveErr := serveAdminPage(c.Database, data)
		if serveErr != "" {
			errStr += "<br>" + serveErr
		}
		data = data.MergeKV("error", errStr)
		mustRender(w, r, "admin", data)
		return
	}
	category.ID = id

	result := c.Database.First(&category)
	if result.Error != nil || result.RecordNotFound() {
		errStr := "tried to modify invalid category"
		data, serveErr := serveAdminPage(c.Database, data)
		if serveErr != "" {
			errStr += "<br>" + serveErr
		}
		data = data.MergeKV("error", errStr)
		mustRender(w, r, "admin", data)
		return
	}

	if name, ok := attr.String("category_name"); ok {
		category.Name = name
	}
	if description, ok := attr.String("category_description"); ok {
		category.Description = description
	}
	data = data.MergeKV("category_edit", category)

	if category.Name == "" {
		data = data.MergeKV("category_errs", map[string][]string{"category_name": {"Missing category name"}})
		data, errStr := serveAdminPage(c.Database, data)
		if errStr != "" {
			data = data.MergeKV("error", errStr)
		}
		mustRender(w, r, "admin", data)
		return
	}

	err = c.Database.Model(&category).Updates(&category).Error
	if err != nil {
		log.WithError(err).Error("internal error when updating a category")
		data = data.MergeKV("error", "internal error")
		data, _ = serveAdminPage(c.Database, data)
		mustRender(w, r, "admin", data)
		return
	}

	http.Redirect(w, r, "/admin", http.StatusFound)
}

// CategoryDeleteHandler removes category from database
type CategoryDeleteHandler struct {
	Database *gorm.DB
}

// ServeHTTP handles deleting the category.
func (c CategoryDeleteHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	data := LayoutData(w, r)

	var errStr string
	data, serveErr := serveAdminPage(c.Database, data)
	if serveErr != "" {
		errStr += "<br>" + serveErr
	}

	idStr, ok := mux.Vars(r)["category"]
	if !ok {
		data = data.MergeKV("error", "missing category ID")
		mustRender(w, r, "admin", data)
		return
	}

	id, err := strconv.ParseUint(idStr, 32, 10)
	if err != nil {
		data = data.MergeKV("error", "category ID is a not valid number")
		mustRender(w, r, "admin", data)
		return
	}

	var category models.Category
	category.ID = uint(id)

	err = c.Database.Delete(&category).Error
	if err != nil {
		log.WithError(err).Error("internal error when deleting a category")
		data = data.MergeKV("error", "internal error")
		data, _ = serveAdminPage(c.Database, data)
		mustRender(w, r, "admin", data)
		return
	}

	http.Redirect(w, r, "/admin", http.StatusFound)
}
