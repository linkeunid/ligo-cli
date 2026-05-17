package module

import (
	"{{.ModulePath}}/internal/usecase"
	"{{.ModulePath}}/internal/worker"

	"github.com/linkeunid/ligo"
)

// RunnerModule returns the background worker module.
// This is a simple implementation without hooks.
func RunnerModule() ligo.Module {
	return ligo.NewModule(
		"runner",
		ligo.Providers(
			ligo.Factory[*usecase.WorkerUseCase](usecase.NewWorkerUseCase),
		),
		ligo.Controllers(worker.NewController),
	)
}
