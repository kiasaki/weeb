package weeb

import (
	"time"
)

func loggingMiddleware(next HandlerFunc) HandlerFunc {
	return func(ctx *Context) error {
		start := time.Now()
		if err := next(ctx); err != nil {
			return err
		}
		code := ctx.Response.(*responseWriterWithStatusCode).statusCode
		if code == 0 {
			code = ctx.StatusCode()
		}
		ctx.Log.Info(ctx.Request.URL.Path, L{
			"method": ctx.Request.Method,
			"code":   code,
			"ms":     time.Now().Unix() - start.Unix(),
		})
		return nil
	}
}
