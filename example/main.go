package main

import (
	"regexp"
	"strings"

	"github.com/kiasaki/weeb"
	"github.com/kiasaki/weeb/example/migrations"
)

var emailRegexp = regexp.MustCompile(`^[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,}$`)

func main() {
	app := weeb.NewApp()
	migrations.AddMigrationsToApp(app)

	app.Router.ErrorHandlers[401] = handleRedirectToLogin
	app.Router.ErrorHandlers[403] = handleRedirectToLogin
	app.Router.ErrorHandlers[404] = handle404
	app.Router.ErrorHandlers[500] = handle500

	app.Router.Get("/", handleHome)

	app.Router.Get("/signin/", handleSignin)
	app.Router.Post("/signin/", handleSignin)
	app.Router.Get("/signup/", handleSignup)
	app.Router.Post("/signup/", handleSignup)
	app.Router.Get("/signout/", handleSignout)

	r := app.Router.Group("/app/")
	r.Use(app.Auth.RequireRoles("user"))
	r.Get("/", handleApp)

	admin := app.Router.Group("/admin/")
	admin.Use(app.Auth.RequireRoles("admin"))
	admin.Get("/", handleAdminHome)

	app.Run()
}

func handleRedirectToLogin(ctx *weeb.Context) error {
	return ctx.Redirect("/signin/")
}

func handle404(ctx *weeb.Context) error {
	return ctx.HTML(404, "404", weeb.J{"title": "Not found"})
}

func handle500(ctx *weeb.Context) error {
	return ctx.HTML(500, "500", weeb.J{"title": "Server error"})
}

func handleHome(ctx *weeb.Context) error {
	return ctx.HTML(200, "home", weeb.J{
		"message": "Hello World!",
	})
}

func handleHello(ctx *weeb.Context) error {
	return ctx.JSON(200, weeb.J{
		"message": "Hello World!",
	})
}

func handleSignin(ctx *weeb.Context) error {
	var err error
	var username, password string
	if ctx.Request.Method == "POST" {
		username = ctx.Param("username", "")
		password = ctx.Param("password", "")

		err = ctx.Auth.Signin(ctx, weeb.AuthSigninInfo{
			Username: username,
			Password: password,
		})
		if err == nil {
			return ctx.Redirect("/app/")
		}
	}

	return ctx.HTML(200, "signin", weeb.J{
		"title":    "Login",
		"hasError": err != nil,
		"username": username,
	})
}

func handleSignup(ctx *weeb.Context) error {
	var errorMessage string
	var user weeb.User
	if ctx.Request.Method == "POST" {
		user = weeb.User{
			Name:     ctx.Param("name", ""),
			Username: ctx.Param("email", ""),
			Password: ctx.Param("password", ""),
		}
		user.Username = strings.ToLower(user.Username)
		if len(user.Name) == 0 {
			errorMessage = "Missing full name"
			goto render
		}
		if !emailRegexp.MatchString(user.Username) {
			errorMessage = "Invalid email provided"
			goto render
		}
		if len(user.Password) < 8 {
			errorMessage = "A password of at least 8 characters is required"
			goto render
		}
		if user.Password != ctx.Param("passwordConfirmation", "") {
			errorMessage = "Password confirmation doesn't match"
			goto render
		}
		err := ctx.Auth.CreateUser(ctx, &user)
		if err != nil {
			errorMessage = "A server error occured while creating account"
			goto render
		}

		ctx.Auth.SigninUser(ctx, &user)
		return ctx.Redirect("/app/")
	}

render:
	return ctx.HTML(200, "signup", weeb.J{
		"title":    "Sign Up",
		"hasError": errorMessage != "",
		"error":    errorMessage,
		"name":     user.Name,
		"email":    user.Username,
	})
}

func handleSignout(ctx *weeb.Context) error {
	ctx.Auth.Signout(ctx)
	return ctx.Redirect("/")
}

func handleApp(ctx *weeb.Context) error {
	return ctx.Text(200, "app")
}

func handleAdminHome(ctx *weeb.Context) error {
	return ctx.Text(200, "admin")
}
