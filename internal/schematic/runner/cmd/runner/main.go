package main

import (
	"os"

	"{{.ModulePath}}/internal/config"
	"{{.ModulePath}}/internal/module"

	"github.com/linkeunid/ligo"
)

func main() {
	cfg := config.NewConfig()
	log := ligo.NewLogger()

	app := ligo.New(
		ligo.WithLogger(log),
		ligo.OnStart(func(ctx any) error {
			log.Info("Runner starting", ligo.LoggerField{Key: "version", Value: cfg.Version})
			return nil
		}),
		ligo.OnStop(func(ctx any) error {
			log.Info("Runner stopping")
			return nil
		}),
	)

	app.Provide(ligo.Value(cfg))
	app.Register(module.RunnerModule())

	log.Info("Press Ctrl+C to stop")
	if err := app.Run(); err != nil {
		log.Error("Runner error", ligo.LoggerField{Key: "error", Value: err.Error()})
		os.Exit(1)
	}

	log.Info("Runner stopped gracefully")
}
