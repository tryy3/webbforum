package handlers

import "net/http"

func AdminHandler(w http.ResponseWriter, r *http.Request) {
	data := LayoutData(w, r)
	mustRender(w, r, "admin", data)
}
