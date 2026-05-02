package controller

import (
	"github.com/linkeunid/ligo"
	"{{.ModulePath}}/internal/infrastructure/http/middleware"
	"{{.ModulePath}}/internal/usecase"
)

// HelloController handles HTTP requests for hello operations.
type HelloController struct {
	uc          *usecase.HelloUseCase
	log         ligo.Logger
	exceptionMW ligo.Middleware
	loggingMW   ligo.Middleware
}

// NewHelloController creates a new HelloController.
func NewHelloController(uc *usecase.HelloUseCase, log ligo.Logger) *HelloController {
	return &HelloController{
		uc:          uc,
		log:         log,
		exceptionMW: middleware.ExceptionMiddleware(log),
		loggingMW:   middleware.LoggingMiddleware(log),
	}
}

// Routes registers all routes for the HelloController.
func (c *HelloController) Routes(r ligo.Router) {
	cr := ligo.NewChainRouter(r.Group("/api"))
	cr.Use(c.exceptionMW, c.loggingMW)

	cr.GET("", c.Hello).Handle()
}

func (c *HelloController) Hello(ctx ligo.Context) error {
	return ctx.OK(c.uc.Hello())
}
