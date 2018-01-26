package weeb

import (
	"errors"
	"fmt"
)

func recoverMiddleware(next HandlerFunc) HandlerFunc {
	return func(ctx *Context) error {
		defer func() {
			if err := recover(); err != nil {
				ctx.HandleError(errors.New(fmt.Sprintf("panic: %v", err)))
			}
		}()
		return next(ctx)
	}
}
