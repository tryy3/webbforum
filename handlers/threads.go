package handlers

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/apex/log"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/volatiletech/authboss"
	"github.com/tryy3/webbforum/models"
)

func serveThreadsPage(db *gorm.DB, data authboss.HTMLData) (authboss.HTMLData, string) {
	var threads []models.Thread
	err := db.Find(&threads).Error
	if err != nil {
		log.WithError(err).Error("internal error when retrieving threads")
		return data, "internal error"
	}

	for i := range threads {
		var post models.Post
		var count int64
		stmt := db.Model(&post).Where("thread_id = ?", threads[i].ID).Order("updated_at desc").First(&post).Count(&count)
		err = stmt.Error
		if err != nil {
			threads[i].LatestPost = time.Time{}
		} else {
			threads[i].LatestPost = post.UpdatedAt
		}

		threads[i].CountPost = count
	}

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

	if r.Method == "GET" {
		mustRender(w, r, "create_threads", data)
		return
	}

	hasPerm, err := models.HasPermission(t.Database, user, models.Permissions.CREATETHREAD)
	if err != nil {
		log.WithError(err).Error("error when trying to retrieve if user has permission")
		data.MergeKV("error", "internal error")
		mustRender(w, r, "create_threads", data)
		return
	}
	if !hasPerm {
		data.MergeKV("error", "You do not have permission to do this.")
		mustRender(w, r, "create_threads", data)
		return
	}

	cat, ok := mux.Vars(r)["category"]
	if !ok {
		data = data.MergeKV("error", "Invalid category name")
		mustRender(w, r, "create_threads", data)
		return
	}

	var category models.Category
	catS := strings.Split(cat, "-")
	if len(catS) <= 0 {
		data = data.MergeKV("error", "Invalid category name")
		mustRender(w, r, "create_threads", data)
		return
	}

	id, err := strconv.ParseUint(catS[0], 32, 10)
	if err != nil {
		data = data.MergeKV("error", "Invalid category ID")
		mustRender(w, r, "create_threads", data)
		return
	}
	category.ID = uint(id)

	err = t.Database.First(&category).Error
	if err != nil {
		log.WithError(err).Error("internal error when retriving category from database")
		data = data.MergeKV("internal error")
		mustRender(w, r, "create_threads", data)
		return
	}

	attr, err := authboss.AttributesFromRequest(r)
	if err != nil {
		log.WithError(err).Error("internal error when parsing request")
		data = data.MergeKV("error", "internal error")
		data, _ = serveThreadsPage(t.Database, data)
		mustRender(w, r, "create_threads", data)
		return
	}

	var errorList = map[string][]string{}
	name, ok := attr.String("thread_name")
	if !ok {
		errorList["thread_name"] = []string{"Invalid thread name"}
	} else {
		data = data.MergeKV("thread_name", name)
	}

	comment, ok := attr.String("thread_message")
	if !ok {
		errorList["thread_message"] = []string{"Invalid thread message"}
	} else {
		data = data.MergeKV("thread_message", comment)
	}

	if len(errorList) > 0 {
		data = data.MergeKV("errs", errorList)
		mustRender(w, r, "create_threads", data)
		return
	}

	var thread models.Thread
	thread.CreatedBy = user
	thread.Category = &category
	thread.Name = name

	var post models.Post
	post.User = user
	post.Thread = &thread
	post.Comment = comment

	tx := t.Database.Begin()

	if tx.Error != nil {
		data = data.MergeKV("error", "internal error")
		log.WithError(tx.Error).Error("error when begining a database transaction")
		mustRender(w, r, "create_threads", data)
		return
	}

	if err := tx.Create(&thread).Error; err != nil {
		tx.Rollback()
		data = data.MergeKV("error", "internal error")
		log.WithError(tx.Error).Error("error creating a new thread")
		mustRender(w, r, "create_threads", data)
		return
	}

	if err := tx.Create(&post).Error; err != nil {
		tx.Rollback()
		data = data.MergeKV("error", "internal error")
		log.WithError(tx.Error).Error("error creating a new post")
		mustRender(w, r, "create_threads", data)
		return
	}

	if err := tx.Commit().Error; err != nil {
		data = data.MergeKV("error", "internal error")
		log.WithError(tx.Error).Error("error commiting transaction")
		mustRender(w, r, "create_threads", data)
		return
	}

	threadName := strings.Replace(thread.Name, " ", "-", -1)
	path := fmt.Sprintf("/forums/thread/%d-%s", thread.ID, threadName)
	urlPath, err := urlEncoded(path)
	if err != nil {
		log.WithField("path", path).Error("error parsing url")
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	http.Redirect(w, r, urlPath, http.StatusFound)
}

type ThreadDeleteHandler struct {
	Database *gorm.DB
	Authboss *authboss.Authboss
}

func (t ThreadDeleteHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	data := LayoutData(w, r)

	u, ok := getUser(w, r, t.Authboss)
	if !ok {
		return
	}
	user, ok := u.(*models.User)
	if !ok {
		return
	}

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

	data = data.MergeKV("category_id", thread.Category.ID)

	if thread.CreatedByID == user.ID {
		hasPerm, err := models.HasPermission(t.Database, user, models.Permissions.DELETESELFTHREAD)
		if err != nil {
			log.WithError(err).Error("error when trying to retrieve if user has permission")
			data.MergeKV("error", "internal error")
			mustRender(w, r, "create_threads", data)
			return
		}
		if !hasPerm {
			hasPerm, err := models.HasPermission(t.Database, user, models.Permissions.DELETETHREAD)
			if err != nil {
				log.WithError(err).Error("error when trying to retrieve if user has permission")
				data.MergeKV("error", "internal error")
				mustRender(w, r, "create_threads", data)
				return
			}
			if !hasPerm {
				data.MergeKV("error", "You do not have permission to do this.")
				mustRender(w, r, "create_threads", data)
				return
			}
		}
	} else {
		hasPerm, err := models.HasPermission(t.Database, user, models.Permissions.DELETETHREAD)
		if err != nil {
			log.WithError(err).Error("error when trying to retrieve if user has permission")
			data.MergeKV("error", "internal error")
			mustRender(w, r, "create_threads", data)
			return
		}
		if !hasPerm {
			data.MergeKV("error", "You do not have permission to do this.")
			mustRender(w, r, "create_threads", data)
			return
		}
	}

	err = t.Database.Delete(&thread).Error
	if err != nil {
		log.WithError(err).Error("Internal error when deleting a thread")
		data = data.MergeKV("error", "internal error")
		data, _ = serveThreadsPage(t.Database, data)
		mustRender(w, r, "threads", data)
		return
	}

	http.Redirect(w, r, "/forums/"+thread.Category.Name, http.StatusFound)
}

type ThreadEditHandler struct {
	Database *gorm.DB
	Authboss *authboss.Authboss
}

func (t ThreadEditHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	data := LayoutData(w, r)

	u, ok := getUser(w, r, t.Authboss)
	if !ok {
		return
	}
	user, ok := u.(*models.User)
	if !ok {
		return
	}

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
	data = data.MergeKV("category_id", thread.Category.ID)



	if thread.CreatedByID == user.ID {
		hasPerm, err := models.HasPermission(t.Database, user, models.Permissions.EDITSELFTHREAD)
		if err != nil {
			log.WithError(err).Error("error when trying to retrieve if user has permission")
			data.MergeKV("error", "internal error")
			mustRender(w, r, "create_threads", data)
			return
		}
		if !hasPerm {
			hasPerm, err := models.HasPermission(t.Database, user, models.Permissions.EDITTHREAD)
			if err != nil {
				log.WithError(err).Error("error when trying to retrieve if user has permission")
				data.MergeKV("error", "internal error")
				mustRender(w, r, "create_threads", data)
				return
			}
			if !hasPerm {
				data.MergeKV("error", "You do not have permission to do this.")
				mustRender(w, r, "create_threads", data)
				return
			}
		}
	} else {
		hasPerm, err := models.HasPermission(t.Database, user, models.Permissions.EDITTHREAD)
		if err != nil {
			log.WithError(err).Error("error when trying to retrieve if user has permission")
			data.MergeKV("error", "internal error")
			mustRender(w, r, "create_threads", data)
			return
		}
		if !hasPerm {
			data.MergeKV("error", "You do not have permission to do this.")
			mustRender(w, r, "create_threads", data)
			return
		}
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

	http.Redirect(w, r, "/forums/"+thread.Category.Name, http.StatusFound)
}

type ThreadShowHandler struct {
	Database *gorm.DB
}

func (t ThreadShowHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	data := LayoutData(w, r)

	cat, ok := mux.Vars(r)["category"]
	if !ok {
		data = data.MergeKV("error", "Invalid category name")
		mustRender(w, r, "threads", data)
		return
	}

	var category models.Category
	catS := strings.Split(cat, "-")
	if len(catS) <= 0 {
		data = data.MergeKV("error", "Invalid category name")
		mustRender(w, r, "threads", data)
		return
	}

	id, err := strconv.ParseUint(catS[0], 32, 10)
	if err != nil {
		data = data.MergeKV("error", "Invalid category ID")
		mustRender(w, r, "threads", data)
		return
	}
	category.ID = uint(id)

	err = t.Database.First(&category).Error
	if err != nil {
		data = data.MergeKV("error", "Internal error")
		mustRender(w, r, "threads", data)
		return
	}

	if len(catS) <= 1 || strings.Join(catS[1:], " ") != category.Name {
		name := strings.Replace(category.Name, " ", "-", -1)
		path := fmt.Sprintf("/forums/%d-%s", category.ID, name)
		http.Redirect(w, r, url.PathEscape(path), http.StatusPermanentRedirect)
		return
	}
	name := strings.Replace(category.Name, " ", "-", -1)
	path := fmt.Sprintf("%d-%s", category.ID, name)
	urlPath, err := urlEncoded(path)
	if err != nil {
		data = data.MergeKV("error", "internal error")
		log.WithField("path", path).Error("error parsing url")
		mustRender(w, r, "threads", data)
		return
	}
	category.DisplayName = urlPath
	data = data.MergeKV("category", category)

	var threads []models.Thread
	err = t.Database.Where("category_id = ?", category.ID).Find(&threads).Error
	if err != nil {
		data = data.MergeKV("error", "Internal error")
		mustRender(w, r, "threads", data)
		return
	}

	for i := range threads {
		name := strings.Replace(threads[i].Name, " ", "-", -1)
		path := fmt.Sprintf("%d-%s", threads[i].ID, name)
		urlPath, err := urlEncoded(path)
		if err != nil {
			data = data.MergeKV("error", "internal error")
			log.WithField("path", path).Error("error parsing url")
			mustRender(w, r, "threads", data)
			return
		}
		threads[i].DisplayName = urlPath

		var post models.Post
		var count int64
		stmt := t.Database.Table("posts").Where("thread_id = ?", threads[i].ID).Order("updated_at desc")
		err = stmt.First(&post).Error
		if err != nil {
			threads[i].LatestPost = time.Time{}
		} else {
			threads[i].LatestPost = post.UpdatedAt
		}

		stmt.Count(&count)
		threads[i].CountPost = count
	}

	data = data.MergeKV("threads", threads)
	mustRender(w, r, "threads", data)
}
