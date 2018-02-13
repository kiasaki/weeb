# weeb

_A weeb framework. For quickly building ambitious weeb applications._

### intro

### features

**there**

- Cache
- Templates
- Routing
- Middewares
- Sessions
- Mails
- Logging
- Authentication
- Database Querying
- Database Migrations

**upcomming**

- Database CRUD
- Encryption
- I18n
- Validation
- Background Jobs
- Cron Jobs
- Deployment
- File Storage
- Security
- Test Helpers

**offered as plugins**

- Auth (w/ signin, signup pages and users table)
- Admin (to manage arbitrary entity `a la` django admin)
- Forms (e.g. contact form, feedback form)
- Subscriptions Billing (manage stripe subscriptions and offer a billing page)

### usage

### concepts

The `weeb.App` object is god. It's passed around everywhere and is the gateway
to most services exposed by weeb.

### starting a new project

```
go get github.com/kiasaki/weeb
```

Create an `main.go` file with:

```js
package main

import (
	"github.com/kiasaki/weeb"
)

func main() {
	app := weeb.NewApp()

	app.Router.ErrorHandlers[404] = handle404
	app.Router.Get("/", handleHome)
	app.Router.Post("/mail", handleMail)

	// **Route groups:**
	r := app.Router.Group("/app/")
	// **Authorization:**
	r.Use(app.Auth.RequireRoles("user"))
	app.Router.Get("/api/", handleApi)

	app.Tasks.Register("say-hello", tasksSayHello)

	app.Run()
}

func tasksSayHello(app *weeb.App, args []string) error {
	fmt.Println("Hello!")
	return nil
}

func handle404(ctx *weeb.Context) error {
	// **Rendering an html template from `templates/*.tmpl`:**
	return ctx.HTML(404, "404", weeb.J{"title": "Not found"})
}

func handleHome(ctx *weeb.Context) error {
	// **Rendering text:**
	return ctx.Text(200, "Hello World!")
}

func handleApi(ctx *weeb.Context) error {
	// **Rendering json:**
	return ctx.JSON(200, weeb.J{
		"version": "3.14",
	})
}

func handleMail(ctx *weeb.Context) error {
	// **Form/Query params with defaults:**
	name := ctx.Param("name", "Unknown")

	// **Structured logging:**
	ctx.Log.Info("sending email", weeb.L{"name": name})

	// **Rendering templates to strings:**
	message, err := ctx.Template("mails/hi", weeb.J{"name": name})
	if err != nil {
		return ctx.HandleError(err)
	}

	// **Sending emails:**
	err = ctx.Mail.Send("from@me.com", "recipient@email.com", "Hi", message)
	if err != nil {
		return ctx.HandleError(err)
	}

	// **Redirects:**
	return ctx.Redirect("/success")
}
```

Start the web server by running:

```
$ go run main.go
{"time": "2017-01-01T00:00:00.000Z", "level": "info", "msg": "starting", "port": "3000"}
```

Check the available command line tasks using:

```
$ go run main.go help
Available commands:

    config
    dev
    generate-session-key
    help
    migrate
    start

```

### license

MIT. See `LICENSE` file.
