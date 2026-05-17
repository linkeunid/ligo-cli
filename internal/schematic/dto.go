package schematic

import (
	"fmt"
	"path/filepath"

	"github.com/linkeunid/ligo-cli/internal/templateutil"
)

type dtoSchematic struct{}

func init() { Register(&dtoSchematic{}) }

func (s *dtoSchematic) Name() string      { return "dto" }
func (s *dtoSchematic) Aliases() []string { return []string{} }

func (s *dtoSchematic) Run(ctx Context) error {
	d := makeData(ctx)
	n := templateutil.NormalizeName(ctx.Name)
	dtoDir := filepath.Join(ctx.WorkDir, "internal", "usecase", "dto")

	if err := templateutil.RenderToFile(
		dtoCreateTmpl,
		filepath.Join(dtoDir, fmt.Sprintf("create_%s.go", n.Snake)),
		d, ctx.DryRun,
	); err != nil {
		return err
	}

	return templateutil.RenderToFile(
		dtoUpdateTmpl,
		filepath.Join(dtoDir, fmt.Sprintf("update_%s.go", n.Snake)),
		d, ctx.DryRun,
	)
}
