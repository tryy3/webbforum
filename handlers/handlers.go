package handlers

import (
	"net/http"
	"strconv"

	"github.com/apex/log"
	"github.com/volatiletech/authboss"
)

func badRequest(w http.ResponseWriter, err error) bool {
	if err == nil {
		return false
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusBadRequest)
	log.WithError(err).Error("bad request")

	return true
}

func getUser(w http.ResponseWriter, r *http.Request, ab *authboss.Authboss) (interface{}, bool) {
	u, err := ab.CurrentUser(w, r)
	if err != nil {
		log.WithError(err).Error("error fetching current user")
		w.WriteHeader(http.StatusInternalServerError)
		return nil, false
	} else if u == nil {
		log.Errorf("redirecting unauthorized user from: %s", r.URL.Path)
		http.Redirect(w, r, "/", http.StatusFound)
		return nil, false
	}
	return u, true
}

func getCategoryID(attr authboss.Attributes) (uint, string) {
	idStr, ok := attr.String("category_id")
	if !ok {
		return 0, "missing category ID"
	}

	id, err := strconv.ParseUint(idStr, 32, 10)
	if err != nil {
		return 0, "category ID is a not valid number"
	}

	return uint(id), ""
}

func getThreadID(attr authboss.Attributes) (uint, string) {
	idStr, ok := attr.String("thread_id")
	if !ok {
		return 0, "missing thread ID"
	}

	id, err := strconv.ParseUint(idStr, 32, 10)
	if err != nil {
		return 0, "thread ID is a not valid number"
	}

	return uint(id), ""
}
