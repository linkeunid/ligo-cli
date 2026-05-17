package auth

import (
	"strings"

	"{{.ModulePath}}/internal/domain/service"
	"{{.ModulePath}}/internal/usecase"

	"github.com/linkeunid/ligo"
)

// AuthGuard creates an authentication guard that validates bearer tokens.
func AuthGuard(authService service.AuthService) ligo.Guard {
	return func(ctx ligo.Context) (bool, error) {
		authHeader := ctx.Request().Header.Get("Authorization")
		if authHeader == "" {
			return false, usecase.ErrUnauthorized
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			return false, usecase.ErrUnauthorized
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")

		user, err := authService.ValidateToken(token)
		if err != nil {
			return false, usecase.ErrUnauthorized
		}

		ctx.Set(service.ContextKeyUser, user)
		return true, nil
	}
}

// AdminGuard creates an admin-only authorization guard.
func AdminGuard() ligo.Guard {
	return func(ctx ligo.Context) (bool, error) {
		user, ok := ctx.Get(service.ContextKeyUser).(service.User)
		if !ok || user == nil {
			return false, usecase.ErrUnauthorized
		}

		if !user.IsAdmin() {
			return false, usecase.ErrForbidden
		}

		return true, nil
	}
}
