package auth

import (
	"fmt"
	"net/http"

	"github.com/apex/log"
	"github.com/gorilla/sessions"
	"github.com/volatiletech/authboss"
)

// SessionManager is the manager for managing the SessionStorer
type SessionManager struct {
	// The cookie name for sessions
	Name string

	// CookieStore interface from gorilla kit
	Store *sessions.CookieStore
}

// NewSessionManager creates a new session manager for managing the SessionStorer
func NewSessionManager(secure *sessions.CookieStore, Name string) *SessionManager {
	return &SessionManager{
		Name:  Name,
		Store: secure,
	}
}

// SessionStorer is the expected ClientStorer for authboss
type SessionStorer struct {
	w       http.ResponseWriter
	r       *http.Request
	manager *SessionManager
}

// NewSessionStorer creates new a new SessionStorer
func (m *SessionManager) NewSessionStorer(w http.ResponseWriter, r *http.Request) authboss.ClientStorer {
	return &SessionStorer{w, r, m}
}

// Get retrieves the session value from the Session
func (s SessionStorer) Get(key string) (string, bool) {
	session, err := s.manager.Store.Get(s.r, s.manager.Name)
	if err != nil {
		log.Error(err.Error())
		return "", false
	}

	strInf, ok := session.Values[key]
	if !ok {
		return "", false
	}

	str, ok := strInf.(string)
	if !ok {
		return "", false
	}

	return str, true
}

// Put modifies the session value of the session
func (s SessionStorer) Put(key, value string) {
	session, err := s.manager.Store.Get(s.r, s.manager.Name)
	if err != nil {
		log.Error(err.Error())
		return
	}

	session.Values[key] = value
	session.Save(s.r, s.w)
}

// Del deletes a session
func (s SessionStorer) Del(key string) {
	session, err := s.manager.Store.Get(s.r, s.manager.Name)
	if err != nil {
		fmt.Println(err)
		return
	}

	delete(session.Values, key)
	session.Save(s.r, s.w)
}
