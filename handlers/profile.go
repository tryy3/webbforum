package handlers

import (
	"net/http"

	"github.com/davecgh/go-spew/spew"
	"github.com/gorilla/mux"
	"github.com/tryy3/webbforum/models"
)

type ProfileHandler struct {
	store models.UserStorer
}

func NewProfileHandler(store models.UserStorer) http.Handler {
	return ProfileHandler{store}
}

func (a ProfileHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	data := LayoutData(w, r)
	username, ok := mux.Vars(r)["username"]
	if !ok {
		mustRender(w, r, "profile", data)
		return
	}

	user, err := a.store.Get(username)
	if err == nil {
		data.MergeKV("user", user)
	}

	spew.Dump(user)

	mustRender(w, r, "profile", data)
}
