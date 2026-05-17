package module

import (
	"{{.ModulePath}}/internal/infrastructure/http/controller"
	"{{.ModulePath}}/internal/usecase"

	"github.com/linkeunid/ligo"
)

// Hello returns the hello module.
func HelloModule() ligo.Module {
	return ligo.NewModule(
		"hello",
		ligo.Providers(
			ligo.Factory[*usecase.HelloUseCase](usecase.NewHelloUseCase),
		),
		ligo.Controllers(controller.NewHelloController),
	)
}
