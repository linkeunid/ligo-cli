package schematic

import (
	"fmt"
	"path/filepath"

	"github.com/linkeunid/ligo-cli/internal/templateutil"
)

type entitySchematic struct{}

func init() { Register(&entitySchematic{}) }

func (s *entitySchematic) Name() string      { return "entity" }
func (s *entitySchematic) Aliases() []string { return []string{"en"} }

func (s *entitySchematic) Run(ctx Context) error {
	d := makeData(ctx)
	n := templateutil.NormalizeName(ctx.Name)
	base := ctx.WorkDir

	if err := templateutil.RenderToFile(
		entityTmpl,
		filepath.Join(base, "internal", "domain", "entity", fmt.Sprintf("%s.go", n.Snake)),
		d, ctx.DryRun,
	); err != nil {
		return err
	}

	return templateutil.RenderToFile(
		repositoryTmpl,
		filepath.Join(base, "internal", "domain", "repository", fmt.Sprintf("%s.go", n.Snake)),
		d, ctx.DryRun,
	)
}
