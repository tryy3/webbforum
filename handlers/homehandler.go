package handlers

import (
	"net/http"
)

// HomeHandler is the http handler for /
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	data := LayoutData(w, r)
	mustRender(w, r, "index", data)
}
