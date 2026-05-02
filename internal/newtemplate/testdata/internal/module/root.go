package module

import (
	"github.com/linkeunid/ligo"
	"{{.ModulePath}}/internal/infrastructure/http/controller"
)

// Root returns the root API info module.
func RootModule() ligo.Module {
	return ligo.NewModule("root",
		ligo.Controllers(controller.NewRootController),
	)
}
