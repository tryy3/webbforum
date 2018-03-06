package middleware

import (
	"net/http"

	"github.com/apex/log"
	"github.com/davecgh/go-spew/spew"

	"github.com/volatiletech/authboss"
)

type PermissionProtector struct {
	perms []string
	f     http.Handler
	ab    *authboss.Authboss
}

func PermissionProtect(f http.Handler, perms []string, ab *authboss.Authboss) http.Handler {
	return LoggedInProtector{PermissionProtector{perms, f, ab}, ab}
}

func (p PermissionProtector) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if u, err := p.ab.CurrentUser(w, r); err != nil {
		log.WithError(err).Error("error fetching current user")
		w.WriteHeader(http.StatusInternalServerError)
	} else if u == nil {
		log.Errorf("redirecting unauthorized user from: %s", r.URL.Path)
		http.Redirect(w, r, "/", http.StatusFound)
	} else {
		spew.Dump(u)
		p.ServeHTTP(w, r)
	}
}

type LoggedInProtector struct {
	f  http.Handler
	ab *authboss.Authboss
}

func LoggedInProtect(f http.Handler, ab *authboss.Authboss) http.Handler {
	return LoggedInProtector{f, ab}
}

func (p LoggedInProtector) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if u, err := p.ab.CurrentUser(w, r); err != nil {
		log.WithError(err).Error("error fetching current user")
		w.WriteHeader(http.StatusInternalServerError)
	} else if u == nil {
		log.Errorf("redirecting unauthorized user from: %s", r.URL.Path)
		http.Redirect(w, r, "/", http.StatusFound)
	} else {
		p.ServeHTTP(w, r)
	}
}
