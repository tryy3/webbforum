package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"github.com/apex/log"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/tryy3/webbforum/models"
	"github.com/volatiletech/authboss"
)

// servePostPage takes care of retriving the general data that post page need
func servePostPage(r *http.Request, data authboss.HTMLData, db *gorm.DB) (authboss.HTMLData, string) {
	cat, ok := mux.Vars(r)["thread"]
	if !ok {
		return data, "Invalid thread name"
	}

	var thread models.Thread
	threadS := strings.Split(cat, "-")
	if len(threadS) <= 0 {
		return data, "Invalid thread name"
	}

	id, err := strconv.ParseUint(threadS[0], 10, 32)
	if err != nil {
		return data, "Invalid thread ID"
	}
	thread.ID = uint(id)

	err = db.First(&thread).Error
	if err != nil {
		log.WithError(err).Error("internal error when retriving thread from database")
		return data, "Internal error"
	}

	threadName := strings.Replace(thread.Name, " ", "-", -1)
	path := fmt.Sprintf("%d-%s", thread.ID, threadName)
	urlPath, err := urlEncoded(path)
	if err != nil {
		log.WithField("path", path).Error("error parsing url")
		return data, "Internal error"
	}
	thread.DisplayName = urlPath
	data = data.MergeKV("thread", &thread)

	var posts []models.Post
	err = db.Where("thread_id = ?", thread.ID).Find(&posts).Error
	if err != nil {
		log.WithError(err).Error("internal error when retriving posts from database")
		return data, "Internal error"
	}

	conv := NewBBCodeConverter()
	for i := range posts {
		comment := conv.Convert(posts[i].Comment)
		posts[i].DisplayComment = template.HTML(comment)

		if posts[i].User != nil {
			posts[i].User = setCustomUserData(posts[i].User).(*models.User)
		}
	}

	data = data.MergeKV("posts", posts)
	return data, ""
}

// PostsShowHandler handler for showing a post
type PostsShowHandler struct {
	Database *gorm.DB
}

// ServeHTTP retrieves data regarding a post and shows it
func (s PostsShowHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	data := LayoutData(w, r)

	data, err := servePostPage(r, data, s.Database)
	if err != "" {
		data = data.MergeKV("error", err)
	}

	mustRender(w, r, "posts", data)
}

// PostCreateHandler handler for creating a post
type PostCreateHandler struct {
	Database *gorm.DB
	Authboss *authboss.Authboss
}

// ServeHTTP creates a new post in the db
func (p PostCreateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	data := LayoutData(w, r)

	data, errStr := servePostPage(r, data, p.Database)
	if errStr != "" {
		data = data.MergeKV("error", errStr)
		mustRender(w, r, "posts", data)
		return
	}

	u, ok := getUser(w, r, p.Authboss)
	if !ok {
		http.Redirect(w, r, strings.TrimRight(r.RequestURI, "/new_post"), http.StatusFound)
		return
	}
	user, ok := u.(*models.User)
	if !ok {
		http.Redirect(w, r, strings.TrimRight(r.RequestURI, "/new_post"), http.StatusFound)
		return
	}

	attr, err := authboss.AttributesFromRequest(r)
	if err != nil {
		log.WithError(err).Error("internal error when parsing request")
		data = data.MergeKV("error", "internal error")
		mustRender(w, r, "posts", data)
		return
	}

	comment, ok := attr.String("comment")
	if !ok {
		data = data.MergeKV("errs", map[string][]string{"comment": {"Invalid comment"}})
		mustRender(w, r, "posts", data)
		return
	}
	data = data.MergeKV("comment", comment)

	var post models.Post
	post.Comment = comment
	post.User = user
	post.Thread = data["thread"].(*models.Thread)

	err = p.Database.Create(&post).Error
	if err != nil {
		log.WithError(err).Error("internal error when creating a new post")
		data = data.MergeKV("error", "internal error")
		mustRender(w, r, "posts", data)
		return
	}

	http.Redirect(w, r, strings.TrimRight(r.RequestURI, "/new_post"), http.StatusFound)
}

// PostEditHandler handler for editing a post
type PostEditHandler struct {
	Database *gorm.DB
	Authboss *authboss.Authboss
}

// ServeHTTP modifies a posts
func (p PostEditHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	data := LayoutData(w, r)

	data, errStr := servePostPage(r, data, p.Database)
	if errStr != "" {
		data = data.MergeKV("error", errStr)
		mustRender(w, r, "posts", data)
		return
	}

	/*u, ok := getUser(w, r, p.Authboss)
	if !ok {
		http.Redirect(w, r, strings.TrimRight(r.RequestURI, "/new_post"), http.StatusFound)
		return
	}
	user, ok := u.(*models.User)
	if !ok {
		http.Redirect(w, r, strings.TrimRight(r.RequestURI, "/new_post"), http.StatusFound)
		return
	}*/

	attr, err := authboss.AttributesFromRequest(r)
	if err != nil {
		log.WithError(err).Error("internal error when parsing request")
		data = data.MergeKV("error", "internal error")
		mustRender(w, r, "posts", data)
		return
	}

	postVar, ok := mux.Vars(r)["post"]
	if !ok {
		data = data.MergeKV("error", "Invalid comment ID")
		mustRender(w, r, "posts", data)
		return
	}

	post := &models.Post{}
	id, err := strconv.ParseUint(postVar, 10, 32)
	if err != nil {
		data = data.MergeKV("error", "Invalid comment ID")
		mustRender(w, r, "posts", data)
		return
	}
	post.ID = uint(id)

	err = p.Database.First(&post).Error
	if err != nil {
		log.WithError(err).Error("internal error when retriving post from database")
		data = data.MergeKV("error", "internal error")
		mustRender(w, r, "posts", data)
		return
	}
	data = data.MergeKV("edit_post", post)

	if r.Method == "GET" {
		mustRender(w, r, "posts", data)
		return
	}

	comment, ok := attr.String("comment")
	if !ok {
		data = data.MergeKV("error", "Invalid comment")
		mustRender(w, r, "posts", data)
		return
	}
	post.Comment = comment

	err = p.Database.Model(&post).Updates(&post).Error
	if err != nil {
		log.WithError(err).Error("internal error when updating post in the database")
		data = data.MergeKV("error", "internal error")
		mustRender(w, r, "posts", data)
		return
	}

	http.Redirect(w, r, "/forums/thread/"+mux.Vars(r)["thread"], http.StatusFound)
}
