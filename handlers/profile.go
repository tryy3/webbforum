package handlers

import "net/http"

type ProfileHandler struct{}

func NewProfileHandler() http.Handler {
	return ProfileHandler{}
}

func (a ProfileHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	data := LayoutData(w, r)
	mustRender(w, r, "admin", data)
}
