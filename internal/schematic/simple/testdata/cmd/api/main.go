package main

import (
	"net/http"

	"github.com/linkeunid/ligo"
	"{{.ModulePath}}/internal/config"
	"{{.ModulePath}}/internal/infrastructure/http/middleware"
	"{{.ModulePath}}/internal/module"
	"github.com/linkeunid/ligo/adapters/echo"
)

func main() {
	cfg := config.NewConfig()
	log := ligo.NewLogger()

	router := echo.NewAdapter()

	app := ligo.New(
		ligo.WithRouter(router),
		ligo.WithAddr(cfg.ServerAddr),
		ligo.WithAutoPort(),
		ligo.WithLogger(log),
		ligo.WithMiddleware(
			middleware.CORS(),
			middleware.Recovery(),
		),
		ligo.OnStart(func(ctx any) error {
			log.Info("Server starting", ligo.LoggerField{Key: "addr", Value: cfg.ServerAddr})
			return nil
		}),
		ligo.OnStop(func(ctx any) error {
			log.Info("Server stopping")
			return nil
		}),
	)

	app.Register(module.Main())

	if err := app.Run(); err != nil {
		if err != http.ErrServerClosed {
			panic(err)
		}
	}
}
