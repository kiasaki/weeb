package weeb

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/kiasaki/weeb/id"
	"github.com/mitchellh/mapstructure"
)

var requestContextKey contextKey = 0

type Context struct {
	app        *App
	statusCode int
	body       string
	Request    *http.Request
	Response   http.ResponseWriter

	Data     J
	DB       DB
	DBHelper *DBHelper
	Log      *Logger
	Session  *Session
	Config   *Config
	Auth     *Auth
	ID       *id.Gen
	Mail     Mailer
}

func NewHTTPContext(app *App, w http.ResponseWriter, r *http.Request) *Context {
	ctx := NewContext(app)

	ctx.Request = r.WithContext(context.WithValue(r.Context(), requestContextKey, ctx))
	ctx.Response = w

	return ctx
}

func NewContext(app *App) *Context {
	ctx := &Context{app: app, statusCode: 200}

	ctx.Data = J{}
	ctx.DB = app.DB
	ctx.DBHelper = NewDBHelper(app.DB)
	ctx.Log = app.Log.WithContext(L{})
	ctx.Session = NewSession(ctx)
	ctx.Config = app.Config
	ctx.Auth = app.Auth
	ctx.ID = app.ID
	ctx.Mail = app.Mail

	return ctx
}

func (ctx *Context) App() *App {
	return ctx.app
}

func (ctx *Context) HandleError(err error) error {
	if err == nil {
		return nil
	}
	ctx.Log.Error("request error", L{"err": err.Error()})
	ctx.Error(500, "internal server error")
	return nil
}

func (ctx *Context) Error(code int, message string) error {
	if handlerFn, ok := ctx.app.Router.ErrorHandlers[code]; ok {
		err := handlerFn(ctx)
		// handle error ourselves to avoid HandleError -> Error -> HandleError loops
		if err != nil {
			ctx.Log.Error("request error", L{"err": err.Error()})
			ctx.Text(500, "internal server error")
		}
		return nil
	}
	return ctx.Text(code, message)
}

func (ctx *Context) Redirect(url string) error {
	return ctx.RedirectWithCode(http.StatusFound, url)
}

func (ctx *Context) RedirectWithCode(code int, url string) error {
	ctx.SetHeader("Location", hexEscapeNonASCII(url))
	ctx.SetStatusCode(code)
	return nil
}

func (ctx *Context) Get(key string) interface{} {
	return ctx.Data[key]
}

func (ctx *Context) Set(key string, value interface{}) {
	ctx.Data[key] = value
}

func (ctx *Context) Param(key string, alt string) string {
	urlVars := ctx.Data["vars"].(map[string]string)
	if value, ok := urlVars[key]; ok {
		return value
	}
	value := ctx.Request.FormValue(key)
	if value == "" {
		return alt
	}
	return value
}

func (ctx *Context) SetHeader(name, value string) {
	ctx.Response.Header().Set(name, value)
}

func (ctx *Context) SetStatusCode(code int) {
	ctx.statusCode = code
}

func (ctx *Context) SetBody(body string) {
	ctx.body = body
}

func (ctx *Context) StatusCode() int {
	return ctx.statusCode
}

func (ctx *Context) finalizeResponse() {
	ctx.Session.save()
	ctx.Response.WriteHeader(ctx.statusCode)
	ctx.Response.Write([]byte(ctx.body))
}

func (ctx *Context) Template(name string, value J) (string, error) {
	data := J{}
	for k, v := range ctx.Data {
		data[k] = v
	}
	for k, v := range value {
		data[k] = v
	}
	return ctx.app.Templates.Render(name, data)
}

// Bind parses the request body into a given entity
func (ctx *Context) Bind(entity interface{}) error {
	defer ctx.Request.Body.Close()
	if strings.Contains(ctx.Request.Header.Get("Content-Type"), "application/x-www-form-urlencoded") {
		formData := map[string]interface{}{}
		err := ctx.Request.ParseForm()
		if err != nil {
			ctx.Log.Error("error parsing form", L{"err": err.Error()})
		}

		for key := range ctx.Request.Form {
			formData[key] = ctx.Request.Form.Get(key)
		}
		return mapstructure.Decode(formData, entity)
	} else if strings.Contains(ctx.Request.Header.Get("Content-Type"), "application/json") {
		return json.NewDecoder(ctx.Request.Body).Decode(entity)
	} else {
		return errors.New("Unsupported Content-Type provided")
	}
}
