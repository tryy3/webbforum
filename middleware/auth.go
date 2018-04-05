package middleware

import (
	"net/http"

	"github.com/apex/log"
	"github.com/volatiletech/authboss"
)

// LoggedInProtector is a middleware handler for checking if a user is logged in or not
type LoggedInProtector struct {
	h  http.Handler
	ab *authboss.Authboss
}

func LoggedInProtect(h http.Handler, ab *authboss.Authboss) http.Handler {
	return LoggedInProtector{h, ab}
}

func (p LoggedInProtector) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	u, err := p.ab.CurrentUser(w, r)
	if err != nil {
		log.WithError(err).Error("error fetching current user")
		w.WriteHeader(http.StatusInternalServerError)
	} else if u == nil {
		log.Errorf("redirecting unauthorized user from: %s", r.URL.Path)
		http.Redirect(w, r, "/", http.StatusFound)
	} else {
		p.h.ServeHTTP(w, r)
	}
}
