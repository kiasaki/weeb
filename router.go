package weeb

import (
	"net/http"

	"github.com/gorilla/mux"
)

type HandlerFunc func(*Context) error

type Router struct {
	app    *App
	router *mux.Router
}

func NewRouter(app *App) *Router {
	return &Router{app: app, router: mux.NewRouter()}
}

func (r *Router) Group(prefix string) *Router {
	subRouter := r.router.PathPrefix(prefix).Subrouter()
	return &Router{app: r.app, router: subRouter}
}

func (r *Router) UseHTTP(middleware func(http.Handler) http.Handler) {
	r.router.Use(middleware)
}

func (r *Router) Use(middleware func(HandlerFunc) HandlerFunc) {
	r.router.Use(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			ctx := r.requestContext(w, req)
			next := func(c *Context) error {
				h.ServeHTTP(c.Response, c.Request)
				return nil
			}
			ctx.HandleError(middleware(next)(ctx))
		})
	})
}

func (r *Router) Static(prefix, dir string) {
	staticFilesHandler := http.StripPrefix(prefix, http.FileServer(http.Dir(dir)))
	r.router.PathPrefix(prefix).Handler(staticFilesHandler)
}

func (r *Router) Handle(method, path string, handler HandlerFunc) {
	r.router.HandleFunc(path, func(w http.ResponseWriter, req *http.Request) {
		ctx := r.requestContext(w, req)
		ctx.HandleError(handler(ctx))
	}).Methods(method)
}

func (r *Router) Head(path string, handler HandlerFunc) {
	r.Handle("HEAD", path, handler)
}

func (r *Router) Options(path string, handler HandlerFunc) {
	r.Handle("OPTIONS", path, handler)
}

func (r *Router) Get(path string, handler HandlerFunc) {
	r.Handle("GET", path, handler)
}

func (r *Router) Post(path string, handler HandlerFunc) {
	r.Handle("POST", path, handler)
}

func (r *Router) Put(path string, handler HandlerFunc) {
	r.Handle("PUT", path, handler)
}

func (r *Router) Patch(path string, handler HandlerFunc) {
	r.Handle("PATCH", path, handler)
}

func (r *Router) Delete(path string, handler HandlerFunc) {
	r.Handle("DELETE", path, handler)
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.router.ServeHTTP(w, req)
}

func (r *Router) requestContext(w http.ResponseWriter, req *http.Request) *Context {
	ctx, ok := req.Context().Value(authUserKey).(*Context)
	if ok {
		return ctx
	}
	return NewContext(r.app, w, req)
}