package handlers

import (
	"net/http"

	"github.com/apex/log"
	"github.com/jinzhu/gorm"
	"github.com/tryy3/webbforum/models"
	"github.com/volatiletech/authboss"
)

func serveThreadsPage(db *gorm.DB, data authboss.HTMLData) (authboss.HTMLData, string) {
	var categories []models.Category
	err := db.Find(&categories).Error
	if err != nil {
		log.WithError(err).Error("internal error when retrieving categories")
		return data, "internal error"
	}

	var threads []models.Thread
	err = db.Find(&threads).Error
	if err != nil {
		log.WithError(err).Error("internal error when retrieving threads")
		return data, "internal error"
	}

	data = data.MergeKV("categories", categories)
	data = data.MergeKV("threads", threads)
	return data, ""
}

type ThreadCreateHandler struct {
	Database *gorm.DB
	Authboss *authboss.Authboss
}

func (t ThreadCreateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	u, ok := getUser(w, r, t.Authboss)
	if !ok {
		return
	}
	user, ok := u.(*models.User)
	if !ok {
		return
	}

	data := LayoutData(w, r)

	attr, err := authboss.AttributesFromRequest(r)
	if err != nil {
		log.WithError(err).Error("internal error when parsing request")
		data = data.MergeKV("error", "internal error")
		data, _ = serveThreadsPage(t.Database, data)
		mustRender(w, r, "threads", data)
		return
	}

	var thread models.Thread
	thread.CreatedBy = user

	var category models.Category
	id, errStr := getCategoryID(attr)
	if errStr != "" {
		data, serveErr := serveThreadsPage(t.Database, data)
		if serveErr != "" {
			errStr += "<br>" + serveErr
		}
		data = data.MergeKV("error", errStr)
		mustRender(w, r, "threads", data)
		return
	}
	category.ID = id

	result := t.Database.First(&category)
	if result.Error != nil || result.RecordNotFound() {
		data = data.MergeKV("thread_errs", map[string][]string{"thread_category": {"Invalid category ID"}})
		data, errStr := serveThreadsPage(t.Database, data)
		if errStr != "" {
			data = data.MergeKV("error", errStr)
		}
		mustRender(w, r, "threads", data)
		return
	}
	thread.Category = &category

	name, ok := attr.String("thread_name")
	if !ok {
		data = data.MergeKV("thread_errs", map[string][]string{"thread_name": {"Invalid thread name"}})
		data, errStr := serveThreadsPage(t.Database, data)
		if errStr != "" {
			data = data.MergeKV("error", errStr)
		}
		mustRender(w, r, "threads", data)
		return
	}
	thread.Name = name

	if err := t.Database.Create(&thread).Error; err != nil {
		log.WithError(err).Error("internal error when parsing request")
		data = data.MergeKV("error", "internal error")
		data, _ = serveThreadsPage(t.Database, data)
		mustRender(w, r, "threads", data)
		return
	}

	http.Redirect(w, r, "/forums/" + category.Name, http.StatusFound)
}

type ThreadDeleteHandler struct {
	Database *gorm.DB
}

func (t ThreadDeleteHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	data := LayoutData(w, r)

	attr, err := authboss.AttributesFromRequest(r)
	if err != nil {
		log.WithError(err).Error("internal error when parsing request")
		data = data.MergeKV("error", "internal error")
		data, _ = serveThreadsPage(t.Database, data)
		mustRender(w, r, "threads", data)
		return
	}

	var thread models.Thread
	id, errStr := getThreadID(attr)
	if errStr != "" {
		data, serveErr := serveThreadsPage(t.Database, data)
		if serveErr != "" {
			errStr += "<br>" + serveErr
		}
		data = data.MergeKV("error", errStr)
		mustRender(w, r, "threads", data)
		return
	}
	thread.ID = id

	err = t.Database.First(&thread).Error
	if err != nil {
		log.WithError(err).Error("Internal error when deleting a thread")
		data = data.MergeKV("error", "internal error")
		data, _ = serveThreadsPage(t.Database, data)
		mustRender(w, r, "threads", data)
		return
	}

	err = t.Database.Delete(&thread).Error
	if err != nil {
		log.WithError(err).Error("Internal error when deleting a thread")
		data = data.MergeKV("error", "internal error")
		data, _ = serveThreadsPage(t.Database, data)
		mustRender(w, r, "threads", data)
		return
	}

	http.Redirect(w, r, "/forums/" + thread.Category.Name, http.StatusFound)
}

type ThreadEditHandler struct {
	Database *gorm.DB
}

func (t ThreadEditHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	data := LayoutData(w, r)

	attr, err := authboss.AttributesFromRequest(r)
	if err != nil {
		log.WithError(err).Error("internal error when parsing request")
		data = data.MergeKV("error", "internal error")
		data, _ = serveThreadsPage(t.Database, data)
		mustRender(w, r, "threads", data)
		return
	}

	var thread models.Thread
	id, errStr := getThreadID(attr)
	if errStr != "" {
		data, serveErr := serveThreadsPage(t.Database, data)
		if serveErr != "" {
			errStr += "<br>" + serveErr
		}
		data = data.MergeKV("error", errStr)
		mustRender(w, r, "threads", data)
		return
	}
	thread.ID = id

	err = t.Database.Find(&thread).Error
	if err != nil {
		log.WithError(err).Error("internal error when trying to find the existing thread")
		data = data.MergeKV("error", "internal error")
		data, _ = serveThreadsPage(t.Database, data)
		mustRender(w, r, "threads", data)
		return
	}

	name, ok := attr.String("thread_name")
	if !ok {
		data = data.MergeKV("thread_errs", map[string][]string{"thread_name": {"Invalid thread name"}})
		data, errStr := serveThreadsPage(t.Database, data)
		if errStr != "" {
			data = data.MergeKV("error", errStr)
		}
		mustRender(w, r, "threads", data)
		return
	}
	thread.Name = name

	catID, errStr := getCategoryID(attr)
	if errStr != "" {
		data, serveErr := serveThreadsPage(t.Database, data)
		if serveErr != "" {
			errStr += "<br>" + serveErr
		}
		data = data.MergeKV("error", errStr)
		mustRender(w, r, "threads", data)
		return
	}

	if thread.CategoryID != catID {
		// TODO Test this
		var category models.Category
		category.ID = catID
		result := t.Database.First(&category)
		if result.Error != nil || result.RecordNotFound() {
			data = data.MergeKV("thread_errs", map[string][]string{"thread_category": {"Invalid category ID"}})
			data, errStr := serveThreadsPage(t.Database, data)
			if errStr != "" {
				data = data.MergeKV("error", errStr)
			}
			mustRender(w, r, "threads", data)
			return
		}

		thread.Category = &category
	}

	err = t.Database.Model(&thread).Updates(&thread).Error
	if err != nil {
		log.WithError(err).Error("internal error when updating a thread")
		data = data.MergeKV("error", "internal error")
		data, _ = serveThreadsPage(t.Database, data)
		mustRender(w, r, "threads", data)
		return
	}

	http.Redirect(w, r, "/forums/" + thread.Category.Name, http.StatusFound)
}
