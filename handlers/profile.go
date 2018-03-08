package handlers

import (
	"net/http"

	"github.com/tryy3/webbforum/models"
	"github.com/volatiletech/authboss"
)

type ProfileHandler struct {
	ab *authboss.Authboss
}

func NewProfileHandler(ab *authboss.Authboss) http.Handler {
	return ProfileHandler{ab}
}

func (p ProfileHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	u, ok := getUser(w, r, p.ab)
	if !ok {
		return
	}
	data := LayoutData(w, r).MergeKV("user", u)
	mustRender(w, r, "profil", data)
}

type ProfileEditHandler struct {
	storer models.UserStorer
	ab *authboss.Authboss
}

func NewProfileEditHandler(storer models.UserStorer, ab *authboss.Authboss) ProfileEditHandler {
	return ProfileEditHandler{storer, ab}
}

func (p ProfileEditHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	u, ok := getUser(w, r, p.ab)
	if !ok {
		return
	}

	ab, err := authboss.AttributesFromRequest(r)
	if badRequest(w, err) {
		return
	}

	err = p.storer.Put(u.(*models.User).Username, ab)
	if badRequest(w, err) {
		return
	}

	http.Redirect(w, r, r.RequestURI, http.StatusFound)
}