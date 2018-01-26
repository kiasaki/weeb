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
	app.Run()
}

func handleHome(ctx *weeb.Context) error {
	return ctx.SendHTML(200, "home", weeb.J{
		"message": "Hello World!",
	})
}

func handleHello(ctx *weeb.Context) error {
	return ctx.SendJSON(200, weeb.J{
		"message": "Hello World!",
	})
}