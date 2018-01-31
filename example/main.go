package main

import (
	"github.com/kiasaki/weeb"
	"github.com/kiasaki/weeb/example/migrations"
)

func main() {
	app := weeb.NewApp()
	migrations.AddMigrationsToApp(app)

	app.Router.Get("/", handleHome)
	app.Router.Get("/hello", handleHello)

	admin := app.Router.Group("/admin")
	admin.Use(app.Auth.RequireRoles("admin"))

	app.Run()
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
