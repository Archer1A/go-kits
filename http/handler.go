package remote

import (
	"fmt"
	"runtime/debug"
)

// Recovery returns a middleware that recovery any panic fired by pending middlewares and save it as error
// in Response
func Recovery() HandlerFunc {
	return func(ctx *Context) {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("panic recovered\nstacktrace: \n" + string(debug.Stack()))
				var ok bool
				var err error
				err, ok = r.(error)
				if !ok {
					err = fmt.Errorf("%v", r)
				}
				ctx.Abort()
				ctx.Response.ErrorSave(err)
			}
		}()
		ctx.Next()
	}
}
