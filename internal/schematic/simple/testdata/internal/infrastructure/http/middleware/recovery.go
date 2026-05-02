package middleware

import "github.com/linkeunid/ligo"

// Recovery returns a middleware that catches panics and responds with 500.
func Recovery() ligo.Middleware {
	return func(next ligo.HandlerFunc) ligo.HandlerFunc {
		return func(ctx ligo.Context) error {
			defer func() {
				if r := recover(); r != nil {
					ctx.Response().WriteHeader(500)
				}
			}()
			return next(ctx)
		}
	}
}
