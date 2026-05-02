package module

import (
	"github.com/linkeunid/ligo"
	ligomemory "github.com/linkeunid/ligo-memory"
	"{{.ModulePath}}/internal/domain/entity"
	"{{.ModulePath}}/internal/domain/repository"
	"{{.ModulePath}}/internal/infrastructure/http/controller"
	"{{.ModulePath}}/internal/infrastructure/persistence/memory"
	"{{.ModulePath}}/internal/usecase"
)

// User returns the user module.
func UserModule() ligo.Module {
	return ligo.NewModule("user",
		ligo.Providers(
			ligomemory.Provider[int, *entity.User](),
			ligo.Factory[repository.UserRepository](memory.NewUserRepository),
			ligo.Factory[*usecase.UserUseCase](usecase.NewUserUseCase),
		),
		ligo.Controllers(controller.NewUserController),
	)
}
