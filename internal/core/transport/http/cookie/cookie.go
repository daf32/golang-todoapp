package core_http_cookie

import (
	"net/http"
	"time"
)

type Manager struct {
	secure bool
}

func NewManager(secure bool) *Manager {
	return &Manager{secure: secure}
}

func (m *Manager) Set(
	rw http.ResponseWriter,
	name, value string,
	maxAge time.Duration,
) {
	http.SetCookie(rw, &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		MaxAge:   int(maxAge.Seconds()),
		Expires:  time.Now().Add(maxAge),
		HttpOnly: true,
		Secure:   m.secure,
		SameSite: http.SameSiteLaxMode,
	})
}

func (m *Manager) Clear(rw http.ResponseWriter, name string) {
	http.SetCookie(rw, &http.Cookie{
		Name:     name,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   m.secure,
		SameSite: http.SameSiteLaxMode,
	})
}
