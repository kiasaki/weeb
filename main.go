package main

import "net/http"

func main() {
	app := NewApp()
	r := app.Router

	r.HandleFunc("/", app.handleHome)

	app.Start()
}

func (app *App) handleHome(w http.ResponseWriter, r *http.Request) {
	app.SendText(w, 200, "Hello World!")
}
