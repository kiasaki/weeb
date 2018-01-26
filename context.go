package weeb

import (
	"context"
	"net/http"
)

var requestContextKey contextKey = 0

// ResponseWriter wrapper that keeps track of the status code we sent
type responseWriterWithStatus struct {
	http.ResponseWriter
	code int
}

func (w *responseWriterWithStatus) WriteHeader(code int) {
	w.code = code
	w.ResponseWriter.WriteHeader(code)
}
func (w *responseWriterWithStatus) Status() int { return w.code }

type Context struct {
	app      *App
	Request  *http.Request
	Response http.ResponseWriter

	Log *Logger
}

func NewContext(app *App, w http.ResponseWriter, r *http.Request) *Context {
	ctx := &Context{app: app}

	ctx.Request = r.WithContext(context.WithValue(r.Context(), requestContextKey, ctx))
	ctx.Response = &responseWriterWithStatus{w, 0}

	ctx.Log = app.Log.WithContext(L{})

	return ctx
}

func (ctx *Context) HandleError(err error) {
	if err == nil {
		return
	}
	ctx.SendText(500, "internal server error")
}

func (ctx *Context) StatusCode() int {
	return ctx.Response.(*responseWriterWithStatus).Status()
}
