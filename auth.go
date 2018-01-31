package weeb

import "time"

var authUserKey contextKey = 1

type AuthUser interface {
	AuthID() string
	AuthUsername() string
	AuthPassword() string
	AuthRoles() []string
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
			return nil
		}
	}
}

func (a *Auth) CurrentUser(ctx *Context) (AuthUser, error) {
	if user, ok := ctx.Data["currentUser"]; ok {
		return user.(AuthUser), nil
	}

	//userID := ctx.Session.Get("userId")

	return nil, nil
}

type User struct {
	ID       string
	Username string
	Password string
	Roles    DBStringArray
	Created  time.Time
	Updated  time.Time
}

func (u *User) Table() string {
	return "weeb_users"
}

func (u *User) AuthID() string {
	return u.ID
}

func (u *User) AuthUsername() string {
	return u.Username
}

func (u *User) AuthPassword() string {
	return u.Password
}

func (u *User) AuthRoles() []string {
	return u.Roles
}

func authDefaultFindByID(ctx *Context, id string) (AuthUser, error) {
	var user User
	err := ctx.DB.QueryOne(&user, "SELECT * FROM weeb_users WHERE id = $1", id)
	return &user, err
}

func authDefaultFindByUsername(ctx *Context, id string) (AuthUser, error) {
	var user User
	err := ctx.DB.QueryOne(&user, "SELECT * FROM weeb_users WHERE username = $1", id)
	return &user, err
}
