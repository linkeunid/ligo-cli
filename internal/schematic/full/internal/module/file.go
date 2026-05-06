package module

import (
	"github.com/linkeunid/ligo"
	ligomemory "github.com/linkeunid/ligo-memory"
	"{{.ModulePath}}/internal/config"
	"{{.ModulePath}}/internal/domain/entity"
	"{{.ModulePath}}/internal/domain/repository"
	"{{.ModulePath}}/internal/infrastructure/http/controller"
	"{{.ModulePath}}/internal/infrastructure/persistence/memory"
	"{{.ModulePath}}/internal/usecase"
)

// File returns the file upload module.
func FileModule() ligo.Module {
	return ligo.NewModule("file",
		ligo.Providers(
			ligomemory.Provider[int, *entity.File](),
			ligo.Factory[repository.FileRepository](func(cfg *config.Config, store *ligomemory.Store[int, *entity.File]) repository.FileRepository {
				return memory.NewFileRepository(cfg.UploadDir, store)
			}),
			ligo.Factory[*usecase.FileUseCase](usecase.NewFileUseCase),
		),
		ligo.Controllers(controller.NewFileController),
	)
}
