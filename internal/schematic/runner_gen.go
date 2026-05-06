package schematic

import (
	"fmt"
	"path/filepath"

	"github.com/linkeunid/ligo-cli/internal/templateutil"
)

type runnerSchematic struct{}

func init() { Register(&runnerSchematic{}) }

func (s *runnerSchematic) Name() string      { return "runner" }
func (s *runnerSchematic) Aliases() []string { return []string{"run"} }

func (s *runnerSchematic) Run(ctx Context) error {
	d := makeData(ctx)
	n := templateutil.NormalizeName(ctx.Name)

	files := []struct {
		tmpl string
		path string
	}{
		{runnerMainTmpl, filepath.Join(ctx.WorkDir, "cmd", "runner", n.Snake, "main.go")},
		{runnerModuleTmpl, filepath.Join(ctx.WorkDir, "internal", n.Snake, "module", fmt.Sprintf("%s.go", n.Snake))},
		{runnerUseCaseTmpl, filepath.Join(ctx.WorkDir, "internal", n.Snake, "usecase", fmt.Sprintf("%s.go", n.Snake))},
		{runnerControllerTmpl, filepath.Join(ctx.WorkDir, "internal", n.Snake, "worker", "controller.go")},
	}

	for _, f := range files {
		if err := templateutil.RenderToFile(f.tmpl, f.path, d, ctx.DryRun); err != nil {
			return err
		}
	}

	return askAndRegisterRunner(ctx, n)
}

func askAndRegisterRunner(ctx Context, n templateutil.Names) error {
	// For now, just inform the user they can run the runner
	fmt.Printf("\n  INFO To run this runner: go run cmd/runner/%s/main.go\n", n.Snake)
	return nil
}

const runnerMainTmpl = `package main

import (
	"os"

	"github.com/linkeunid/ligo"
	"{{.ModulePath}}/internal/{{.Snake}}"
)

func main() {
	log := ligo.NewLogger()

	app := ligo.New(
		ligo.WithLogger(log),
		ligo.OnStart(func(ctx any) error {
			log.Info("{{.Pascal}} runner starting")
			return nil
		}),
		ligo.OnStop(func(ctx any) error {
			log.Info("{{.Pascal}} runner stopping")
			return nil
		}),
	)

	app.Register({{.Snake}}.Module())

	log.Info("Press Ctrl+C to stop")
	if err := app.Run(); err != nil {
		log.Error("Runner error", ligo.LoggerField{Key: "error", Value: err.Error()})
		os.Exit(1)
	}

	log.Info("Runner stopped gracefully")
}
`

const runnerModuleTmpl = `package module

import (
	"github.com/linkeunid/ligo"
	"{{.ModulePath}}/internal/{{.Snake}}/usecase"
	"{{.ModulePath}}/internal/{{.Snake}}/worker"
)

// Module returns the {{.Pascal}} runner module.
func Module() ligo.Module {
	return ligo.NewModule("{{.Snake}}",
		ligo.Providers(
			ligo.Factory[*usecase.{{.Pascal}}UseCase](usecase.New{{.Pascal}}UseCase),
		),
		ligo.Controllers(worker.NewController),
	)
}
`

const runnerUseCaseTmpl = `package usecase

import (
	"github.com/linkeunid/ligo"
)

type {{.Pascal}}UseCase struct {
	log ligo.Logger
}

func New{{.Pascal}}UseCase(log ligo.Logger) *{{.Pascal}}UseCase {
	return &{{.Pascal}}UseCase{log: log}
}

func (uc *{{.Pascal}}UseCase) Execute() {
	uc.log.Info("{{.Pascal}} worker executing task")
	// Add your worker logic here
}
`

const runnerControllerTmpl = `package worker

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/linkeunid/ligo"
	"{{.ModulePath}}/internal/{{.Snake}}/usecase"
)

type Controller struct {
	usecase *usecase.{{.Pascal}}UseCase
	log     ligo.Logger
	cancel  context.CancelFunc
	running atomic.Bool
}

func NewController(uc *usecase.{{.Pascal}}UseCase, log ligo.Logger) *Controller {
	return &Controller{
		usecase: uc,
		log:     log,
	}
}

func (c *Controller) Start() error {
	if c.running.Load() {
		return nil
	}

	c.log.Info("{{.Pascal}} worker controller starting")

	ctx, cancel := context.WithCancel(context.Background())
	c.cancel = cancel
	c.running.Store(true)

	go c.run(ctx)
	return nil
}

func (c *Controller) Stop() error {
	if !c.running.Load() {
		return nil
	}

	c.log.Info("{{.Pascal}} worker controller stopping")
	c.running.Store(false)
	if c.cancel != nil {
		c.cancel()
	}
	return nil
}

func (c *Controller) run(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	c.usecase.Execute()

	for {
		select {
		case <-ctx.Done():
			c.log.Info("{{.Pascal}} worker stopped")
			return
		case <-ticker.C:
			c.usecase.Execute()
		}
	}
}
`
