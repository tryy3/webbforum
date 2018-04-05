package handlers

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/apex/log"
	"github.com/jinzhu/gorm"
	"github.com/volatiletech/authboss"
	"github.com/tryy3/webbforum/models"
)

type idModel struct {
	ID uint
}

type updateAtModel struct {
	idModel
	UpdateAt time.Time
}

func serveHomePage(db *gorm.DB, data authboss.HTMLData) (authboss.HTMLData, string) {
	var categories []models.Category
	err := db.Find(&categories).Error
	if err != nil {
		log.WithError(err).Error("internal error when retrieving categories")
		return data, "internal error"
	}

	for i := range categories {
		name := strings.Replace(categories[i].Name, " ", "-", -1)
		path := fmt.Sprintf("%d-%s", categories[i].ID, name)
		urlPath, err := urlEncoded(path)
		if err != nil {
			log.WithField("path", path).Error("error parsing url")
			return data, "internal error"
		}
		categories[i].DisplayName = urlPath

		var threads []idModel
		db.Table("threads").Where("category_id = ?", categories[i].ID).Find(&threads)

		var postIds = make([]uint, len(threads))
		for i, thread := range threads {
			postIds[i] = thread.ID
		}

		var posts []updateAtModel
		db.Table("posts").Where("thread_id IN (?)", postIds).Order("updated_at desc").Find(&posts)

		categories[i].CountThreads = int64(len(threads))
		categories[i].CountPost = int64(len(posts))

		if len(posts) > 0 {
			var latest = models.Post{ID: posts[0].ID}
			db.Find(&latest)
			categories[i].LatestUpdate = &latest

			name := strings.Replace(latest.Thread.Name, " ", "-", -1)
			path := fmt.Sprintf("%d-%s", latest.Thread.ID, name)
			urlPath, err := urlEncoded(path)
			if err != nil {
				log.WithField("path", path).Error("error parsing url")
				return data, "internal error"
			}
			latest.Thread.DisplayName = urlPath
		}
	}

	data = data.MergeKV("categories", categories)
	return data, ""
}

// HomeHandler is the http handler for /
type HomeHandler struct {
	Database *gorm.DB
}

func (h HomeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	data := LayoutData(w, r)
	data, err := serveHomePage(h.Database, data)
	if err != "" {
		data = data.MergeKV("error", err)
	}
	mustRender(w, r, "index", data)
}
