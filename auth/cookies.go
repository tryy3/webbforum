package auth

import (
	"net/http"
	"time"

	"github.com/apex/log"
	"github.com/gorilla/securecookie"
	"github.com/volatiletech/authboss"
)

// CookieManager is the manager for managing the CookieStorer
type CookieManager struct {
	// Expiry is used for setting when a cookie should be expired (current time + Expiry)
	Expiry time.Duration

	// Store is the SecureCookie interface for accessing cookies
	Store *securecookie.SecureCookie
}

// NewCookieManager creates a new cookie manager for managing the CookieStorer
func NewCookieManager(secure *securecookie.SecureCookie, Expiry time.Duration) *CookieManager {
	return &CookieManager{
		Expiry: Expiry,
		Store:  secure,
	}
}

// CookieStorer is the expected ClientStorer for authboss
type CookieStorer struct {
	w       http.ResponseWriter
	r       *http.Request
	manager *CookieManager
}

// NewCookieStorer creates new a new CookieStorer
func (m *CookieManager) NewCookieStorer(w http.ResponseWriter, r *http.Request) authboss.ClientStorer {
	return &CookieStorer{w, r, m}
}

// Get retrieves the cookie value
func (s CookieStorer) Get(key string) (string, bool) {
	cookie, err := s.r.Cookie(key)
	if err != nil {
		if err.Error() != "http: named cookie not present" {
			log.Error(err.Error())
		}
		return "", false
	}

	var value string
	err = s.manager.Store.Decode(key, cookie.Value, &value)
	if err != nil {
		return "", false
	}
	return value, true
}

// Put modifies/creates a cookie
func (s CookieStorer) Put(key, value string) {
	encoded, err := s.manager.Store.Encode(key, value)
	if err != nil {
		log.Error(err.Error())
		return
	}

	cookie := &http.Cookie{
		Expires: time.Now().UTC().Add(s.manager.Expiry),
		Name:    key,
		Value:   encoded,
		Path:    "/",
	}
	http.SetCookie(s.w, cookie)
}

// Del will set the MaxAge of the cookie to -1 (it will remove the cookie)
func (s CookieStorer) Del(key string) {
	cookie := &http.Cookie{
		MaxAge: -1,
		Name:   key,
		Path:   "/",
	}
	http.SetCookie(s.w, cookie)
}
