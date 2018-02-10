package main

import (
	"regexp"

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
	app.Router.Get("/mail", handleMail)

	/*
		r := app.Router.Group("/app/")
		r.Use(app.Auth.RequireRoles("user"))
		r.Get("/", handleApp)

		admin := app.Router.Group("/admin/")
		admin.Use(app.Auth.RequireRoles("admin"))
		admin.Get("/", handleAdminHome)
	*/

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

func handleMail(ctx *weeb.Context) error {
	message, err := ctx.Template("mail_hi", weeb.J{"name": "Fred"})
	if err != nil {
		return ctx.Text(500, "error rendering email template")
	}
	err = ctx.Mail.Send("", "frederic@gingras.cc", "Hi", message)
	if err != nil {
		return ctx.Text(500, "error sending email")
	}
	return ctx.Text(200, "emailed")
}
