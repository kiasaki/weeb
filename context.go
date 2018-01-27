package weeb

import (
	"context"
	"net/http"
)

var requestContextKey contextKey = 0

type Context struct {
	app        *App
	statusCode int
	body       string
	Request    *http.Request
	Response   http.ResponseWriter

	Data    J
	Log     *Logger
	Session *Session
}

func NewContext(app *App, w http.ResponseWriter, r *http.Request) *Context {
	ctx := &Context{app: app, statusCode: 200}

	ctx.Request = r.WithContext(context.WithValue(r.Context(), requestContextKey, ctx))
	ctx.Response = w

	ctx.Data = J{}
	ctx.Log = app.Log.WithContext(L{})
	ctx.Session = NewSession(ctx)

	return ctx
}

func (ctx *Context) App() *App {
	return ctx.app
}

func (ctx *Context) HandleError(err error) {
	if err == nil {
		return
	}
	ctx.Text(500, "internal server error")
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
