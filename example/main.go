package main

import (
	"net/http"

	"github.com/kiasaki/weeb"
	"github.com/kiasaki/weeb/example/migrations"
)

type App struct {
	*weeb.App
}

func main() {
	app := &App{weeb.NewApp()}
	migrations.AddMigrationsToApp(app.App)
	app.Router.HandleFunc("/", app.handleHome)
	app.Run()
}

func (app *App) handleHome(w http.ResponseWriter, r *http.Request) {
	app.SendJSON(w, 200, weeb.J{
		"message": "Hello World!",
	})
}
