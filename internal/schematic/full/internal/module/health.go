package module

import (
	"{{.ModulePath}}/internal/infrastructure/http/controller"

	"github.com/linkeunid/ligo"
)

// Health returns the health check module.
// This module does not use hooks - it's a simple example without lifecycle management.
func HealthModule() ligo.Module {
	return ligo.NewModule(
		"health",
		ligo.Controllers(controller.NewHealthController),
	)
}
