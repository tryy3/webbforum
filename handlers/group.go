package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/apex/log"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/tryy3/webbforum/models"
	"github.com/volatiletech/authboss"
)

// GroupCreateHandler creates a handler for creating group in database
type GroupCreateHandler struct {
	Database *gorm.DB
}

// ServeHTTP handle creating new group
func (c GroupCreateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	data := LayoutData(w, r)

	attr, err := authboss.AttributesFromRequest(r)
	if err != nil {
		log.WithError(err).Error("internal error when parsing request")
		data = data.MergeKV("error", "internal error")
		data, _ = serveAdminPage(c.Database, data)
		mustRender(w, r, "admin", data)
		return
	}

	var group models.Group
	if name, ok := attr.String("group_name"); ok {
		group.Name = name
	}
	if description, ok := attr.String("group_description"); ok {
		group.Description = description
	}
	data = data.MergeKV("group", group)

	if group.Name == "" {
		data = data.MergeKV("errs", map[string][]string{"group_name": {"Missing group name"}})
		data, errStr := serveAdminPage(c.Database, data)
		if errStr != "" {
			data = data.MergeKV("error", errStr)
		}
		mustRender(w, r, "admin", data)
		return
	}

	err = c.Database.Create(&group).Error
	if err != nil {
		log.WithError(err).Error("internal error when creating new group")
		data = data.MergeKV("error", "internal error")
		data, _ = serveAdminPage(c.Database, data)
		mustRender(w, r, "admin", data)
		return
	}

	http.Redirect(w, r, "/admin", http.StatusFound)
}

// GroupEditHandler handles editing group in database
type GroupEditHandler struct {
	Database *gorm.DB
}

// ServeHTTP handles editing group in database
func (c GroupEditHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	data := LayoutData(w, r)

	attr, err := authboss.AttributesFromRequest(r)
	if err != nil {
		log.WithError(err).Error("internal error when parsing request")
		data = data.MergeKV("error", "internal error")
		data, _ = serveAdminPage(c.Database, data)
		mustRender(w, r, "admin", data)
		return
	}

	var group models.Group
	id, errStr := getGroupID(attr)
	if errStr != "" {
		data, serveErr := serveAdminPage(c.Database, data)
		if serveErr != "" {
			errStr += "<br>" + serveErr
		}
		data = data.MergeKV("error", errStr)
		mustRender(w, r, "admin", data)
		return
	}
	group.ID = id

	result := c.Database.First(&group)
	if result.Error != nil || result.RecordNotFound() {
		errStr := "tried to modify invalid group"
		data, serveErr := serveAdminPage(c.Database, data)
		if serveErr != "" {
			errStr += "<br>" + serveErr
		}
		data = data.MergeKV("error", errStr)
		mustRender(w, r, "admin", data)
		return
	}

	if name, ok := attr.String("group_name"); ok {
		group.Name = name
	}
	if description, ok := attr.String("group_description"); ok {
		group.Description = description
	}
	data = data.MergeKV("group_edit", group)

	if group.Name == "" {
		data = data.MergeKV("group_errs", map[string][]string{"group_name": {"Missing group name"}})
		data, errStr := serveAdminPage(c.Database, data)
		if errStr != "" {
			data = data.MergeKV("error", errStr)
		}
		mustRender(w, r, "admin", data)
		return
	}

	err = c.Database.Model(&group).Updates(&group).Error
	if err != nil {
		log.WithError(err).Error("internal error when updating a group")
		data = data.MergeKV("error", "internal error")
		data, _ = serveAdminPage(c.Database, data)
		mustRender(w, r, "admin", data)
		return
	}

	var perms []models.Permission
	err = c.Database.Where("group_id = ?", group.ID).Find(&perms).Error
	if err != nil {
		log.WithError(err).Error("internal error when trying to get existing group permissions")
		data = data.MergeKV("error", "internal error")
		data, _ = serveAdminPage(c.Database, data)
		mustRender(w, r, "admin", data)
		return
	}

	var deleteIDs []uint

	for _, perm := range perms {
		if _, ok := attr["permission_"+perm.Permission]; !ok {
			deleteIDs = append(deleteIDs, perm.ID)
		}
	}

	if len(deleteIDs) > 0 {
		err = c.Database.Where("id IN (?)", deleteIDs).Delete(models.Permission{}).Error
		if err != nil {
			log.WithError(err).Error("internal error when trying to delete permission")
			data = data.MergeKV("error", "internal error")
			data, _ = serveAdminPage(c.Database, data)
			mustRender(w, r, "admin", data)
			return
		}
	}

	for perm := range attr {
		if !strings.Contains(perm, "permission_") {
			continue
		}
		parsed := strings.TrimPrefix(perm, "permission_")

		var found bool
		for _, p := range perms {
			if parsed == p.Permission {
				found = true
				break
			}
		}

		if !found {
			permission := models.Permission{
				GroupID:    group.ID,
				Permission: parsed,
			}

			err = c.Database.Create(&permission).Error
			if err != nil {
				log.WithError(err).Error("internal error when trying to insert a new permission")
				data = data.MergeKV("error", "internal error")
				data, _ = serveAdminPage(c.Database, data)
				mustRender(w, r, "admin", data)
				return
			}
			continue
		}
	}

	http.Redirect(w, r, "/admin", http.StatusFound)
}

// GroupDeleteHandler handles deleting a group in database
type GroupDeleteHandler struct {
	Database *gorm.DB
}

// ServeHTTP handles deleting a group in database
func (c GroupDeleteHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	data := LayoutData(w, r)

	var errStr string
	data, serveErr := serveAdminPage(c.Database, data)
	if serveErr != "" {
		errStr += "<br>" + serveErr
	}

	idStr, ok := mux.Vars(r)["group"]
	if !ok {
		data = data.MergeKV("error", "missing group ID")
		mustRender(w, r, "admin", data)
		return
	}

	id, err := strconv.ParseUint(idStr, 32, 10)
	if err != nil {
		data = data.MergeKV("error", "group ID is a not valid number")
		mustRender(w, r, "admin", data)
		return
	}

	var group models.Group
	group.ID = uint(id)

	err = c.Database.Where("group_id = ?", group.ID).Delete(models.Group{}).Error
	if err != nil {
		log.WithError(err).Error("internal error when deleting group permissions")
		data = data.MergeKV("error", "internal error")
		data, _ = serveAdminPage(c.Database, data)
		mustRender(w, r, "admin", data)
		return
	}

	err = c.Database.Delete(&group).Error
	if err != nil {
		log.WithError(err).Error("internal error when deleting a group")
		data = data.MergeKV("error", "internal error")
		data, _ = serveAdminPage(c.Database, data)
		mustRender(w, r, "admin", data)
		return
	}

	http.Redirect(w, r, "/admin", http.StatusFound)
}
