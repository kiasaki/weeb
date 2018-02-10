package weeb

import (
	"bytes"
	"context"
	"net/http"

	"github.com/kiasaki/weeb/id"
)

var requestContextKey contextKey = 0

type Context struct {
	app        *App
	statusCode int
	body       string
	Request    *http.Request
	Response   http.ResponseWriter

	Data    J
	DB      DB
	Log     *Logger
	Session *Session
	Auth    *Auth
	ID      *id.Gen
	Mail    Mailer
}

func NewContext(app *App, w http.ResponseWriter, r *http.Request) *Context {
	ctx := &Context{app: app, statusCode: 200}

	ctx.Request = r.WithContext(context.WithValue(r.Context(), requestContextKey, ctx))
	ctx.Response = w

	ctx.Data = J{}
	ctx.DB = app.DB
	ctx.Log = app.Log.WithContext(L{})
	ctx.Session = NewSession(ctx)
	ctx.Auth = app.Auth
	ctx.ID = app.ID
	ctx.Mail = app.Mail

	return ctx
}

func (ctx *Context) App() *App {
	return ctx.app
}

func (ctx *Context) HandleError(err error) {
	if err == nil {
		return
	}
	ctx.Log.Error("request error", L{"err": err.Error()})
	ctx.Error(500, "internal server error")
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

func (ctx *Context) Template(template string, value J) (string, error) {
	data := J{}
	for k, v := range ctx.Data {
		data[k] = v
	}
	for k, v := range value {
		data[k] = v
	}

	var b bytes.Buffer
	err := ctx.app.templates.ExecuteTemplate(&b, template, data)
	return b.String(), err
}
