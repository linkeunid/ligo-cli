package middleware

import "github.com/linkeunid/ligo"

// ExceptionMiddleware handles errors and converts them to HTTP responses.
func ExceptionMiddleware(log ligo.Logger) ligo.Middleware {
	return func(next ligo.HandlerFunc) ligo.HandlerFunc {
		return func(ctx ligo.Context) error {
			err := next(ctx)
			if err == nil {
				return nil
			}
			log.Error("Request error",
				ligo.LoggerField{Key: "method", Value: ctx.Request().Method},
				ligo.LoggerField{Key: "path", Value: ctx.Request().URL.Path},
				ligo.LoggerField{Key: "error", Value: err.Error()},
			)
			return ctx.JSON(500, map[string]string{"error": "Internal Server Error"})
		}
	}
}
