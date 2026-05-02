package middleware

import (
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/linkeunid/ligo"
	"{{.ModulePath}}/internal/usecase"
)

type fieldErr struct {
	Field string `json:"field"`
	Tag   string `json:"tag"`
	Param string `json:"param,omitempty"`
}

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
				return ctx.Unauthorized()
			case errors.Is(err, usecase.ErrForbidden):
				return ctx.Forbidden()
			case errors.Is(err, usecase.ErrNotFound):
				return ctx.NotFound()
			case errors.As(err, &ve):
				errs := make([]fieldErr, 0, len(ve))
				for _, fe := range ve {
					errs = append(errs, fieldErr{Field: fe.Field(), Tag: fe.Tag(), Param: fe.Param()})
				}
				return ctx.JSON(http.StatusUnprocessableEntity, map[string]any{"errors": errs})
			case errors.Is(err, ligo.ErrBadRequest):
				return ctx.BadRequest(err.Error())
			case errors.Is(err, usecase.ErrValidation):
				return ctx.BadRequest()
			default:
				return ctx.InternalServerError()
			}
		}
	}
}
