package weeb

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

var authUserKey contextKey = 1

var ErrorUserNotFound = errors.New("User not found")
var ErrorPasswordsDontMatch = errors.New("Passwords don't match")

type AuthUser interface {
	AuthID() string
	AuthUsername() string
	AuthPassword() string
	AuthRoles() []string
}

type AuthSigninInfo struct {
	Username     string
	Password     string
	OnlyValidate bool
}

type Auth struct {
	app            *App
	FindByID       func(ctx *Context, id string) (AuthUser, error)
	FindByUsername func(ctx *Context, username string) (AuthUser, error)
}

func NewAuth(app *App) *Auth {
	return &Auth{
		app:            app,
		FindByID:       authDefaultFindByID,
		FindByUsername: authDefaultFindByUsername,
	}
}

func (a *Auth) RequireRoles(roles ...string) func(HandlerFunc) HandlerFunc {
	return func(next HandlerFunc) HandlerFunc {
		return func(ctx *Context) error {
			user, err := ctx.Auth.CurrentUser(ctx)
			if err != nil {
				return err
			}
			if user == nil {
				return ctx.Error(401, "unauthorized")
			}
			userRoles := user.AuthRoles()
			for _, role := range roles {
				if !containsString(userRoles, role) {
					return ctx.Error(403, "forbidden")
				}
			}

			return next(ctx)
		}
	}
}

func (a *Auth) Signin(ctx *Context, info AuthSigninInfo) error {
	user, err := a.FindByUsername(ctx, info.Username)
	if err != nil {
		return err
	}
	if user == nil {
		return ErrorUserNotFound
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.AuthPassword()), []byte(info.Password))
	if err != nil {
		return ErrorPasswordsDontMatch
	}

	if !info.OnlyValidate {
		a.SigninUser(ctx, user)
	}

	return nil
}

func (a *Auth) SigninUser(ctx *Context, user AuthUser) {
	ctx.Set("currentUser", user)
	ctx.Session.Set("userID", user.AuthID())
}

func (a *Auth) Signout(ctx *Context) {
	ctx.Session.Set("userID", "")
}

func (a *Auth) CurrentUser(ctx *Context) (AuthUser, error) {
	if user, ok := ctx.Data["currentUser"]; ok {
		return user.(AuthUser), nil
	}

	userID := ctx.Session.Get("userID")
	if userID == "" {
		return nil, nil
	}

	user, err := a.FindByID(ctx, userID)
	if err == nil && user != nil {
		ctx.Set("currentUser", user)
	}
	return user, err
}

func authDefaultFindByID(ctx *Context, id string) (AuthUser, error) {
	panic("Auth: FindByID was not configured")
}

func authDefaultFindByUsername(ctx *Context, id string) (AuthUser, error) {
	panic("Auth: FindByUsername was not configured")
}
