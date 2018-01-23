# weeb

_A weeb framework. For quickly building ambitious weeb applications._

### intro

### features

**there**

- Authentication
- Controllers
- DI / IOC Container
- Database CRUD
- Database Querying
- Encryption
- I18n
- Logging
- Mails
- Middewares
- Pluggable Providers
- Routing
- Templates
- Validation

**upcomming**

- Background Jobs
- Cache
- Cron Jobs
- Database Migrations
- Deployment
- File Storage
- Security
- Sessions
- Subscriptions Billing
- Test Helpers

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
	"net/http"

	"github.com/kiasaki/weeb"
)

type App struct {
	*weeb.App
}

func main() {
	app := &App{weeb.NewApp()}
	app.Router.HandleFunc("/", app.handleHome)
	app.Start()
}

func (app *App) handleHome(w http.ResponseWriter, r *http.Request) {
	app.SendJSON(w, 200, weeb.J{
		"message": "Hello World!",
	})
}
```

Start the web server by running:

```
$ go run main.go
{"time": "2017-01-01T00:00:00.000Z", "level": "info", "msg": "starting", "port": "3000"}
```

### license

MIT. See `LICENSE` file.
