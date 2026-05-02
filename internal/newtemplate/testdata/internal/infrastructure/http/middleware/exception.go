package middleware

import (
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/linkeunid/ligo"
	"{{.ModulePath}}/internal/usecase"
)

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

			var ve validator.ValidationErrors
			switch {
			case errors.Is(err, usecase.ErrUnauthorized):
				return ctx.JSON(401, map[string]string{"error": "Unauthorized"})
			case errors.Is(err, usecase.ErrForbidden):
				return ctx.JSON(403, map[string]string{"error": "Forbidden"})
			case errors.Is(err, usecase.ErrNotFound):
				return ctx.JSON(404, map[string]string{"error": "Not Found"})
			case errors.Is(err, ligo.ErrBadRequest):
				return ctx.JSON(400, map[string]string{"error": "Bad Request"})
			case errors.As(err, &ve):
				return ctx.JSON(422, map[string]string{"error": "Unprocessable Entity"})
			case errors.Is(err, usecase.ErrValidation):
				return ctx.JSON(400, map[string]string{"error": "Bad Request"})
			default:
				return ctx.JSON(500, map[string]string{"error": "Internal Server Error"})
			}
		}
	}
}
