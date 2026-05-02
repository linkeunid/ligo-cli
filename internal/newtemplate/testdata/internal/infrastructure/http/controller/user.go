package controller

import (
	"github.com/linkeunid/ligo"
	infraauth "{{.ModulePath}}/internal/infrastructure/auth"
	"{{.ModulePath}}/internal/infrastructure/http/middleware"
	"{{.ModulePath}}/internal/infrastructure/http/presenter"
	"{{.ModulePath}}/internal/usecase"
	"{{.ModulePath}}/internal/usecase/dto"
)

// UserController handles HTTP requests and route bindings for user operations.
type UserController struct {
	userUseCase *usecase.UserUseCase
	presenter   *presenter.UserPresenter
	log         ligo.Logger
	authGuard   ligo.Guard
	adminGuard  ligo.Guard
	exceptionMW ligo.Middleware
	loggingMW   ligo.Middleware
	auditMW     ligo.Middleware
}

// NewUserController creates a new user controller.
func NewUserController(uc *usecase.UserUseCase, pre *presenter.UserPresenter, jwt *infraauth.JWTAuth, log ligo.Logger) *UserController {
	return &UserController{
		userUseCase: uc,
		presenter:   pre,
		log:         log,
		authGuard:   infraauth.AuthGuard(jwt),
		adminGuard:  infraauth.AdminGuard(),
		exceptionMW: middleware.ExceptionMiddleware(log),
		loggingMW:   middleware.LoggingMiddleware(log),
		auditMW:     middleware.AuditMiddleware(log),
	}
}

// Routes registers all routes for the user controller.
func (c *UserController) Routes(r ligo.Router) {
	cr := ligo.NewChainRouter(r.Group("/users"))
	cr.Use(c.exceptionMW, c.loggingMW)

	cr.GET("", c.GetAllUsers).Handle()

	cr.GET("/:id", c.GetUserByID).
		Guard(c.authGuard).
		Pipe(ligo.ParseIntPipe("id")).
		Handle()

	cr.POST("", c.CreateUser).
		Guard(c.authGuard).
		Pipe(ligo.ValidationPipe(&dto.CreateUserInput{})).
		Handle()

	cr.PUT("/:id", c.UpdateUser).
		Guard(c.authGuard).
		Pipe(ligo.ParseIntPipe("id"), ligo.ValidationPipe(&dto.UpdateUserInput{})).
		Handle()

	cr.DELETE("/:id", c.DeleteUser).
		Guard(c.authGuard, c.adminGuard).
		Use(c.auditMW).
		Pipe(ligo.ParseIntPipe("id")).
		Handle()
}

// GetAllUsers handles GET /users
func (c *UserController) GetAllUsers(ctx ligo.Context) error {
	users := c.userUseCase.GetAllUsers()
	return ctx.OK(c.presenter.ToListResponse(users))
}

// GetUserByID handles GET /users/:id
func (c *UserController) GetUserByID(ctx ligo.Context) error {
	id := ligo.Get[int](ctx, "id")
	user, err := c.userUseCase.GetUserByID(id)
	if err != nil {
		return err
	}
	return ctx.OK(c.presenter.ToResponse(user))
}

// CreateUser handles POST /users
func (c *UserController) CreateUser(ctx ligo.Context) error {
	input := ligo.ValidatedBody[dto.CreateUserInput](ctx)
	user, err := c.userUseCase.CreateUser(*input)
	if err != nil {
		return err
	}
	return ctx.Created(c.presenter.ToResponse(user))
}

// UpdateUser handles PUT /users/:id
func (c *UserController) UpdateUser(ctx ligo.Context) error {
	id := ligo.Get[int](ctx, "id")
	input := ligo.ValidatedBody[dto.UpdateUserInput](ctx)
	user, err := c.userUseCase.UpdateUser(id, *input)
	if err != nil {
		return err
	}
	return ctx.OK(c.presenter.ToResponse(user))
}

// DeleteUser handles DELETE /users/:id
func (c *UserController) DeleteUser(ctx ligo.Context) error {
	id := ligo.Get[int](ctx, "id")
	if err := c.userUseCase.DeleteUser(id); err != nil {
		return err
	}
	return ctx.NoContent()
}
