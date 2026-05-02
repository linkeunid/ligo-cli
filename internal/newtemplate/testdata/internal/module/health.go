package module

import (
	"github.com/linkeunid/ligo"
	"{{.ModulePath}}/internal/infrastructure/http/controller"
)

// Health returns the health check module.
func Health() ligo.Module {
	return ligo.NewModule("health",
		ligo.Controllers(controller.NewHealthController),
	)
}
