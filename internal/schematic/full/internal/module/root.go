package module

import (
	"{{.ModulePath}}/internal/infrastructure/http/controller"

	"github.com/linkeunid/ligo"
)

// Root returns the root API info module.
// This module does not use hooks - it's a simple example without lifecycle management.
func RootModule() ligo.Module {
	return ligo.NewModule(
		"root",
		ligo.Controllers(controller.NewRootController),
	)
}
