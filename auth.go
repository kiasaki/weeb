package weeb

import (
	"errors"
	"strconv"
	"time"

	"github.com/lib/pq"

	"golang.org/x/crypto/bcrypt"
)

var authUserKey contextKey = 1

var ErrorPasswordsDontMatch = errors.New("Passwords don't match")

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

type AuthSigninInfo struct {
	Username     string
	Password     string
	OnlyValidate bool
}

func (a *Auth) Signin(ctx *Context, info AuthSigninInfo) error {
	user, err := a.FindByUsername(ctx, info.Username)
	if err != nil {
		return err
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

func (a *Auth) CreateUser(ctx *Context, user *User) error {
	user.ID = ctx.ID.Next()
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)
	if user.Roles == nil {
		user.Roles = []string{"user"}
	}
	insertSQL := `INSERT INTO weeb_users (id, name, username, password, roles) VALUES (:id, :name, :username, :password, :roles)`
	return ctx.DB.ExecNamed(insertSQL, user)
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

type User struct {
	ID       int64
	Name     string
	Username string
	Password string
	Roles    pq.StringArray
	Created  time.Time
	Updated  time.Time
}

func (u *User) Table() string {
	return "weeb_users"
}

func (u *User) AuthID() string {
	return strconv.FormatInt(u.ID, 10)
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
