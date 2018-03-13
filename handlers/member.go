package handlers

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/tryy3/webbforum/models"
)

// MemberHandler handles serving of the /anvandare page
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
		mustRender(w, r, "user", data)
		return
	}

	u, err := a.store.Get(username)
	if err == nil {
		u = setImageProfile(u)
		data.MergeKV("user", u)
	}

	mustRender(w, r, "user", data)
}

// setImageProfile will attempt to set the ProfileImageURL if ProfileImage exists on the user
func setImageProfile(u interface{}) interface{} {
	user, ok := u.(*models.User)
	if !ok {
		return u
	}

	if user.ProfileImage == nil {
		return u
	}

	profileURL := "/images/" + user.ProfileImage.Base64Hash
	user.ProfileImageURL = profileURL
	return user
}
