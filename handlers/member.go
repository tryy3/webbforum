package handlers

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/tryy3/webbforum/models"
)

type MemberHandler struct {
	store models.UserStorer
}

func NewMemberHandler(store models.UserStorer) http.Handler {
	return MemberHandler{store}
}

func (a MemberHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	data := LayoutData(w, r)
	username, ok := mux.Vars(r)["username"]
	if !ok {
		mustRender(w, r, "medlem", data)
		return
	}

	user, err := a.store.Get(username)
	if err == nil {
		data.MergeKV("user", user)
	}

	mustRender(w, r, "medlem", data)
}
