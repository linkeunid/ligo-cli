package schematic

import (
	"fmt"
	"path/filepath"

	"github.com/linkeunid/ligo-cli/internal/templateutil"
)

type usecaseSchematic struct{}

func init() { Register(&usecaseSchematic{}) }

func (s *usecaseSchematic) Name() string      { return "usecase" }
func (s *usecaseSchematic) Aliases() []string { return []string{"uc"} }

func (s *usecaseSchematic) Run(ctx Context) error {
	d := makeData(ctx)
	n := templateutil.NormalizeName(ctx.Name)

	if ctx.Full {
		if err := (&dtoSchematic{}).Run(ctx); err != nil {
			return err
		}
		if err := templateutil.RenderToFile(
			usecaseErrorsTmpl,
			filepath.Join(ctx.WorkDir, "internal", "usecase", "errors.go"),
			nil, ctx.DryRun,
		); err != nil {
			return err
		}
		return templateutil.RenderToFile(
			usecaseTmpl,
			filepath.Join(ctx.WorkDir, "internal", "usecase", fmt.Sprintf("%s.go", n.Snake)),
			d, ctx.DryRun,
		)
	}

	return templateutil.RenderToFile(
		simpleUsecaseTmpl,
		filepath.Join(ctx.WorkDir, "internal", "usecase", fmt.Sprintf("%s.go", n.Snake)),
		d, ctx.DryRun,
	)
}
