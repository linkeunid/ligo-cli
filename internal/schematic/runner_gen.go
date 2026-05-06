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
		{runnerModuleTmpl, filepath.Join(ctx.WorkDir, "internal", n.Snake, "module.go")},
		{runnerUseCaseTmpl, filepath.Join(ctx.WorkDir, "internal", n.Snake, "usecase.go")},
		{runnerControllerTmpl, filepath.Join(ctx.WorkDir, "internal", n.Snake, "worker.go")},
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
	fmt.Printf("\n  INFO To run this runner: ligo work %s\n", n.Snake)
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

	app.Register({{.Snake}}.{{.Pascal}}Module())

	log.Info("Press Ctrl+C to stop")
	if err := app.Run(); err != nil {
		log.Error("Runner error", ligo.LoggerField{Key: "error", Value: err.Error()})
		os.Exit(1)
	}

	log.Info("Runner stopped gracefully")
}
`

const runnerModuleTmpl = `package {{.Snake}}

import (
	"github.com/linkeunid/ligo"
)

// {{.Pascal}}Module returns the {{.Pascal}} runner module.
func {{.Pascal}}Module() ligo.Module {
	return ligo.NewModule("{{.Snake}}",
		ligo.Providers(
			ligo.Factory[*{{.Pascal}}UseCase](New{{.Pascal}}UseCase),
		),
		ligo.Controllers(ligo.HookedController(New{{.Pascal}}Controller)),
	)
}
`

const runnerUseCaseTmpl = `package {{.Snake}}

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

const runnerControllerTmpl = `package {{.Snake}}

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/linkeunid/ligo"
)

type {{.Pascal}}Controller struct {
	{{.Snake}}UseCase *{{.Pascal}}UseCase
	log              ligo.Logger
	cancel           context.CancelFunc
	mu               sync.Mutex
	running          atomic.Bool
	wg               sync.WaitGroup
}

func New{{.Pascal}}Controller(uc *{{.Pascal}}UseCase, log ligo.Logger) *{{.Pascal}}Controller {
	return &{{.Pascal}}Controller{
		{{.Snake}}UseCase: uc,
		log:               log,
	}
}

func (c *{{.Pascal}}Controller) Initialize() error {
	c.log.Info("{{.Pascal}} worker controller initializing")
	return nil
}

func (c *{{.Pascal}}Controller) Start() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.running.Load() {
		return nil
	}

	c.log.Info("{{.Pascal}} worker controller starting")

	ctx, cancel := context.WithCancel(context.Background())
	c.cancel = cancel
	c.running.Store(true)

	c.wg.Add(1)
	go c.run(ctx)
	return nil
}

func (c *{{.Pascal}}Controller) Drain() error {
	c.log.Info("{{.Pascal}} worker controller draining - waiting for current work to complete")
	if c.cancel != nil {
		c.cancel()
	}
	c.wg.Wait()
	c.log.Info("{{.Pascal}} worker controller drained")
	return nil
}

func (c *{{.Pascal}}Controller) Stop() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.running.Load() {
		return nil
	}

	c.log.Info("{{.Pascal}} worker controller stopping")
	c.running.Store(false)
	if c.cancel != nil {
		c.cancel()
	}
	c.wg.Wait()
	c.log.Info("{{.Pascal}} worker controller stopped")
	return nil
}

func (c *{{.Pascal}}Controller) Register(registry *ligo.HookRegistry) {
	registry.OnInit(c.Initialize)
	registry.OnBootstrap(c.Start)
	registry.BeforeShutdown(c.Drain)
	registry.OnShutdown(c.Stop)
}

func (c *{{.Pascal}}Controller) run(ctx context.Context) {
	defer c.wg.Done()

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	c.{{.Snake}}UseCase.Execute()

	for {
		select {
		case <-ctx.Done():
			c.log.Info("{{.Pascal}} worker stopped")
			return
		case <-ticker.C:
			c.{{.Snake}}UseCase.Execute()
		}
	}
}
`
