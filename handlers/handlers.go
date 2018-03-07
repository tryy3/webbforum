package handlers

import (
	"net/http"
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