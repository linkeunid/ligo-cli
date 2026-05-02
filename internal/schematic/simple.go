package schematic

import (
	"fmt"
	"path/filepath"

	"github.com/linkeunid/ligo-cli/internal/templateutil"
)

type simpleSchematic struct{}

func init() { Register(&simpleSchematic{}) }

func (s *simpleSchematic) Name() string      { return "simple" }
func (s *simpleSchematic) Aliases() []string { return []string{"si"} }

func (s *simpleSchematic) Run(ctx Context) error {
	d := makeData(ctx)
	n := templateutil.NormalizeName(ctx.Name)

	if err := templateutil.RenderToFile(simpleUsecaseTmpl,
		filepath.Join(ctx.WorkDir, "internal", "usecase", fmt.Sprintf("%s.go", n.Snake)),
		d, ctx.DryRun,
	); err != nil {
		return err
	}

	if err := templateutil.RenderToFile(simpleControllerTmpl,
		filepath.Join(ctx.WorkDir, "internal", "infrastructure", "http", "controller",
			fmt.Sprintf("%s.go", n.Snake)),
		d, ctx.DryRun,
	); err != nil {
		return err
	}

	if err := templateutil.RenderToFile(simpleModuleTmpl,
		filepath.Join(ctx.WorkDir, "internal", "module", fmt.Sprintf("%s.go", n.Snake)),
		d, ctx.DryRun,
	); err != nil {
		return err
	}

	return askAndRegister(ctx, n.Pascal)
}
