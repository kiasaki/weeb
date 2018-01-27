package weeb

import (
	"fmt"
)

func recoverMiddleware(next HandlerFunc) HandlerFunc {
	return func(ctx *Context) error {
		defer func() {
			if err := recover(); err != nil {
				ctx.HandleError(fmt.Errorf("panic: %v", err))
			}
		}()
		return next(ctx)
	}
}
