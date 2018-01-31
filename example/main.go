package main

import (
	"github.com/kiasaki/weeb"
	"github.com/kiasaki/weeb/example/migrations"
)

func main() {
	app := weeb.NewApp()
	migrations.AddMigrationsToApp(app)

	app.Router.ErrorHandlers[401] = handleRedirectToLogin
	app.Router.ErrorHandlers[403] = handleRedirectToLogin

	app.Router.Get("/", handleHome)
	app.Router.Get("/hello", handleHello)
	app.Router.Get("/signin", handleSignin)

	admin := app.Router.Group("/admin")
	admin.Use(app.Auth.RequireRoles("admin"))
	admin.Get("/", handleAdminHome)

	app.Run()
}

func handleRedirectToLogin(ctx *weeb.Context) error {
	return ctx.Redirect("/signin")
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
	return ctx.Text(200, "signin")
}

func handleAdminHome(ctx *weeb.Context) error {
	return ctx.Text(200, "admin")
}
