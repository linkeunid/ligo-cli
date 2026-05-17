package module

import (
	"{{.ModulePath}}/internal/domain/entity"
	"{{.ModulePath}}/internal/domain/repository"
	"{{.ModulePath}}/internal/infrastructure/http/controller"
	"{{.ModulePath}}/internal/infrastructure/http/presenter"
	"{{.ModulePath}}/internal/infrastructure/persistence/memory"
	"{{.ModulePath}}/internal/usecase"

	"github.com/linkeunid/ligo"
	ligomemory "github.com/linkeunid/ligo-memory"
)

// User returns the user module with compile-time safe hook registration.
func UserModule() ligo.Module {
	return ligo.NewModule(
		"user",
		ligo.Providers(
			ligomemory.Provider[int, *entity.User](),
			ligo.Factory[*presenter.UserPresenter](presenter.NewUserPresenter),
			// Use HookedFactory for compile-time safe hook registration.
			// The UserRepository.Register method will be called automatically
			// to register its lifecycle hooks.
			ligo.HookedFactory[repository.UserRepository](memory.NewUserRepository),
			ligo.Factory[*usecase.UserUseCase](usecase.NewUserUseCase),
		),
		// Use HookedController for compile-time safe hook registration.
		// The UserController.Register method will be called automatically
		// to register its lifecycle hooks.
		ligo.Controllers(ligo.HookedController(controller.NewUserController)),
	)
}
