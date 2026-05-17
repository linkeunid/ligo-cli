package middleware

import "github.com/linkeunid/ligo"

// CORS returns a middleware that sets permissive CORS headers and handles preflight requests.
func CORS() ligo.Middleware {
	return func(next ligo.HandlerFunc) ligo.HandlerFunc {
		return func(ctx *ligo.Context) error {
			ctx.Response().Header().Set("Access-Control-Allow-Origin", "*")
			ctx.Response().Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			ctx.Response().Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

			if ctx.Request().Method == "OPTIONS" {
				return ctx.String(204, "")
			}

			return next(ctx)
		}
	}
}
