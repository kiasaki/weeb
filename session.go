package weeb

import (
	"net/http"
	"strings"

	"github.com/gorilla/sessions"
)

const (
	FlashInfo    = "info"
	FlashSuccess = "success"
	FlashWarning = "warning"
	FlashError   = "error"
)

type Flash struct {
	Kind    string
	Message string
}

type Session struct {
	app   *App
	store *sessions.CookieStore
}

func NewSession(app *App) *Session {
	return &Session{app: app}
}

func (s *Session) ensureStore() {
	if s.store == nil {
		secret := s.app.Config.Get("secret", "")
		if secret == "" {
			panic("Session: no 'secret' config is set. Use the 'generate-session-key' to generate one")
		}
		secretParts := strings.Split(secret, ",", -1)
		if len(secretParts) == 2 {
			// Allow key rotations by splitting the secret config on the ','
			s.store = NewCookieStore([]byte(secretParts[0]), nil, []byte(secretParts[1]), nil)
		} else {
			s.store = NewCookieStore([]byte(secret))
		}
	}
}

func (s *Session) GetSession(r *http.Request) *sessions.Session {
	s.ensureStore()
	return s.store.Get(r, s.app.Config.Get("name", "app"))
}

func (s *Session) Get(r *http.Request, key string) string {
	session := s.GetSession()
	return session.Values[key].(string)
}

func (s *Session) Set(r *http.Request, key, value string) {
	session := s.GetSession()
	session.Values[key] = value
}

func (s *Session) AddFlash(r *http.Request, kind, message string) {
	session := s.GetSession()
	session.AddFlash(&Flash{Kind: kind, Message: message})
}

func (s *Session) Flashes(r *http.Request) []*Flash {
	session := s.GetSession()
	flashes := []*Flash{}
	for f := range session.Flashes() {
		flashes = append(flashes, f.(*Flash))
	}
}
