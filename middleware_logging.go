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
		ctx.Log.Info(ctx.Request.URL.Path, L{
			"method": ctx.Request.Method,
			"code":   ctx.StatusCode(),
			"ms":     time.Now().Unix() - start.Unix(),
		})
		return nil
	}
}
