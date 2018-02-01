package weeb

import (
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
	ctx   *Context
	store *sessions.CookieStore
}

func NewSession(ctx *Context) *Session {
	return &Session{ctx: ctx}
}

func (s *Session) ensureStore() {
	if s.store == nil {
		secret := s.ctx.App().Config.Get("secret", "")
		if secret == "" {
			panic("Session: no 'secret' config is set. Use the 'generate-session-key' to generate one")
		}
		secretParts := strings.Split(secret, ",")
		if len(secretParts) == 2 {
			// Allow key rotations by splitting the secret config on the ','
			s.store = sessions.NewCookieStore([]byte(secretParts[0]), nil, []byte(secretParts[1]), nil)
		} else {
			s.store = sessions.NewCookieStore([]byte(secret))
		}
	}
}

func (s *Session) save() {
	if s.store == nil {
		return
	}
	err := s.GetSession().Save(s.ctx.Request, s.ctx.Response)
	if err != nil {
		s.ctx.Log.Error("error saving session", L{"err": err})
	}
}

func (s *Session) GetSession() *sessions.Session {
	s.ensureStore()
	sessionName := s.ctx.App().Config.Get("name", "_app_session")
	session, err := s.store.Get(s.ctx.Request, sessionName)
	if err != nil {
		s.ctx.Log.Error("error parsing session", L{"err": err})
	}
	session.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: true,
	}
	return session
}

func (s *Session) Get(key string) string {
	session := s.GetSession()
	if value, ok := session.Values[key]; ok {
		return value.(string)
	}
	return ""
}

func (s *Session) Set(key, value string) {
	session := s.GetSession()
	session.Values[key] = value
}

func (s *Session) AddFlash(kind, message string) {
	session := s.GetSession()
	session.AddFlash(&Flash{Kind: kind, Message: message})
}

func (s *Session) Flashes() []*Flash {
	session := s.GetSession()
	flashes := []*Flash{}
	for _, f := range session.Flashes() {
		flashes = append(flashes, f.(*Flash))
	}
	return flashes
}
