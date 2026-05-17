package middleware

import (
	"{{.ModulePath}}/internal/domain/service"

	"github.com/linkeunid/ligo"
)

// AuditMiddleware logs admin actions for security auditing.
func AuditMiddleware(log ligo.Logger) ligo.Middleware {
	return func(next ligo.HandlerFunc) ligo.HandlerFunc {
		return func(ctx ligo.Context) error {
			user, _ := ctx.Get(service.ContextKeyUser).(service.User)

			err := next(ctx)

			if user != nil && user.IsAdmin() {
				log.Info(
					"Admin action performed",
					ligo.LoggerField{Key: "admin_id", Value: user.GetID()},
					ligo.LoggerField{Key: "action", Value: ctx.Request().Method},
					ligo.LoggerField{Key: "path", Value: ctx.Request().URL.Path},
					ligo.LoggerField{Key: "success", Value: err == nil},
				)
			}

			return err
		}
	}
}
