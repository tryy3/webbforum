package middleware

import (
	"net/http"

	"github.com/apex/log"
	"github.com/justinas/nosurf"
)

// NoSurfingMiddleware is a middleware for protecting against CSRF attacks
func NoSurfingMiddleware(h http.Handler) http.Handler {
	surfing := nosurf.New(h)
	surfing.SetFailureHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Errorf("failed to validate XSRF Token: %s", nosurf.Reason(r))
		w.WriteHeader(http.StatusBadRequest)
	}))
	return surfing
}
