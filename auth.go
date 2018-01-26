package weeb

import "net/http"

var authUserKey contextKey = 1

type AuthUser interface {
	ID() string
	Username() string
	Password() string
	Roles() []string
}

type Auth struct {
	app *App
}

func NewAuth(app *App) *Auth {
	return &Auth{app: app}
}

func (a *Auth) CurrentUser(r *http.Request) AuthUser {
	ctx := r.Context()
	user, ok := ctx.Value(authUserKey).(AuthUser)
	if ok {
		return user
	}

	//userID := a.app.Session.Get(r, "")
	return nil
}
